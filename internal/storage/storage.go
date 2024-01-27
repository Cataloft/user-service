package storage

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/Cataloft/user-service/internal/model"
	"github.com/Cataloft/user-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func New(dbURL string) *Storage {
	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Error: parse db url")
	}

	var connPool *pgxpool.Pool

	for {
		connPool, _ = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if connPool.Ping(context.Background()) == nil {
			log.Println("postgres upped")

			break
		}

		time.Sleep(1 * time.Second)
	}

	sqlDB := stdlib.OpenDBFromPool(connPool)
	UpMigrations(sqlDB)

	return &Storage{Conn: connPool}
}

func (s *Storage) SaveUser(u *model.User) error {
	queryCreate := "INSERT INTO public.users (name, surname, patronymic, gender, age, nationality) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := s.Conn.Exec(context.Background(),
		queryCreate, u.Name, u.Surname, u.Patronymic, u.Gender, u.Age, u.Nationality)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	exists, err := s.ExistsByID(id)
	if err != nil {
		return err
	}

	if !exists {
		log.Printf("User with id=%d is not exist", id)

		return nil
	}

	queryDelete := "DELETE FROM users where id = $1"
	_, err = s.Conn.Exec(context.Background(), queryDelete, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(id int, user *model.User) error {
	exists, err := s.ExistsByID(id)

	if err != nil {
		return err
	}

	if !exists {
		log.Printf("User with id=%d is not exist", id)

		return nil
	}

	queryUpdate := "UPDATE users SET"
	tail, args := utils.ProcessUserFields(id, user)
	queryUpdate += tail
	_, err = s.Conn.Exec(context.Background(), queryUpdate, args...)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUsers(filters []string, logger *slog.Logger) (*[]model.User, error) {
	var users []model.User

	var u model.User

	queryGet := "SELECT * FROM users"

	if len(filters) != 0 {
		queryGet = queryGet + " " + "WHERE"

		for i, filter := range filters {
			if i > 0 {
				queryGet = queryGet + " " + "AND"
			}

			queryGet = queryGet + " " + filter
		}

		logger.Debug("Sql request: %s", "sql", queryGet)

		rows, err := s.Conn.Query(context.Background(), queryGet)

		if err != nil {
			logger.Error("Error querying db", "error", err)

			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Gender, &u.Age, &u.Nationality)
			if err != nil {
				logger.Error("Error scanning row", "error", err)

				return nil, err
			}

			users = append(users, u)
		}

		return &users, nil
	}

	rows, err := s.Conn.Query(context.Background(), queryGet)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Gender, &u.Age, &u.Nationality)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return &users, nil
}

func (s *Storage) ExistsByID(id int) (bool, error) {
	queryExist := "SELECT EXISTS(SELECT * FROM users WHERE id = $1)"

	var exists bool

	row := s.Conn.QueryRow(context.Background(), queryExist, id)
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

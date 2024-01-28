package storage

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/Cataloft/user-service/internal/config"
	"github.com/Cataloft/user-service/internal/model"
	"github.com/Cataloft/user-service/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	Conn   *pgxpool.Pool
	logger *slog.Logger
}

func New(cfg config.Database, logger *slog.Logger) *Storage {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error: parse db url")
	}

	var connPool *pgxpool.Pool

	for i := cfg.MaxAttempts; i > 0; i-- {
		connPool, _ = pgxpool.NewWithConfig(context.Background(), poolCfg)
		if connPool.Ping(context.Background()) == nil {
			logger.Debug("postgres upped")

			break
		}

		time.Sleep(cfg.DurationAttempts)
	}

	sqlDB := stdlib.OpenDBFromPool(connPool)
	UpMigrations(sqlDB)

	return &Storage{
		Conn:   connPool,
		logger: logger,
	}
}

func (s *Storage) SaveUser(u *model.User) error {
	queryCreate := "INSERT INTO public.users (name, surname, patronymic, gender, age, nationality) VALUES ($1, $2, $3, $4, $5, $6)"
	s.logger.Debug("Create query", "query", queryCreate)

	_, err := s.Conn.Exec(context.Background(),
		queryCreate, u.Name, u.Surname, u.Patronymic, u.Gender, u.Age, u.Nationality)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(id int) (string, error) {
	exists, err := s.ExistsByID(id)
	if err != nil {
		return "", err
	}

	if !exists {
		s.logger.Debug("User is not exist", "id", id)

		return "user not exist", nil
	}

	queryDelete := "DELETE FROM users where id = $1"
	s.logger.Debug("Delete query", "query", queryDelete)

	_, err = s.Conn.Exec(context.Background(), queryDelete, id)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (s *Storage) UpdateUser(id int, user *model.User) (string, error) {
	exists, err := s.ExistsByID(id)
	if err != nil {
		return "", err
	}

	if !exists {
		s.logger.Debug("User is not exist", "id", id)

		return "user not exist", err
	}

	queryUpdate := "UPDATE users SET"
	s.logger.Debug("Update query", "query", queryUpdate)

	tail, args := utils.ProcessUserFields(id, user)
	queryUpdate += tail
	_, err = s.Conn.Exec(context.Background(), queryUpdate, args...)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (s *Storage) GetUsers(filters []string) ([]model.User, error) {
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

		s.logger.Debug("Get query", "query", queryGet)

		rows, err := s.Conn.Query(context.Background(), queryGet)
		if err != nil {
			s.logger.Error("Error querying db", "error", err)

			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Gender, &u.Age, &u.Nationality)
			if err != nil {
				s.logger.Error("Error scanning row", "error", err)

				return nil, err
			}

			users = append(users, u)
		}

		return users, nil
	}

	rows, err := s.Conn.Query(context.Background(), queryGet)
	if err != nil {
		s.logger.Error("Error querying db", "error", err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Gender, &u.Age, &u.Nationality)
		if err != nil {
			s.logger.Error("Error scanning row", "error", err)

			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
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

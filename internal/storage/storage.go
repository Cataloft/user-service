package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"strings"
	"test-task/internal/model"
	"test-task/internal/utils"
)

type Storage struct {
	Conn *pgxpool.Pool
}

func New(dbUrl string) *Storage {
	poolCfg, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalln("ERROR: parse db url")
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		log.Fatalln("ERROR: connect to db ")
	}

	return &Storage{Conn: connPool}
}

func (s *Storage) SaveUser(u *model.User) error {
	sqlCreate := "INSERT INTO public.users (name, surname, patronymic, gender, age, nationality) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := s.Conn.Exec(context.Background(), sqlCreate, u.Name, u.Surname, u.Patronymic, u.Gender, u.Age, u.Nationality)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	sqlDelete := "DELETE FROM users where id = $1"
	_, err := s.Conn.Exec(context.Background(), sqlDelete, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUser(id int, user *model.User) error {
	sqlUpdate, args := utils.ProcessUserFields(id, user)

	_, err := s.Conn.Exec(context.Background(), sqlUpdate, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUsers(filters []string, log *slog.Logger) (*[]model.User, error) {
	var users []model.User
	var u model.User
	sqlGet := "SELECT * FROM users"

	if len(filters) != 0 {
		sqlGet = strings.Join([]string{sqlGet, "WHERE"}, " ")
		for i, filter := range filters {
			if i > 0 {
				sqlGet = strings.Join([]string{sqlGet, "AND"}, " ")
			}
			sqlGet = strings.Join([]string{sqlGet, filter}, " ")
		}

		log.Debug(fmt.Sprintf("sql request: %s", sqlGet))

		rows, err := s.Conn.Query(context.Background(), sqlGet)
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

	rows, err := s.Conn.Query(context.Background(), sqlGet)
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

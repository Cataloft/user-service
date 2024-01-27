package utils

import (
	"fmt"
	"github.com/Cataloft/user-service/internal/model"
)

func ProcessUserFields(id int, u *model.User) (string, []any) {
	query := "UPDATE users SET "
	var args []any
	argIndex := 1

	if u.Name != "" {
		query += fmt.Sprintf("name=$%d, ", argIndex)
		args = append(args, u.Name)
		argIndex++
	}
	if u.Surname != "" {
		query += fmt.Sprintf("surname=$%d, ", argIndex)
		args = append(args, u.Surname)
		argIndex++
	}
	if u.Patronymic != "" {
		query += fmt.Sprintf("patronymic=$%d, ", argIndex)
		args = append(args, u.Patronymic)
		argIndex++
	}
	if u.Gender != "" {
		query += fmt.Sprintf("gender=$%d, ", argIndex)
		args = append(args, u.Gender)
		argIndex++
	}
	if u.Age != 0 {
		query += fmt.Sprintf("age=$%d, ", argIndex)
		args = append(args, u.Age)
		argIndex++
	}
	if u.Nationality != "" {
		query += fmt.Sprintf("nationality=$%d, ", argIndex)
		args = append(args, u.Nationality)
		argIndex++
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id=$%d", argIndex)
	args = append(args, id)

	return query, args
}

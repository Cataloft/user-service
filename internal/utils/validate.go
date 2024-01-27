package utils

import (
	"fmt"
	"github.com/Cataloft/user-service/internal/model"
)

func ProcessUserFields(id int, u *model.User) (string, []any) {
	tail := ""
	var args []any
	argIndex := 1

	if u.Name != "" {
		tail += fmt.Sprintf(" name=$%d, ", argIndex)
		args = append(args, u.Name)
		argIndex++
	}
	if u.Surname != "" {
		tail += fmt.Sprintf(" surname=$%d,", argIndex)
		args = append(args, u.Surname)
		argIndex++
	}
	if u.Patronymic != "" {
		tail += fmt.Sprintf(" patronymic=$%d,", argIndex)
		args = append(args, u.Patronymic)
		argIndex++
	}
	if u.Gender != "" {
		tail += fmt.Sprintf(" gender=$%d,", argIndex)
		args = append(args, u.Gender)
		argIndex++
	}
	if u.Age != 0 {
		tail += fmt.Sprintf(" age=$%d,", argIndex)
		args = append(args, u.Age)
		argIndex++
	}
	if u.Nationality != "" {
		tail += fmt.Sprintf(" nationality=$%d,", argIndex)
		args = append(args, u.Nationality)
		argIndex++
	}
	args = append(args, id)
	tail = tail[:len(tail)-1] + fmt.Sprintf(" WHERE id=$%d", argIndex)

	return tail, args
}

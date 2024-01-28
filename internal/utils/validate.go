package utils

import (
	"fmt"

	"github.com/Cataloft/user-service/internal/model"
)

func addCondition(tail string, argIndex int, field string, value any, args []any) (string, int, []any) {
	if value != "" && value != 0 {
		tail += fmt.Sprintf(" %s=$%d,", field, argIndex)

		args = append(args, value)

		argIndex++
	}

	return tail, argIndex, args
}

func ProcessUserFields(id int, u *model.User) (tail string, args []any) {
	tail = ""
	argIndex := 1
	tail, argIndex, args = addCondition(tail, argIndex, "name", u.Name, args)
	tail, argIndex, args = addCondition(tail, argIndex, "surname", u.Surname, args)
	tail, argIndex, args = addCondition(tail, argIndex, "patronymic", u.Patronymic, args)
	tail, argIndex, args = addCondition(tail, argIndex, "gender", u.Gender, args)
	tail, argIndex, args = addCondition(tail, argIndex, "age", u.Age, args)
	tail, argIndex, args = addCondition(tail, argIndex, "nationality", u.Nationality, args)

	if len(tail) > 0 {
		tail = tail[:len(tail)-1] + fmt.Sprintf(" WHERE id=$%d", argIndex)

		args = append(args, id)
	}

	return tail, args
}

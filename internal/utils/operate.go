package utils

import (
	"strings"
)

func OperateStrings(field, fieldVal, op string) string {
	filterOps := map[string]string{
		"greater": ">",
		"lower":   "<",
		"equal":   "=",
		"like":    "LIKE",
	}

	fieldAliases := map[string]string{
		"ageGreater":  "age",
		"ageLower":    "age",
		"nameContain": "name",
	}

	switch op {
	case "greater", "lower":
		return fieldAliases[field] + filterOps[op] + fieldVal
	case "like":
		return strings.Join(
			[]string{fieldAliases[field], filterOps[op], "'%" + fieldVal + "%'", ""}, " ")
	default:
		return field + filterOps[op] + "'" + fieldVal + "'"
	}
}

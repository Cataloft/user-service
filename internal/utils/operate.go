package utils

import (
	"strings"
)

func OperateStrings(field string, fieldVal string, op string) string {
	filterOps := map[string]string{
		"greater": ">",
		"lower":   "<",
		"equal":   "=",
		"like":    "LIKE",
	}

	switch op {
	case "greater", "lower":
		return field + filterOps[op] + fieldVal
	case "like":
		return strings.Join(
			[]string{field, filterOps[op], strings.Join([]string{"'%", fieldVal, "%'"}, "")}, " ")
	default:
		return field + filterOps[op] + fieldVal
	}
}

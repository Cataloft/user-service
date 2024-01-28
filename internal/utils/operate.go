package utils

import (
	"strings"
)

func OperateStrings(field, fieldVal string) string {
	filterOps := map[string]string{
		"ageGreater":  ">",
		"ageLower":    "<",
		"nameContain": "LIKE",
	}

	fieldAliases := map[string]string{
		"ageGreater":  "age",
		"ageLower":    "age",
		"nameContain": "name",
	}

	switch field {
	case "ageGreater", "ageLower":
		return fieldAliases[field] + filterOps[field] + fieldVal
	case "nameContain":
		return strings.Join([]string{fieldAliases[field], filterOps[field], "'%" + fieldVal + "%'", ""}, " ")
	default:
		return field + "=" + "'" + fieldVal + "'"
	}
}

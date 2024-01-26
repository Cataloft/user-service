package utils

import "strings"

func OperateStrings(field string, fieldVal string, op string) string {
	filterOps := map[string]string{
		"greater": ">",
		"lower":   "<",
		"equal":   "=",
	}

	if field == "age" && op != "" {
		handledStr := strings.Join([]string{field, filterOps[op], fieldVal}, "")
		return handledStr
	}

	handledStr := strings.Join([]string{field, filterOps["equal"], fieldVal}, "")

	return handledStr
}

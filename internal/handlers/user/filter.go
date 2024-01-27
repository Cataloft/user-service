package user

import (
	"log/slog"
	"net/http"

	"github.com/Cataloft/user-service/internal/utils"

	"github.com/Cataloft/user-service/internal/model"
	"github.com/gin-gonic/gin"
)

type GetterUser interface {
	GetUsers(filters []string, log *slog.Logger) (*[]model.User, error)
}

func GetList(userGetter GetterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		filterParams := map[string]string{
			"ageGreater":  "greater",
			"ageLower":    "lower",
			"nameContain": "like",
			"age":         "equal",
			"name":        "equal",
			"surname":     "equal",
			"patronymic":  "equal",
			"gender":      "equal",
			"nationality": "equal",
		}

		var filters []string

		for param, op := range filterParams {
			if paramFilter := c.Query(param); paramFilter != "" {
				filters = append(filters, utils.OperateStrings(param, paramFilter, op))
			}
		}

		log.Debug("Getting filters", "filters", filters)

		filteredUsers, err := userGetter.GetUsers(filters, log)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting data from server", "error": err})
			log.Error("Error getting data from server", "error", err)

			return
		}

		items := *filteredUsers
		log.Debug("Filtered data", "items", items)

		page, ok := c.MustGet("page").(int)
		if !ok {
			log.Error("Type assertion failed")
		}

		pageSize, ok := c.MustGet("pageSize").(int)
		if !ok {
			log.Error("Type assertion failed")
		}

		startIndex := (page - 1) * pageSize
		endIndex := startIndex + pageSize

		if endIndex > len(items) {
			endIndex = len(items)
		}

		paginatedItems := items[startIndex:endIndex]
		log.Debug("Paginated data", "items", paginatedItems)

		c.JSON(http.StatusOK, gin.H{"page": page, "pageSize": pageSize, "filteredUsers": paginatedItems})
		log.Info("Paginated data successfully sent")
	}
}

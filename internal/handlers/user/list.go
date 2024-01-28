package user

import (
	"log/slog"
	"net/http"

	"github.com/Cataloft/user-service/internal/model"
	"github.com/Cataloft/user-service/internal/utils"
	"github.com/gin-gonic/gin"
)

type GetterUser interface {
	GetUsers(filters []string) ([]model.User, error)
}

func GetList(userGetter GetterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log = log.With(slog.String("op", "handlers.user.GetList"))
		params := []string{"ageGreater", "ageLower", "nameContain", "age", "name", "surname", "patronymic", "gender", "nationality"}

		var filters []string

		for _, param := range params {
			if paramVal := c.Query(param); paramVal != "" {
				filters = append(filters, utils.OperateStrings(param, paramVal))
			}
		}

		log.Debug("Getting filters", "filters", filters)

		filteredUsers, err := userGetter.GetUsers(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting data from server", "error": err})
			log.Error("Error getting data from server", "error", err)

			return
		}

		log.Debug("Filtered data", "items", filteredUsers)

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

		if endIndex > len(filteredUsers) {
			endIndex = len(filteredUsers)
		}

		if startIndex > endIndex {
			c.JSON(http.StatusOK, gin.H{"message": "This page doesn't exist"})

			return
		}

		paginatedItems := filteredUsers[startIndex:endIndex]
		log.Debug("Paginated data", "items", paginatedItems)

		c.JSON(http.StatusOK, gin.H{"page": page, "pageSize": pageSize, "filteredUsers": paginatedItems})
		log.Info("Paginated data successfully sent")
	}
}

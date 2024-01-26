package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"test-task/internal/model"
	"test-task/internal/utils"
)

type GetterUser interface {
	GetUsers(filters []string, log *slog.Logger) (*[]model.User, error)
}

func Filter(userGetter GetterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters []string

		if paramAge := c.Query("age"); paramAge != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("age", paramAge, paramFilter))
		}
		if paramName := c.Query("name"); paramName != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("name", fmt.Sprintf("'%s'", paramName), paramFilter))
		}
		if paramSurname := c.Query("surname"); paramSurname != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("surname", fmt.Sprintf("'%s'", paramSurname), paramFilter))
		}
		if paramPatronymic := c.Query("patronymic"); paramPatronymic != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("patronymic", fmt.Sprintf("'%s'", paramPatronymic), paramFilter))
		}
		if paramGender := c.Query("gender"); paramGender != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("gender", fmt.Sprintf("'%s'", paramGender), paramFilter))
		}
		if paramNationality := c.Query("nationality"); paramNationality != "" {
			paramFilter := c.Query("filter")
			filters = append(filters, utils.OperateStrings("nationality", fmt.Sprintf("'%s'", paramNationality), paramFilter))
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

		page := c.MustGet("page").(int)
		pageSize := c.MustGet("pageSize").(int)

		startIndex := (page - 1) * pageSize
		endIndex := startIndex + pageSize

		if endIndex > len(items) {
			endIndex = len(items)
		}

		paginatedItems := items[startIndex:endIndex]
		log.Debug("Paginated data", "items", paginatedItems)

		c.JSON(http.StatusOK, gin.H{
			"page":          page,
			"pageSize":      pageSize,
			"filteredUsers": paginatedItems,
		})
		log.Info("Paginated data successfully sent")
	}
}

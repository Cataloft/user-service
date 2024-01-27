package user

import (
	"fmt"
	"github.com/Cataloft/user-service/internal/model"
	"github.com/Cataloft/user-service/internal/utils"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type GetterUser interface {
	GetUsers(filters []string, log *slog.Logger) (*[]model.User, error)
}

func GetFiltered(userGetter GetterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filters []string

		if paramFilter := c.Query("ageGreater"); paramFilter != "" {
			filters = append(filters, utils.OperateStrings("age", paramFilter, "greater"))
		}
		if paramFilter := c.Query("ageLower"); paramFilter != "" {
			filters = append(filters, utils.OperateStrings("age", paramFilter, "lower"))
		}
		if paramFilter := c.Query("nameContain"); paramFilter != "" {
			filters = append(filters, utils.OperateStrings("name", paramFilter, "like"))
		}
		if paramAge := c.Query("age"); paramAge != "" {
			filters = append(filters, utils.OperateStrings("age", paramAge, "equal"))
		}
		if paramName := c.Query("name"); paramName != "" {
			filters = append(filters, utils.OperateStrings("name", fmt.Sprintf("'%s'", paramName), "equal"))
		}
		if paramSurname := c.Query("surname"); paramSurname != "" {
			filters = append(filters, utils.OperateStrings("surname", fmt.Sprintf("'%s'", paramSurname), "equal"))
		}
		if paramPatronymic := c.Query("patronymic"); paramPatronymic != "" {
			filters = append(filters, utils.OperateStrings("patronymic", fmt.Sprintf("'%s'", paramPatronymic), "equal"))
		}
		if paramGender := c.Query("gender"); paramGender != "" {
			filters = append(filters, utils.OperateStrings("gender", fmt.Sprintf("'%s'", paramGender), "equal"))
		}
		if paramNationality := c.Query("nationality"); paramNationality != "" {
			filters = append(filters, utils.OperateStrings("nationality", fmt.Sprintf("'%s'", paramNationality), "equal"))
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

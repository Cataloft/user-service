package user

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"test-task/internal/model"
)

//type UpdateRequest struct {
//	Name        string `json:"name"`
//	Surname     string `json:"surname"`
//	Patronymic  string `json:"patronymic"`
//	Gender      string `json:"gender"`
//	Age         int    `json:"age"`
//	Nationality string `json:"nationality"`
//}

type UpdaterUser interface {
	UpdateUser(id int, user *model.User) error
}

func Update(updater UpdaterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramKey := c.Param("id")
		id, err := strconv.Atoi(paramKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing param", "error": err})
			log.Error("Error parsing param", "error", err)
			return
		}

		var user model.User
		err = c.ShouldBindJSON(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error request data", "error": err})
			log.Error("Error request data", "error", err)
			return
		}
		log.Debug("Updating fields", "fields:", user)

		err = updater.UpdateUser(id, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating data", "error": err})
			log.Error("Error updating data", "error", err)
			return
		}

		c.JSON(http.StatusOK, id)
		log.Info("User successfully updated")
	}
}

package user

import (
	"github.com/Cataloft/user-service/internal/model"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

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
		log.Debug("Updating fields", "fields", user)

		err = updater.UpdateUser(id, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating data", "error": err})
			log.Error("Error updating data", "error", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
		log.Info("User successfully updated")
	}
}

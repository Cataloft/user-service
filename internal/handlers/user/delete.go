package user

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleterUser interface {
	DeleteUser(id int) (string, error)
}

func Delete(deleter DeleterUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramKey := c.Param("id")

		id, err := strconv.Atoi(paramKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing param", "error": err})
			log.Error("Error parsing param", "error", err)

			return
		}

		exist, err := deleter.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user", "error": err})
			log.Error("Error deleting user", "error", err)

			return
		}

		if exist != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "User is not exist", "id": id})
			log.Info("User is not exist", "id", id)

			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
		log.Info("User successfully deleted")
	}
}

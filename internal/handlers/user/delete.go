package user

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type DeleterUser interface {
	DeleteUser(id int) error
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

		err = deleter.DeleteUser(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting user", "error": err})
			log.Error("Error deleting user", "error", err)

			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
		log.Info("User successfully deleted")
	}
}

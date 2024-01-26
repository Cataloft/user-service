package user

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sync"
	"test-task/internal/handlers/enrich"
	"test-task/internal/model"
)

type Request struct {
	Name       string `json:"name" binding:"required"`
	Surname    string `json:"surname" binding:"required"`
	Patronymic string `json:"patronymic"`
}

type SaverUser interface {
	SaveUser(user *model.User) error
}

func EnrichFIO(saver SaverUser, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Request
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			log.Error("Error requested data", "error", err)
			return
		}

		user := model.User{
			Name:       req.Name,
			Surname:    req.Surname,
			Patronymic: req.Patronymic,
		}

		wg := sync.WaitGroup{}
		wg.Add(3)
		go func(c *gin.Context) {
			defer wg.Done()
			user.Age = enrich.EnrichAge(c, user.Name)
		}(c)
		go func(c *gin.Context) {
			defer wg.Done()
			user.Gender = enrich.EnrichGender(c, user.Name)
		}(c)
		go func(c *gin.Context) {
			defer wg.Done()
			user.Nationality = enrich.EnrichNationality(c, user.Name)
		}(c)
		wg.Wait()

		log.Debug("User data enriched", "enriched user", user)

		err = saver.SaveUser(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user", "error": err})
			log.Error("Error saving user", "error", err)

			return
		}

		c.JSON(http.StatusCreated, user.ID)
		log.Info("User successfully created")
	}
}

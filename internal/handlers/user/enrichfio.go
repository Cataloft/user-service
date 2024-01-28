package user

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/Cataloft/user-service/internal/apis"
	"github.com/Cataloft/user-service/internal/model"
	"github.com/gin-gonic/gin"
)

const numAPIs = 3

type Request struct {
	Name       string `binding:"required" json:"name"`
	Surname    string `binding:"required" json:"surname"`
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

		enricher := apis.New(log)

		wg := sync.WaitGroup{}
		wg.Add(numAPIs)

		go func() {
			defer wg.Done()

			user.Age = enricher.EnrichAge(c, user.Name)
		}()
		go func() {
			defer wg.Done()

			user.Gender = enricher.EnrichGender(c, user.Name)
		}()
		go func() {
			defer wg.Done()

			user.Nationality = enricher.EnrichNationality(c, user.Name)
		}()
		wg.Wait()

		log.Debug("User data enriched", "enriched user", user)

		err = saver.SaveUser(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user", "error": err})
			log.Error("Error saving user", "error", err)

			return
		}

		c.JSON(http.StatusCreated, gin.H{"created": user})
		log.Info("User successfully created")
	}
}

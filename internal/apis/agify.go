package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgeResp struct {
	Age int `json:"age"`
}

func (e *Enricher) EnrichAge(c *gin.Context, name string) int {
	var age AgeResp

	apiURL := "https://api.agify.io/?name="
	url := apiURL + name

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		e.logger.Error("Error executing request", "error", err)

		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		e.logger.Error("Error getting data", "error", err)

		return -1
	}

	err = json.NewDecoder(resp.Body).Decode(&age)

	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		e.logger.Error("Error decoding data", "error", err)

		return -1
	}

	e.logger.Debug("Got age", "age", age.Age)

	return age.Age
}

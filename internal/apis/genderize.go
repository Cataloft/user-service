package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenderResp struct {
	Gender string `json:"gender"`
}

func (e *Enricher) EnrichGender(c *gin.Context, name string) string {
	var gender GenderResp

	apiURL := "https://api.genderize.io/?name="
	url := apiURL + name

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		e.logger.Error("Error executing request", "error", err)

		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		e.logger.Error("Error getting data", "error", err)

		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&gender)

	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		e.logger.Error("Error decoding data", "error", err)

		return ""
	}

	e.logger.Debug("Got gender", "gender", gender.Gender)

	return gender.Gender
}

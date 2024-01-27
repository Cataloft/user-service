package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

func EnrichGender(c *gin.Context, name string, log *slog.Logger) string {
	type GenderResp struct {
		Gender string `json:"gender"`
	}
	var gender GenderResp

	apiUrl := "https://api.genderize.io/?name="
	url := strings.Join([]string{apiUrl, name}, "")

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		log.Error("Error executing request", "error", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		log.Error("Error getting data", "error", err)
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&gender)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Error("Error decoding data", "error", err)
		return ""
	}

	log.Debug("Got gender", "gender", gender.Gender)
	return gender.Gender
}

package apis

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

func EnrichAge(c *gin.Context, name string, log *slog.Logger) int {
	type AgeResp struct {
		Age int `json:"age"`
	}

	apiUrl := "https://api.agify.io/?name="
	url := strings.Join([]string{apiUrl, name}, "")

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		log.Error("Error executing request", "error", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		log.Error("Error getting data", "error", err)
		return -1
	}

	var age AgeResp
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Error("Error decoding data", "error", err)
		return -1
	}

	log.Debug("Got age", "age", age.Age)
	return age.Age
}

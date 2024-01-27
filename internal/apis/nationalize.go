package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

func EnrichNationality(c *gin.Context, name string, log *slog.Logger) string {
	type NationalityResp struct {
		Country []struct {
			CountryId   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	var nationality NationalityResp

	apiUrl := "https://api.nationalize.io/?name="
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

	err = json.NewDecoder(resp.Body).Decode(&nationality)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Error("Error decoding data", "error", err)
		return ""
	}

	maxProbability := 0.0
	maxNationality := ""
	for _, v := range nationality.Country {
		if v.Probability > maxProbability {
			maxProbability = v.Probability
			maxNationality = v.CountryId
		}
	}

	log.Debug("Got nationality", "nationality", maxNationality)
	return maxNationality
}

package apis

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NationalityResp struct {
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func (e *Enricher) EnrichNationality(c *gin.Context, name string) string {
	var nationality NationalityResp

	apiURL := "https://api.nationalize.io/?name="
	url := apiURL + name

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		e.logger.Error("Error executing request", "error", err)

		return ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "apiURL": apiURL})
		e.logger.Error("Error getting data", "apiURL", apiURL)

		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&nationality)

	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		e.logger.Error("Error decoding data", "error", err)

		return ""
	}

	maxProbability := 0.0
	maxNationality := ""

	for _, v := range nationality.Country {
		if v.Probability > maxProbability {
			maxProbability = v.Probability
			maxNationality = v.CountryID
		}
	}

	e.logger.Debug("Got nationality", "nationality", maxNationality)

	return maxNationality
}

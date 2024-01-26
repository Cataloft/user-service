package enrich

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func EnrichAge(c *gin.Context, name string) int {
	type AgeResp struct {
		Age int `json:"age"`
	}

	apiUrl := "https://api.agify.io/?name="
	url := strings.Join([]string{apiUrl, name}, "")

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		log.Println("Error executing request")
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		return -1
	}

	var age AgeResp
	err = json.NewDecoder(resp.Body).Decode(&age)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Println("failed decoding json data")
		return -1
	}

	return age.Age
}

func EnrichGender(c *gin.Context, name string) string {
	type GenderResp struct {
		Gender string `json:"gender"`
	}
	var gender GenderResp

	apiUrl := "https://api.genderize.io/?name="
	url := strings.Join([]string{apiUrl, name}, "")

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error executing request", "error": err})
		log.Println("Error executing request")
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&gender)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Println("failed decoding json data")
		return ""
	}

	return gender.Gender
}

func EnrichNationality(c *gin.Context, name string) string {
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
		log.Println("Error executing request")
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"message": "Error getting data", "error": err})
		return ""
	}

	err = json.NewDecoder(resp.Body).Decode(&nationality)
	if err != nil {
		c.JSON(resp.StatusCode, gin.H{"message": "Error decoding data", "error": err})
		log.Println("failed decoding json data")
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

	return maxNationality
}

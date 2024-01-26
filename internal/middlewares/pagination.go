package middlewares

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.Query("pageSize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		c.Set("page", page)
		c.Set("pageSize", pageSize)
		c.Next()
	}
}

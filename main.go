package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Saludo struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

func main() {
	var err error
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
	router.POST("/saludar", func(c *gin.Context) {
		var s Saludo
		if err := c.BindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		mensaje := fmt.Sprintf("Hola %v %v", s.Name, s.LastName)

		c.JSON(http.StatusOK, gin.H{"mensaje": mensaje})
	})
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
	router.Run()
}

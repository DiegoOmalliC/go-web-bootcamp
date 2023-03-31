package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var products []Product

func readJson(path string) error {
	var (
		err       error
		file      *os.File
		byteValue []byte
	)
	file, err = os.Open(path)
	if err != nil {
		return errors.New("Not found file")
	}
	defer file.Close()
	byteValue, err = io.ReadAll(file)
	if err != nil {
		return errors.New("Cannot read File")
	}
	err = json.Unmarshal(byteValue, &products)
	if err != nil {
		return errors.New("Not load json")
	}
	return nil
}

func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"products": products})
}
func getId(c *gin.Context) {
	id := c.Param("id")
	for _, a := range products {
		idProduct, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		if a.ID == idProduct {
			c.JSON(http.StatusOK, gin.H{"product": a})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
}
func getProduct(c *gin.Context) {
	price := c.Query("price")
	if price == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Price is required"})
		return
	}
	priceProduct, err := strconv.ParseFloat(price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var filteredProducts []Product
	for _, v := range products {
		if v.Price > priceProduct {
			filteredProducts = append(filteredProducts, v)
		}
	}
	c.IndentedJSON(http.StatusOK, gin.H{"products": filteredProducts})
}

func main() {
	var err error
	err = readJson("./exercise_2/products.json")
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	router.GET("/products", getAll)
	router.GET("/products/:id", getId)
	router.GET("/products/search", getProduct)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	CodeValue   string  `json:"code_value" binding:"required"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
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
func lastID() int {
	if len(products) == 0 {
		return 0
	}
	return products[len(products)-1].ID
}
func validCode(codeValue string) bool {
	for _, p := range products {
		if p.CodeValue == codeValue {
			return false
		}
	}
	return true
}
func validDate(date string) bool {
	var dayInt, monthInt, yearInt int
	if len(strings.Split(date, "/")) != 3 {
		return false
	}

	day := strings.Split(date, "/")[0]
	month := strings.Split(date, "/")[1]
	year := strings.Split(date, "/")[2]

	if len(day) != 2 || len(month) != 2 || len(year) != 4 {
		return false
	}
	if v, err := strconv.Atoi(strings.Split(date, "/")[0]); err == nil {
		dayInt = v
	}
	if v, err := strconv.Atoi(strings.Split(date, "/")[1]); err == nil {
		monthInt = v
	}
	if v, err := strconv.Atoi(strings.Split(date, "/")[2]); err == nil {
		yearInt = v
	}
	if dayInt < 1 || dayInt > 31 {
		return false
	}
	if monthInt < 1 || monthInt > 12 {
		return false
	}
	if yearInt < 1 || yearInt > 9999 {
		return false
	}
	return true
}
func addProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot bind json"})
		log.Printf("Error: %v", err)
		return
	}
	if !validCode(product.CodeValue) {
		c.JSON(http.StatusConflict, gin.H{"error": "Code value Product already exists"})
		return
	}
	if !validDate(product.Expiration) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Expiration Date is invalid"})
		return
	}
	product.ID = lastID() + 1
	products = append(products, product)
	c.JSON(http.StatusCreated, gin.H{"product": product})
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
	router.POST("/product", addProduct)
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
	if err = router.Run(":8080"); err != nil {
		panic(err)
	}
}

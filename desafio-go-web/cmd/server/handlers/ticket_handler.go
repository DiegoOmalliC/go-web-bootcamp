package handlers

import (
	"net/http"

	"github.com/bootcamp-go/desafio-go-web/internal/tickets"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service tickets.Service
}

func NewTicketHandler(s tickets.Service) *TicketHandler {
	return &TicketHandler{
		service: s,
	}
}

func (s *TicketHandler) GetTicketsByCountry() gin.HandlerFunc {
	return func(c *gin.Context) {

		destination := c.Param("dest")

		quantity, err := s.service.GetTotalTickets(c, destination)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		c.JSON(200, gin.H{
			"quantity":    quantity,
			"destination": destination,
		})
	}
}

func (s *TicketHandler) AverageDestination() gin.HandlerFunc {
	return func(c *gin.Context) {

		destination := c.Param("dest")

		avg, err := s.service.AverageDestination(c, destination)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error(), nil)
			return
		}

		c.JSON(200, gin.H{
			"avg":         avg,
			"destination": destination,
		})
	}
}

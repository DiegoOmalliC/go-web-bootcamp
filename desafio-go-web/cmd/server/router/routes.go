package router

import (
	"github.com/bootcamp-go/desafio-go-web/cmd/server/handlers"
	"github.com/bootcamp-go/desafio-go-web/internal/domain"
	"github.com/bootcamp-go/desafio-go-web/internal/tickets"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
	Data   []domain.Ticket
}

func NewRouter(server *gin.Engine, data []domain.Ticket) *Router {
	return &Router{
		Engine: server,
		Data:   data,
	}
}
func (r Router) MapRoutes() {
	repository := tickets.NewRepository(r.Data)
	service := tickets.NewService(repository)
	handler := handlers.NewTicketHandler(service)
	r.Engine.GET("/ticket/getByCountry/:dest", handler.GetTicketsByCountry())
	r.Engine.GET("/ticket/getAverage/:dest", handler.AverageDestination())
}

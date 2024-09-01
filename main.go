package main

import (
	"github.com/LetsFocus/goLF/goLF"
	invoice2 "github.com/LetsFocus/template-service/handlers"
	invoice3 "github.com/LetsFocus/template-service/services"
	"github.com/LetsFocus/template-service/stores"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	g := gin.Default()

	golf := goLF.New()
	store := stores.New(golf.Postgres)
	service := invoice3.New(store)
	handler := invoice2.New(service)

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "tenantId", "X-Correlation-Id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	g.Use(cors.New(config), SetResponseHeaders())

	// user endpoints
	g.GET("/template", handler.Get)
	g.POST("/template", handler.Create)
	g.GET("/template/:id", handler.GetByID)
	g.DELETE("/template/:id", handler.Delete)
	g.PATCH("/template/:id", handler.Patch)

	err := g.Run(":"+g.Configs.Get("HTTP_PORT"))
	if err != nil {
		return
	}
}

func SetResponseHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

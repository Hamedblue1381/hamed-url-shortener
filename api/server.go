package api

import (
	"fmt"
	"net/http"

	"github.com/Hamedblue1381/hamed-url-shortener/token"
	"github.com/Hamedblue1381/hamed-url-shortener/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cant create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	router.POST("/register", server.Register)
	router.POST("/login", server.Login)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	{
		authRoutes.POST("shorten", server.ShortenURL)
		authRoutes.DELETE("shorten/:shortenedURL", server.DeleteShortURL)
	}

	router.GET(":shortenedURL", server.RedirectShortURL)
	router.GET("/views/:shortenedURL", GetShortUrlViews)

	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package api

import (
	"net/http"

	"github.com/Hamedblue1381/hamed-url-shortener/db/model"
	"github.com/Hamedblue1381/hamed-url-shortener/token"
	"github.com/gin-gonic/gin"
)

type shortenURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

func (server *Server) ShortenURL(c *gin.Context) {
	var request shortenURLRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := model.GetUserByID(authPayload.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	shortenedURL, err := model.ShortenURL(request.OriginalURL, user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"shortened_url": server.config.HTTPServerAddress + "/" + shortenedURL})
}

func (server *Server) RedirectShortURL(c *gin.Context) {
	shortenedURL := c.Param("shortenedURL")

	originalURL, err := model.RedirectURL(shortenedURL, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shortened URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (server *Server) DeleteShortURL(c *gin.Context) {
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	shortenedURL := c.Param("shortenedURL")

	err := model.DeleteShortenedURL(shortenedURL, authPayload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shortened URL deleted successfully"})
}

func GetShortUrlViews(c *gin.Context) {
	shortenedURL := c.Param("shortenedURL")

	clickedCount, err := model.GetClickedCount(shortenedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"views": clickedCount})
}

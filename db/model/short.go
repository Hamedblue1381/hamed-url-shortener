package model

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Hamedblue1381/hamed-url-shortener/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShortenURL(originalURL string, userID uint64) (string, error) {
	var count int64
	result := db.Model(&ShortUrl{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		return "failed to shorten: ", result.Error
	}
	if count >= 10 {
		err := fmt.Errorf("you reached the maximum of shortened url's allowed")
		return "failed to shorten: ", err
	}
	shortenedURL := util.RandomURL(6)

	// Store the URL mapping with the shortened URL.
	shortURL := ShortUrl{
		Redirect:  originalURL,
		Shortened: shortenedURL,
		UserID:    userID,
	}

	result = db.Create(&shortURL)
	if result.Error != nil {
		return "failed to shorten: ", result.Error
	}

	return shortenedURL, nil
}

// RedirectURL takes a shortened URL and redirects to the original URL with parameters.
func RedirectURL(shortenedURL string, c *gin.Context) (string, error) {
	// Retrieve the original URL from the database.
	originalURL, err := getOriginalURL(shortenedURL)
	if err != nil {
		return "", err
	}

	if err := db.Model(&ShortUrl{}).Where("shortened = ?", shortenedURL).Update("clicked", gorm.Expr("clicked + ?", 1)).Error; err != nil {
		return "", err
	}
	// Extract and append query parameters to the original URL.
	originalURLWithParams, err := appendQueryParameters(originalURL, getQueryParameters(c.Request.URL))
	if err != nil {
		return "", err
	}

	return originalURLWithParams, nil
}

func appendQueryParameters(originalURL string, params url.Values) (string, error) {
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "failed to parse url:", err
	}

	query := parsedURL.Query()
	for key, values := range params {
		for _, value := range values {
			query.Add(key, value)
		}
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}
func getQueryParameters(shortenedURL *url.URL) url.Values {
	params := shortenedURL.Query()

	queryParameters := url.Values(params)

	return queryParameters
}

func getOriginalURL(shortenedURL string) (string, error) {
	var shortURL ShortUrl
	result := db.Where("shortened = ?", shortenedURL).First(&shortURL)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("shortened url not found")
		}
		return "", result.Error
	}

	return shortURL.Redirect, nil
}
func GetClickedCount(shortenedURL string) (uint64, error) {
	var clickedCount uint64

	result := db.Model(&ShortUrl{}).Select("clicked").Where("shortened = ?", shortenedURL).Scan(&clickedCount)
	if result.Error != nil {
		return 0, result.Error
	}
	return clickedCount, nil
}
func DeleteShortenedURL(shortenedURL string, userID uint64) error {
	// Verify user
	var shortURL ShortUrl
	result := db.Where("shortened = ? AND user_id = ?", shortenedURL, userID).First(&shortURL)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return errors.New("shortened url not found or does not belong to the user")
		}
		return result.Error
	}

	// Delete the shortened URL from the database.
	result = db.Delete(&shortURL)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

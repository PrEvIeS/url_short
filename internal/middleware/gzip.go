package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	contentEncoding = "Content-Encoding"
	acceptEncoding  = "Accept-Encoding"
	gzipEncoding    = "gzip"
	contentType     = "Content-Type"
	applicationJSON = "application/json"
	textHTML        = "text/html"
)

func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader(contentEncoding), gzipEncoding) {
			decompressGzip(c)
		}

		if shouldCompressResponse(c) {
			compressResponse(c)
		}

		c.Next()
	}
}

func decompressGzip(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	gz, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid gzip body"})
		return
	}
	defer func() {
		if closeErr := gz.Close(); closeErr != nil {
			log.Printf("gzip close error: %v", closeErr)
		}
	}()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "decompression failed"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewReader(decompressed))
	c.Request.Header.Del(contentEncoding)
}

func shouldCompressResponse(c *gin.Context) bool {
	return strings.Contains(c.GetHeader(acceptEncoding), gzipEncoding) &&
		(strings.Contains(c.GetHeader(contentType), applicationJSON) ||
			strings.Contains(c.GetHeader(contentType), textHTML))
}

func compressResponse(c *gin.Context) {
	gz := gzip.NewWriter(c.Writer)
	defer func() {
		if closeErr := gz.Close(); closeErr != nil {
			log.Printf("gzip close error: %v", closeErr)
		}
	}()

	c.Header(contentEncoding, gzipEncoding)
	c.Writer = &gzipResponseWriter{Writer: gz, ResponseWriter: c.Writer}
}

type gzipResponseWriter struct {
	io.Writer
	gin.ResponseWriter
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	n, err := g.Writer.Write(data)
	if err != nil {
		return n, fmt.Errorf("gzip write failed: %w", err)
	}
	return n, nil
}

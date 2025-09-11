package web

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cursor Gen

//go:embed view/*
var viewFS embed.FS

//go:embed view/static/*
var staticFS embed.FS

func ServeEmbeddedFile(c *gin.Context, filePath string) {
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	data, err := viewFS.ReadFile("view" + filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	ext := strings.ToLower(path.Ext(filePath))
	var contentType string
	switch ext {
	case ".html":
		contentType = "text/html; charset=utf-8"
	case ".css":
		contentType = "text/css; charset=utf-8"
	case ".js":
		contentType = "application/javascript; charset=utf-8"
	case ".json":
		contentType = "application/json; charset=utf-8"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	default:
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}

func ServeEmbeddedStatic(c *gin.Context) {
	filePath := c.Param("filepath")
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}
	data, err := staticFS.ReadFile("view/static" + filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	ext := strings.ToLower(path.Ext(filePath))
	var contentType string
	switch ext {
	case ".html":
		contentType = "text/html; charset=utf-8"
	case ".css":
		contentType = "text/css; charset=utf-8"
	case ".js":
		contentType = "application/javascript; charset=utf-8"
	case ".json":
		contentType = "application/json; charset=utf-8"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	default:
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}

func GetEmbeddedFileSystem() fs.FS {
	return viewFS
}

func GetEmbeddedStaticFileSystem() fs.FS {
	return staticFS
}

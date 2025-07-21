package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) renewAccessToken(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
} 
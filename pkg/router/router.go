package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xesina/golang-realworld/pkg/env"
)

// Engine setup and return a gin Engine
func Engine(e string) *gin.Engine {
	if env.IsRelease(e) {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	if env.IsDevelopment(e) {
		r.Use(gin.Logger())
	}
	return r
}

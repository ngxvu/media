package router

import (
	limit "github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"media-service/services/errors"
)

func Route() *gin.Engine {
	r := gin.Default()
	applyMiddlewares(r)
	r.Static("/images", "./images")
	Routes(r)
	return r
}

func applyMiddlewares(r *gin.Engine) {
	r.Use(limit.MaxAllowed(200))
	r.Use(errors.HandlerError)
}

func Routes(r *gin.Engine) {

	v1 := r.Group("/v1")
	{
		api := v1.Group("/api")
		{
			MediaRoutes(api, NewMediaHandler())
		}
	}
}

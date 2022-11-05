package web

import (
	"net/http"
	"salmon-fish/web/controller"
	"salmon-fish/web/handler"
	"time"

	"github.com/gin-gonic/gin"
)

// 启动Web服务并指定路由信息
func WebStart(app controller.Application) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	h := handler.NewHandler()
	apiv1 := router.Group("/api/v1")
	{
		apiv1.POST("/addFish", h.AddFish)
		apiv1.GET("/queryInfoByFishId", h.QueryInfoByFishId)
	}

	s := &http.Server{
		Addr:           ":9000",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}

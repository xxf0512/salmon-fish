package handler

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AddFish(c *gin.Context) {

}

func (h *Handler) QueryInfoByFishId(c *gin.Context) {

}

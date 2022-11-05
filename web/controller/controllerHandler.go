package controller

import (
	"net/http"
	"salmon-fish/service"

	"github.com/gin-gonic/gin"
)

func (app *Application) AddFish(c *gin.Context) {
	var fish service.SalmonFish
	if err := c.ShouldBindJSON(&fish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// for i:=0; i<len(Users); i++ {
	// 	if Users[i].LoginName == ""
	// }
	// if json.User != "manu" || json.Password != "123" {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	// 	return
	// }

	payload, err := app.Setup.SaveFish(fish)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"payload": payload,
	}
	c.JSON(http.StatusOK, data)
}

func (app *Application) FindFishByFishId(c *gin.Context) {
	fishId := c.Query("id")
	payload, err := app.Setup.FindFishByFishId(fishId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"payload": payload,
	}
	c.JSON(http.StatusOK, data)
}
package main

import (
	"github.com/gin-gonic/gin"
	"./api"
)

var PORT string = ":8099"

func main() {

	// MONGO DB
	api.DBConnection()

	r := gin.Default()
	r.Use(api.CORSMiddleware())
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "SUCCESS",
		})
	})
	
	// Send koimetrics script
	r.GET( "/API/v1/:key/script.js", 			api.Stats_script) 
	
	// Receive client statistics
	r.POST("/API/v1/statistics/", 			 	api.Statistics )  
	
	// Update client session alive
	r.POST("/API/v1/heartbeats/", 			 	api.Heart_beats)  
	r.POST("/DJANGO/APIKEY/register", 			api.Register_apikey)
	r.POST("/DJANGO/asked_location_websites/",  api.Update_asked_location_websites)
	r.GET( "/DJANGO/get_analytics/", 			api.Analytics_between_dates)
	r.GET( "/DJANGO/get_sessions/", 			api.Current_sessions)
	r.Run(PORT)
}

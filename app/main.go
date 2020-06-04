package main

import (
	"github.com/gin-gonic/gin"
	"github.com/koimetrics-api/app/api"
)


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
	r.GET( "/API/v1/:key/koimetrics.js", 			api.VisitorScript) 
	
	// Receive client statistics
	r.POST("/API/v1/statistics/", 			 	api.VisitorResults )  
	
	// Update client session alive
	r.POST("/API/v1/heartbeats/", 			 	api.HeartBeats)  
	
	// Register users api key
	r.GET("/DJANGO/new_apikey", 			api.RegisterApikey)
	
	// Deprecated
	//r.POST("/DJANGO/asked_location_websites/",  api.Update_asked_location_websites)
	
	// Deprecated
	r.GET( "/DJANGO/get_analytics/", 			api.AnalyticsBetweenDates)
	
	r.GET( "/DJANGO/get_sessions/", 			api.CurrentSessions)

	var PORT string = api.GoDotEnvVariable("GOAPI_PORT")
	r.Run(PORT)
}

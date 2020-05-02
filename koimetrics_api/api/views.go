package api

import (
	"context"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"strconv"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"time"
	"github.com/twinj/uuid"
)

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
  

// CLIENT BROWSER ENDPOINTS
func js_parser(template_name string, context map[string]string) string {
	client_flea, err := ioutil.ReadFile(template_name)
	if err == nil {
		parsed_js := string(client_flea)
		for s, r := range context {
			parsed_js = strings.Replace(parsed_js, s, r, -1)
		}
		return parsed_js
	} else {
		return ""
	}
}

func Stats_script(c *gin.Context) {
	key := c.Param("key")
	var apikey ApiKey
	filter := bson.D{{
		"$and", bson.A{
			bson.M{"key": key},
			bson.D{{"enddate", bson.D{{ "$gte", time.Now().Format(("2006-01-02")) }}}},
		},
	}}
	err := apikeys.FindOne(context.TODO(), filter).Decode(&apikey)
	if err == nil {
		goapi_host := goDotEnvVariable("GOAPI_HOST")
		session_id := uuid.NewV4()
		fmt.Println(session_id.String())
		c.Header("Content-Type", "application/javascript; charset=utf-8")
		
		context := map[string]string{
			"{{.key}}": key,
			"{{.session_id}}": session_id.String(),
			"{{.goapi_host}}": goapi_host,
			"{{.ask_location_to}}" : strings.Join(apikey.AskLocationTo, ","),
		}
		client_flea_str := js_parser("templates/client_flea.js", context)
		c.String(http.StatusOK, client_flea_str)
	} else {
		c.JSON(200, gin.H{
			"status":     "ERROR",
			"statusCode": 0,
		})
	}
}

func Statistics(c *gin.Context) {
	analytic := AnalyticResult{}
	analytic.Key  			= c.PostForm("Key")
	analytic.Host 			= c.PostForm("Host")
	analytic.Path 			= c.PostForm("Path")
	analytic.Date 			= c.PostForm("Date")
	analytic.Referrer 		= c.PostForm("Referrer")
	analytic.ReferrerPath 	= c.PostForm("ReferrerPath")
	analytic.Time 			= c.PostForm("Time")
	analytic.Performance, _ = strconv.ParseFloat(c.PostForm("Performance"), 64)
	analytic.Latitude, _ 	= strconv.ParseFloat(c.PostForm("Latitude"), 64)
	analytic.Longitude, _ 	= strconv.ParseFloat(c.PostForm("Longitude"), 64)
	isPhoneInt, _ 			:= strconv.Atoi(c.PostForm("IsPhone"))
	analytic.IsPhone 		= isPhoneInt != 0
	analytic.Country 		= c.PostForm("Country")
	analytic.City 			= c.PostForm("City")
	analytic.Region 		= c.PostForm("Region")
	analytic.Session_id		= c.PostForm("session_id")
	analytic.Session_start 	= time.Now().Format("2006-01-02 15:04:05")
	
	insertResult, _ := analytics.InsertOne(context.TODO(), analytic)
	fmt.Println("New session started with session_id: ", analytic.Session_id)
	fmt.Println("New document inserted with ID: ", insertResult.InsertedID)
	
	c.JSON(200, gin.H{
		"status": "SUCCESS",
	})
}

func Heart_beats(c *gin.Context){
	session_id := c.PostForm("session_id")
	filter := bson.D{{"session_id", session_id}}
	session_end := time.Now().Format("2006-01-02 15:04:05")
	analytics.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"session_end": session_end}})
	fmt.Println(" Updated session end for session_id: ", session_id)
	c.JSON(200, gin.H{
		"status": "SUCCESS",
		"session_id": session_id,
	})
}


// WEB APP ENDPOINTS
//register_apikey: 
//  Receives a keyCode and an endDate in format "YYYY-mm-dd".
//   If keycode doesnt exist, creates a new apikey record in database with end date.
//   If keyCode does exist, send an error message telling the key code isnt available.

func Register_apikey(c *gin.Context) {
	newKeyCode := c.PostForm("keyCode")
	newEndDate := c.PostForm("endDate")
	var apikey ApiKey
	// Check if key exists
	filter := bson.D{{"key", newKeyCode}}
	err := apikeys.FindOne(context.TODO(), filter).Decode(&apikey)
	if err != nil {
		// err != nil means that couldnt find an api key with same code => KeyCode is available
		apikey.Key = newKeyCode
		apikey.AskLocationTo = []string{}
		apikey.EndDate = newEndDate
		insertResult, _ := apikeys.InsertOne(context.TODO(), apikey)
		fmt.Println("Inserted a New Key: ", insertResult.InsertedID)
		c.JSON(200, gin.H{
			"status":     "SUCCESS",
			"statusCode": 1,
			"data": gin.H{
				"api_key": newKeyCode,
				"end_date": newEndDate,
				"ID":  insertResult.InsertedID,
			},
		})
	} else {
		fmt.Println(err)
		fmt.Println(newKeyCode + "Already exists.")
		c.JSON(200, gin.H{
			"status":     "ERROR",
			"statusCode": 0,
		})
	}
}



func Update_asked_location_websites(c *gin.Context) {
	akey := c.PostForm("ApiKey")
	askLocationToData := c.PostForm("AskLocationTo")
	askLocationTo := strings.Split(askLocationToData, ",")
	filter := bson.D{{"key", akey}}
	updatedResult, err := apikeys.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"asklocationto": askLocationTo}})

	if err != nil {
		c.JSON(200, gin.H{
			"status": "ERROR",
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "SUCCESS",
			"updated": updatedResult,
		})
	}
}

func Analytics_between_dates(c *gin.Context) {
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	host := c.Query("host")
	apikey := c.Query("apikey")
	
	pipeline := bson.M{
		"$and": []interface{}{
			bson.M{
				"key": apikey,
			},
			bson.M{
				"host": host,
			},
			bson.M{"date": bson.M{
				"$gte": start_date,
			}},
			bson.M{"date": bson.M{
				"$lte": end_date,
			}},
		},
	}
	

	filterCursor, err := analytics.Find(context.TODO(), pipeline)

	var analyticsFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &analyticsFiltered); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"status":  "SUCCESS",
		"results": analyticsFiltered,
	})
}

func Current_sessions(c *gin.Context){
	host := c.Query("host")
	apikey := c.Query("apikey")
	seconds_ago, err := strconv.Atoi( c.Query("seconds") )
	min_session_end := time.Now().Add(time.Duration(-seconds_ago) * time.Second).Format("2006-01-02 15:04:05")
	pipeline := bson.M{
		"$and": []interface{}{
			bson.M{
				"key": apikey,
			},
			bson.M{
				"host": host,
			},
			bson.M{"session_end": bson.M{
				"$gte": min_session_end,
			}},
		},
	}
	
	filterCursor, err := analytics.Find(context.TODO(), pipeline)
	if err != nil {
		fmt.Println(err)
	}
	var analyticsFiltered []bson.M
	if err = filterCursor.All(context.TODO(), &analyticsFiltered); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"status":  "SUCCESS",
		"results": analyticsFiltered,
	})
}
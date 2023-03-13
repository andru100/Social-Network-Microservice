package main

import (
	"context"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
    "github.com/andru100/Social-Network-Microservices/backend/services/PostFile/utils"
    "github.com/andru100/Social-Network-Microservices/backend/services/PostFile/model"
)

func main() {
    router := gin.New()  
    router.Use(utils.CORSMiddleware())
    router.POST("/postfile/:userid", Postfile)// posts profile pic and users media
    router.Run(":4010")
    //router.RunTLS(":4001", "./server.pem", "./server.key")
}


func Postfile(c *gin.Context) {//Takes file from request form, runs upload func, puts in s3, returns s3 address.

    file, header, err := c.Request.FormFile("file") // get file from request body
    
    if err != nil {
        log.Println(err)
    }
 
    filename := header.Filename
    
    fileread, err := ioutil.ReadAll(file) 

    userid := c.Param("userid") // get id from url request
    
    requestType := c.PostForm("type")

    collection := utils.Client.Database("datingapp").Collection("userdata")

    imgaddress, err := utils.Uploaditem(&userid, &filename, &fileread)// upload func returns uploaded img url
    
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, err)
    }

    currentDoc := model.MongoFields{}

    ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
    
    // Find the document that mathes the id from the request.
    err = collection.FindOne(ctx, bson.M{"Username": userid}).Decode(&currentDoc)
    
    //create filter 
    filter := bson.M{"Username": userid}
    update := bson.D{}
    if requestType == "profPic" {
        currentDoc.Profpic= *imgaddress //replace url to profile pic URL
        update = bson.D{
            {"$set", bson.D{
                {"Profpic", currentDoc.Profpic},
            }},
        }
    } else if requestType == "addPhotos" {
        currentDoc.Photos= append(currentDoc.Photos, *imgaddress) //append to list of users photo urls
        update = bson.D{
            {"$set", bson.D{
                {"Photos", currentDoc.Photos},
            }},
        }
    }
    
    //put to db
    _, err = collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        log.Fatal(err)
    }

    c.IndentedJSON(http.StatusOK, currentDoc)
}
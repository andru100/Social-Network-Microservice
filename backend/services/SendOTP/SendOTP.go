package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/andru100/Social-Network/backend/social"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	router.Use(social.CORSMiddleware())
	router.POST("/SendOTP", SendOTP) // checks for authentication using jwt
	router.POST("/RecieveOTP", RecieveOTP) // checks for authentication using jwt
	router.Run(":4009")
	//router.RunTLS(":4001", "./server.pem", "./server.key")
}

func SendOTP (c *gin.Context) { //

	var reqbody social.Verification

	if err := c.BindJSON(&reqbody); err != nil {
		fmt.Println(err)
		return
	}


	//create otp
	nums := []rune("123456789")

	rand.Seed(time.Now().UnixNano())

    b := make([]rune, 6)
    for i := range b {
        b[i] = nums[rand.Intn(len(nums))]
    }

	otp := string(b) 
	fmt.Println("randon otp is", otp, "this isnt safe, wiill need some secret key to truly randomize")
	a := social.SendTxt("+447307868951", otp)
	fmt.Println(a, a==false)


}

func RecieveOTP (c *gin.Context) { //

	var request social.Verification

	if err := c.BindJSON(&reqbody); err != nil {
		fmt.Println(err)
		return
	}


	if !request.PwordReq {
		collection := social.Client.Database("datingapp").Collection("OTP") // connect to db and collection.

		currentDoc := social.MongoFields{}

		ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

		err := collection.FindOne(ctx, bson.M{"Username": qry.UserName}).Decode(&currentDoc)
		social.CheckError(err)
		c.IndentedJSON(http.StatusOK, currentDoc)
	}

	 
	fmt.Println("randon otp is", otp, "this isnt safe, wiill need some secret key to truly randomize")
	a := social.SendTxt("+447307868951", otp)
	fmt.Println(a, a==false)


}
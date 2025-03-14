package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Meh-Mehul/db-config-service/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson"
	// "context"
)
var collectionH *mongo.Collection
var collectionU *mongo.Collection
func init(){
	collectionH = controllers.Connect("hash")
	collectionU = controllers.Connect("uri")
}
func addHashHandler(c *gin.Context) {
	var input controllers.Hash
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// randHash := controllers.GetRandomHash()
	res, err := controllers.AddHashtoDB(collectionH, input.Rand, input.Filename,input.Path, input.Ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add hash"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Hash added", "id": res.InsertedID})
}

func addURIHandler(c *gin.Context){
	var ip controllers.URI
	if err := c.ShouldBindJSON(&ip); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := controllers.AddURItoDB(collectionU, ip.Uri, ip.Hashes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add URI"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "URI added", "id": res.InsertedID})



}

func getHashHandler(c *gin.Context) {
	hashID := c.Param("id") // Read hash ID from URL
	// objID, err := primitive.ObjectIDFromHex(hashID)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hash ID"})
	// 	return
	// }
	hash, err := controllers.GetHashFromDB(collectionH, hashID)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get or find hash"})
		return
	}
	c.JSON(http.StatusOK, hash)
}


func getURIHandler(c *gin.Context){
	uriId := c.Param("id");
	uri, err := controllers.GetURIFromDB(collectionU, uriId)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get or find URI"})
		return
	}
	c.JSON(http.StatusOK, uri)

}


func main(){
	r := gin.Default()
	r.POST("/hash", addHashHandler)
	r.GET("/hash/:id", getHashHandler)
	r.POST("/uri", addURIHandler)
	r.GET("/uri/:id", getURIHandler)
	if err := r.Run(":6000"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is running on port 8080...")
}


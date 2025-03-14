package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
	"log"
	// "github.com/gohugoio/hugo/tpl/collections"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// we need to first get a random hash
func GetRandomHash() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) 
	}
	return hex.EncodeToString(bytes)
}
// a functiom to create hashes objects and add to collection
func AddHashtoDB(collection *mongo.Collection, randhash string, filename string, path string,ext string) (*mongo.InsertOneResult, error) {
	newHash := Hash{
		ID: primitive.NewObjectID(),
		Rand: randhash,
		Filename: filename,
		Ext: ext,
		Path: path,
		Time: time.Now(),
	}
	res, err := collection.InsertOne(context.Background(), newHash);
	if err != nil {
		log.Println("Error inserting hash:", err)
		return nil, err
	}
	log.Println("Inserted hash with ID:", res.InsertedID)
	return res, nil
}


func AddURItoDB(collection *mongo.Collection, uri string, hashes []string) (*mongo.InsertOneResult, error) {
	newURI := URI{
		ID: primitive.NilObjectID,
		Uri: uri,
		Hashes: hashes,
	}
	res, err := collection.InsertOne(context.Background(), newURI)
	if err != nil {
		log.Println("Error inserting URI", err)
		return nil, err
	}
	log.Println("Insterted URI", newURI.ID)
	return res, nil
}



func GetHashFromDB(collection *mongo.Collection, hash string) (Hash, error) {
	var gottenHash Hash
	err := collection.FindOne(context.Background(), bson.M{"rand":hash}).Decode(&gottenHash)
	if(err!=nil){
		log.Println("Could not Find the Hash, Error: ", err)
		return gottenHash, err
	}
	return gottenHash, nil
}

func GetURIFromDB(collection *mongo.Collection, uri string) (URI, error) {
	var gottenURI URI
	err := collection.FindOne(context.Background(), bson.M{"uri":uri}).Decode(&gottenURI)
	if(err!=nil){
		log.Println("Could not Find the URI, Error: ", err)
		return gottenURI, err
	}
	return gottenURI, nil

}


// These are the Models for Hash and URI,
// TODO: Do same for reader-client
type Hash struct {
	ID 	 primitive.ObjectID `bson:"_id,omitempty"`
	Rand string				`bson:"rand"`
	Filename string 		`bson:"filename"`
	Path 	 string			`bson:"path"`
	Ext 	 string 		`bson:"ext"`
	Time 	 time.Time 		`bson:"time"`
}

// MAJOR ISSUE: I have not indexed the DB to have unique uri (although in pipeline, they will enncounter this situation will super low probablity)

type URI struct{
	ID	primitive.ObjectID 		`bson:"_id,omitempty"`
	Uri string 					`bson:"uri"`
	Hashes []string				`bson:"hashes"`
}













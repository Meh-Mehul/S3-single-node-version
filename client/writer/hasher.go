package writer

import (
	// "context"
	"crypto/rand"
	"encoding/hex"
	// "time"
	// "log"
	// "github.com/gohugoio/hugo/tpl/collections"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
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









// for refrence : 
// type Hash struct {
// 	ID 	 primitive.ObjectID `bson:"_id,omitempty"`
// 	Rand string				`bson:"rand"`
// 	Filename string 		`bson:"filename"`
// 	Path 	 string			`bson:"path"`
// 	Ext 	 string 		`bson:"ext"`
// 	Time 	 time.Time 		`bson:"time"`
// }

// // MAJOR ISSUE: I have not indexed the DB to have unique uri (although in pipeline, they will enncounter this situation will super low probablity)

// type URI struct{
// 	ID	primitive.ObjectID 		`bson:"_id,omitempty"`
// 	Uri string 					`bson:"uri"`
// 	Hashes []string				`bson:"hashes"`
// }













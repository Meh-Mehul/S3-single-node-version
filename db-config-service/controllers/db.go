package controllers


import(
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"context"
	"log"
	"fmt"
)


func Connect(name string) *mongo.Collection {
    err := godotenv.Load(".env")
    // if err != nil {
    //     log.Fatalf("Error loading .env file: %s", err)
    // }
    MONGO_URI := os.Getenv("MONGO_URI")
    clientOptions := options.Client().ApplyURI(MONGO_URI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }else{
      fmt.Println("Connected to mongoDB!!!")
   }
   collection := client.Database("configdb").Collection(name)
   if err != nil {
	   log.Fatal(err)
   }
   return collection;
}
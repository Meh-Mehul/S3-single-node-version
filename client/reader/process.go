package reader


import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"net/http"
	"errors"
)
// Location for config-db service
var s5 = "http://localhost:6000/"
type Hash struct {
	ID 	 primitive.ObjectID `json:"ID,omitempty"`
	Rand string				`json:"Rand"`
	Filename string 		`json:"Filename"`
	Path 	 string			`json:"Path"`
	Ext 	 string 		`json:"Ext"`
	Time 	 time.Time 		`json:"Time"`
}
type URI struct{
	ID	primitive.ObjectID 		`json:"ID,omitempty"`
	Uri string 					`json:"Uri"`
	Hashes []string				`json:"Hashes"`
}
func GetHashDetailsFromDB(hash string) (*Hash, error){
	url := s5 + "hash/"+hash;
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch hash details")
	}
	var hashDetails Hash
	if err := json.NewDecoder(resp.Body).Decode(&hashDetails); err != nil {
		return nil, err
	}
	return &hashDetails, nil
}


func GetURIDeatilsFromDB(uri string) (*URI, error){
	url := s5 + "uri/" + uri;
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch uri details")
	}
	var URIdetails URI;
	if err := json.NewDecoder(resp.Body).Decode(&URIdetails); err != nil {
		return nil, err
	}
	return &URIdetails, nil;
}


func CheckContainability(uri string, hash string) (*Hash, error){
	the_uri , err:= GetURIDeatilsFromDB(uri)
	if(err != nil){
		return nil, err
	}
	the_hash , err:= GetHashDetailsFromDB(hash)
	if(err != nil){
		return nil, err
	}
	chk := false
	for _, item := range the_uri.Hashes {
        if item == the_hash.Rand {
            chk = true
        }
    }
	if(!chk){
		return nil, errors.New("[Auth_ERR]please Make sure the Hash and the URI's match")
	}
	return the_hash, nil
}
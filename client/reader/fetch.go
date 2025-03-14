package reader
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)
var s1 = "http://localhost:3000/get"
var s2 = "http://localhost:3001/get"
var s3 = "http://localhost:3002/get"
var s4 = "http://localhost:3003/get"
func FetchFile(hash Hash, baseDir string) {
	// url := "http://localhost:3002/get"
	var url string
	rand := hash.Rand
	Node_offset := string(rand[0])
	if(Node_offset >= "0" && Node_offset <= "3"){
		url = s1
	}
	if(Node_offset >= "4" && Node_offset <= "7"){
		url = s2
	}
	if(Node_offset >= "8" && Node_offset <= "b"){
		url = s3
	}
	if(Node_offset >= "c" && Node_offset <= "f"){
		url = s4
	}
	reqBody, err := json.Marshal(map[string]string{"hash": rand})
	if err != nil {
		fmt.Println("Error creating request body:", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error fetching file:", err)
		return
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	if contentType == "application/json" {
		var rawData bytes.Buffer
		_, err := io.Copy(&rawData, resp.Body)
		if err != nil {
			fmt.Println("Error reading JSON response:", err)
			return
		}
		var jsonResponse map[string]string
		err = json.Unmarshal(rawData.Bytes(), &jsonResponse)
		if err != nil {
			fmt.Println("Failed to parse JSON error response")
			return
		}
		fmt.Println(jsonResponse["Message"])
		return
	}
	SaveFile(hash, resp.Body, baseDir)
}

func SaveFile(hash Hash, stream io.Reader, baseDir string) {
	filePath := filepath.Join(baseDir, hash.Filename+hash.Ext)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file, ensure that BASE_DIR Already Exists:", err)
		return
	}
	defer file.Close()
	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Printf("File saved as: %s\n", filePath)
}
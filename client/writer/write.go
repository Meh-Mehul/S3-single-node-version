// Main file which sends the write file to appropriate storage service
// NOTE: In current (shitty) implementation, i have exposed each of the storages in BSS, although i can change it so that there is a custom Alloc-service in between which only reads first byte and pipes the rest to appropriate system, this Alloc-service itself will then be Load-balanced

package writer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"os"
	// "path/filepath"
)
// These SHOULD be the ports open for BSS
var s1 = "http://localhost:5000/stream"
var s2 = "http://localhost:5001/stream"
var s3 = "http://localhost:5002/stream"
var s4 = "http://localhost:5003/stream"
// This should be open for config-db
var s5 = "http://localhost:6000/"

type StreamResponse struct{
	Status string `json:"status"`
	Msg string 		`json:"msg"`
}

func StreamFile(filePath string, hash string) error {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()
	nameBuffer := make([]byte, 32)
	copy(nameBuffer, hash)
	payload := append(nameBuffer, fileContent...)

	// url := "http://localhost:5002/stream"
	// Super Primitive (and Not-recommended) Allocation logic
	var url string
	Node_offset := string(hash[0])
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
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize+32))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("upload failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// fmt.Println("Response:", string(body))
	var respjson StreamResponse
	err1 := json.Unmarshal(body, &respjson)
	if err1 != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// fmt.Println("File Successfully Written to Storage: ", respjson.Msg[len(respjson.Msg)-32:])
	return nil
}

func SendFileMetaData(file File) error {
	url := s5 + "hash"
	jsonData, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-OK status: %d", resp.StatusCode)
	}
	// fmt.Println("File metadata sent successfully for File: ", file.Name + file.Ext)
	return nil
}


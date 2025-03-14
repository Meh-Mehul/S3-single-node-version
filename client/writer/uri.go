package writer

import(
	"errors"
	"path/filepath"
	"os"
	"log"
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"sync"
)

type URI struct{
	Uri string `json:"uri"`
	Hashes []string `json:"hashes"`
}


// btw,uri is also a 16 byte random hash
func DoDir(dirpath string) error {
	res, err := CheckDir(dirpath)
	if err != nil {
		return err
	}
	if res != "directory" {
		return errors.New("this should never have been called")
	}
	uri := GetRandomHash()
	var hashes []string
	var wg sync.WaitGroup
	errChan := make(chan error, 10)
	err_ := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error accessing path:", err)
			return nil
		}
		if !info.IsDir() {
			file, err2 := GenerateFileMetaData(path)
			if err2 != nil {
				log.Println("Error Accessing file: ", err2)
				return nil
			}
			wg.Add(2)
			go func() {
				defer wg.Done()
				if err := StreamFile(file.Path, file.Rand); err != nil {
					log.Println("Error Streaming file: ", err, "\nYou might want to consider uploading files upto only 1GB, Please?")
					errChan <- err
				}
			}()
			go func() {
				defer wg.Done()
				if err := SendFileMetaData(file); err != nil {
					log.Println("Error Sending Metadata: ", err)
					errChan <- err
				}
			}()
			hashes = append(hashes, file.Rand)
			fmt.Println("Successfully Uploaded:", file.Name+file.Ext)
		}
		return nil
	})
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	if err_ != nil {
		log.Fatal("Error walking the directory:", err)
	}
	resuri := URI{
		Uri:    uri,
		Hashes: hashes,
	}
	fmt.Println("Successfully Uploaded Entire Bucket. Here is Your URI:", uri, "Please Keep it Safe!")
	if err := SendURIMetaData(resuri); err != nil {
		fmt.Println("Error sending URI Data to services: ", err)
		return err
	}
	return nil
}

func DoFile(filePath string) error {
	res, err := CheckDir(filePath)
	if err != nil {
		return err
	}
	if res != "file" {
		return errors.New("this should never have been called")
	}
	uri := GetRandomHash()
	var hashes []string
	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	file, err := GenerateFileMetaData(filePath)
	if err != nil {
		log.Println("Error generating metadata:", err)
		return err
	}
	// err = StreamFile(file.Path, file.Rand)
	// if err != nil {
	// 	log.Println("Error streaming file:", err)
	// 	return err
	// }
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := SendFileMetaData(file); err != nil {
			log.Println("Error sending metadata:", err)
			errChan <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := StreamFile(file.Path, file.Rand); err != nil {
			log.Println("Error streaming file:", err)
			errChan <- err
		}
	}()
	// err = SendFileMetaData(file)
	// if err != nil {
	// 	log.Println("Error sending metadata:", err)
	// 	return err
	// }
	wg.Wait()
	close(errChan) 
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	hashes = append(hashes, file.Rand)
	fmt.Println("Successfully uploaded:", file.Name+file.Ext)
	// // var resuri URI
	resuri := URI{
		Uri: uri,
		Hashes: hashes,
	}
	err = SendURIMetaData(resuri)
	if(err != nil){
		fmt.Println("Error sending URI Data to services: ", err)
	}
	// fmt.Println(resuri)
	fmt.Println("Successfully Uploaded Entire Directory. Here is Your URI:", uri)
	return nil
}


func SendURIMetaData(resuri URI) error {
	url := "http://localhost:6000/uri" // yes, i hardcoded supposed port AGAIN!!!
	jsonData, err := json.Marshal(resuri)
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
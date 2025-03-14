// This is the code for the CLI-tool for interacting with the S3-BSS and IW layers
//	NO Auth-system implemented yet, i think that will have to be in inter-layer ig
// Steps it needs to do for writing:
//		1. Check whether given target is folder or not?
// 		2. Traverse the directory
//		3. for each of the files, get filename, ext, path
//		4. Get a Hash, send to config-DB-Service
//		|-----> Parallely stream data to BSS
//		|-----> Parallely log stuff

// Steps it needs to do for reading (from cli-client):
// 		1. Get URI from arg, get HASH from Arg
//		2. Check if Hash is actually for that URI
//		3. Go to Hash's properties and extraxt filename, ext, path
//		4. Write(append) these in a metadata.yaml file (I have not done this yet)
//		5. Read file data from IWL.

// Steps it needs to do for Status Check:
//  1. Get URI
//  2. Get Hash data for each hash
//  3. Write to metadata.yaml , the details for each filename, ext, and path
package main

import (
	"fmt"
	"os"
	"github.com/Meh-Mehul/client/writer"
	"github.com/Meh-Mehul/client/reader"
)


func main(){
	if len(os.Args) < 2 {
		fmt.Println("Error: Please specify a command (e.g., write, read, status).")
		os.Exit(1)
	}
	QueryType := os.Args[1]
	switch QueryType{
	case "write":
		if len(os.Args) < 3 {
			fmt.Println("CALLING_ERR: Please provide the (realtive) path of directory/file you wanna upload.")
			return
		}
		path := os.Args[2]
		chk , err:= writer.CheckDir(path)
		if(err !=nil){
			fmt.Println("Please Check file path")
			return 
		}
		var res error
		if(chk == "file"){
			res = writer.DoFile(path)
		}
		if(chk == "directory"){
			res = writer.DoDir(path)
		}
		if(res != nil){
			fmt.Println("Error Uploading Please re-try or check for errors :(", res)
		}
	case "read":
		if len(os.Args) < 5 {
			fmt.Println("CALLING_ERR: Please provide the URI and the HASH and the BASE_DIR in the form:\nmain.exe read <URI> <HASH> <BASE_DIR>")
			return
		}
		uri := os.Args[2]
		hash := os.Args[3]
		base_dir := os.Args[4]
		the_hash , err:= reader.CheckContainability(uri, hash)
		if(err != nil){
			fmt.Println("Error Occured (Could be due to invalid Hash and/or URI):", err)
		}
		reader.FetchFile(*the_hash, base_dir)
	case "status":
		if len(os.Args) < 3 {
			fmt.Println("CALLING_ERR: Please provide the URI also")
			return
		}
		uri := os.Args[2]
		the_uri, err := reader.GetURIDeatilsFromDB(uri)
		if(err != nil){
			fmt.Println("Error Occured (Could be due to invalid URI):", err)
		}
		hashes := the_uri.Hashes
		var frmt_hash []reader.Hash
		for _,h := range hashes{
			frmt, err := reader.GetHashDetailsFromDB(h)
			if(err != nil){
				fmt.Println("Error Occured (Could be due to invalid HASH):", err)
			}
			frmt_hash = append(frmt_hash, *frmt)
		}
		fmt.Println("Bulk-Upload Status for URI:", uri)
		for _,f := range frmt_hash{
			print("Filename:", f.Filename+f.Ext,"|-> Hash:", f.Rand, "\n")
		}

	default:
		fmt.Println("Please provide one of write, read, status")
		return
	}
}

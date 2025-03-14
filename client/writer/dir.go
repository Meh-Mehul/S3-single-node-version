// This is having functions related to directory checking and traversal functions
package writer

import (
	"os"
	"fmt"
	"path/filepath"
)

func CheckDir(relPath string) (string, error){
	basePath, err := os.Executable();
	if err != nil {
		fmt.Println("Some Error Occurred");
		return "", err
	}
	baseDir := filepath.Dir(basePath)
	absolutePath := filepath.Join(baseDir, relPath)
	info, err := os.Stat(absolutePath)
	if os.IsNotExist(err) {
		return "does not exist", nil
	} else if err != nil {
		return "", fmt.Errorf("error checking path: %v", err)
	}
	if info.IsDir() {
		return "directory", nil
	}
	return "file", nil
}


type File struct {
	Name string `json:"filename"`
	Ext  string `json:"ext"`
	Rand string `json:"rand"`
	Path string `json:"path"`
}


func GenerateFileMetaData(filePath string) (File, error) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return File{}, err
	}
	
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return File{}, err
	}
	name := fileInfo.Name()
	ext := filepath.Ext(name)
	nameWithoutExt := name[:len(name)-len(ext)]
	hash:= GetRandomHash()

	return File{
		Name: nameWithoutExt,
		Ext:  ext,
		Path: absPath,
		Rand: hash,
	}, nil
}





package utils

import (
	"github.com/pingcap/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"sync"
)

var DIRPATHNOTDIRERRIR = errors.New("the path is not dir")

var log =  zap.L().With(zap.String("util", "file"))

func CheckDirExist(path string) (bool,error){
	s,err:=os.Stat(path)
	if err!=nil{
		log.Info("Check dir exist fail , " + err.Error())
		return false,err
	}
	ok :=s.IsDir()
	if !ok{
		log.Info("Check dir exist fail , " + path + "is not dir")
		return ok,DIRPATHNOTDIRERRIR
	}
	return ok,nil
}


func GetDataFile(filePath string,files map[string]int,mu *sync.Mutex) error {

	err := GetFilesFromPath(filePath,files,mu)
	if err!=nil{
		return err
	}
	return nil
}

func GetFilesFromPath(filePath string,files map[string]int,mu *sync.Mutex) error {

	fs, err := ioutil.ReadDir(filePath)
	if err!=nil{
		return err
	}

	for _, file := range fs {
		if file.IsDir() {
			continue
		} else {
			mu.Lock()
			files[file.Name()]=0
			mu.Unlock()
		}
	}

	return nil
}







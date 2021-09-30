package utils

import (
	"github.com/pingcap/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
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


func CloseFile(m map[string]*os.File)  error{
	for _,v := range m{
		err:=v.Close()
		if err!=nil {
			return err
		}
	}
	return nil
}

func GetDataFile(filePath string) (map[string]*os.File,error) {

	m,err := GetFilesFromPath(filePath)
	if err!=nil{
		return nil,err
	}
	for k,_ := range m{
		m[k],err =os.Open(filePath+"/"+k)
		if err!=nil{
			return nil,err
		}
	}
	return m,nil
}

func GetFilesFromPath(filePath string) (map[string]*os.File,error) {

	m :=make(map[string]*os.File)

	files, err := ioutil.ReadDir(filePath)
	if err!=nil{
		return nil ,err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			m[file.Name()]=nil
		}
	}

	return m,nil
}






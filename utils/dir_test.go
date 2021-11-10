package utils

import (
	"github.com/agiledragon/gomonkey"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/fs"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

var logger *zap.Logger


func init (){
	cfg := zap.NewDevelopmentConfig()
	//cfg.Level = zap.NewAtomicLevelAt()
	cfg.DisableStacktrace = !cfg.Level.Enabled(zap.DebugLevel)
	logger, _ = cfg.Build()
	zap.ReplaceGlobals(logger)
	logger = zap.L().With(zap.String("conn","test-mysql.go"))
	logger = logger.Named("test")
}


//test check dir exist success
func TestUtil_CheckDirExist_Succ(t *testing.T) {
	path :="./"

	ok,err:=CheckDirExist(path)
	ast:=assert.New(t)
	ast.Equal(ok,true)
	ast.Nil(err)
}

//test check dir exist success
func TestUtil_CheckDirExist_Stat_Fail(t *testing.T) {
	path :=""
	patch := gomonkey.ApplyFunc(os.Stat, func (name string) (os.FileInfo, error){
		return nil,DIRPATHNOTDIRERRIR
	})
	defer patch.Reset()
	ok,err:=CheckDirExist(path)
	ast:=assert.New(t)
	ast.Equal(ok,false)
	ast.Equal(err,DIRPATHNOTDIRERRIR)
}


func TestUtil_CheckDirExist_IsDir_Fail(t *testing.T) {
	path :="./dir_test.go"
	ok,err:=CheckDirExist(path)
	ast:=assert.New(t)
	ast.Equal(ok,false)
	ast.Equal(err,DIRPATHNOTDIRERRIR)
}

func TestUitl_GetFilesFromPath_With_ReadDir_Fail(t *testing.T){
	var filePath ="./"
	err :=errors.New("dir is not exist")
	patch := gomonkey.ApplyFunc(ioutil.ReadDir, func (dirname string) ([]fs.FileInfo, error){
		return nil,err
	})
	defer patch.Reset()

	m := make(map[string]int)
	mu := new(sync.Mutex)
	err1:=GetFilesFromPath(filePath,m,mu)

	ast:=assert.New(t)

	//ast.Nil(m)
	ast.Equal(err,err1)

}

func TestUitl_GetFilesFromPath_succ(t *testing.T){
	var filePath ="./"
	m := make(map[string]int )
	mu := new(sync.Mutex)
	err:= GetFilesFromPath(filePath,m,mu)

	ast:=assert.New(t)

	ast.Nil(err)
}




func TestGetDataFile_With_GetFilesFromPath_Fail(t *testing.T) {
	filePath:="./"
	err :=errors.New("dir is not exist")
	patch := gomonkey.ApplyFunc(GetFilesFromPath, func (filePath string,files map[string]int,mu *sync.Mutex) error{
		return err
	})
	defer patch.Reset()
	m := make(map[string]int)
	mu :=new(sync.Mutex)
	err1:= GetDataFile(filePath,m,mu)

	assert.New(t).Equal(err,err1)
}


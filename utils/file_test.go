/*******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 ******************************************************************************/

/**
 * @Author: guobob
 * @Description:
 * @File:  file_test.go
 * @Version: 1.0.0
 * @Date: 2021/11/9 20:13
 */

package utils

import (
	"github.com/agiledragon/gomonkey"
	"github.com/pingcap/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCloseFile(t *testing.T) {
	file ,_ :=os.Open("./file_test.go")
	type args struct {
		f *os.File
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "close file success ",
			args: args{
				f: file,
			},
			wantErr: false,
		},
		{
			name: "close file fail ",
			args: args{
				f: new(os.File),
			},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CloseFile(tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("CloseFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMoveFileToBackupDir_fail (t *testing.T){

	dataDir:="./"
	fileName:="test"
	backDir:="./"
	err:=errors.New("do not have privileges")
	patch := gomonkey.ApplyFunc(os.Rename, func  (oldpath  ,newpath string ) error{
		return err
	})
	defer patch.Reset()
	err1:= MoveFileToBackupDir(dataDir,fileName,backDir)
	assert.New(t).Equal(err,err1)

}

func TestMoveFileToBackupDir_succ(t *testing.T){

	dataDir:="./"
	fileName:="test"
	backDir:="./"
	patch := gomonkey.ApplyFunc(os.Rename, func  (oldpath  ,newpath string ) error{
		return nil
	})
	defer patch.Reset()
	err1:= MoveFileToBackupDir(dataDir,fileName,backDir)
	assert.New(t).Nil(err1)

}

func TestOpenFile_succ(t *testing.T){


	fileName:="test"

	patch := gomonkey.ApplyFunc(os.Open, func  (name string) (*os.File, error){
		return new(os.File),nil
	})
	defer patch.Reset()
	_,err1:= OpenFile(fileName)
	assert.New(t).Nil(err1)

}

func TestOpenFile_fail(t *testing.T){


	fileName:="test"
	err:= errors.New("do not have privileges")
	patch := gomonkey.ApplyFunc(os.Open, func  (name string) (*os.File, error){
		return nil,err
	})
	defer patch.Reset()
	f,err1:= OpenFile(fileName)
	assert.New(t).Equal(err,err1)
	assert.New(t).Nil(f)

}

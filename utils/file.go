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
 * @File:  file.go
 * @Version: 1.0.0
 * @Date: 2021/11/9 18:11
 */

package utils

import "os"

func CloseFile (f *os.File) error {
	return f.Close()
}

func MoveFileToBackupDir (dataDir  ,fileName ,backupDir string ) error {
	return os.Rename(dataDir+"/"+fileName,backupDir+"/"+fileName)
}

func OpenFile (fn string ) (*os.File,error){
	return os.Open(fn)
}

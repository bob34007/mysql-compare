/**
 * @Author: guobob
 * @Description:
 * @File:  utils.go
 * @Version: 1.0.0
 * @Date: 2021/11/23 15:07
 */

package cmd

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

func printTime(log *zap.Logger) {
	t := time.NewTicker(time.Second * 60)
	ts := time.Now()
	for {
		select {
		case <-t.C:
			fmt.Println(time.Now().String() + ":" + fmt.Sprintf("program run %v seconds", time.Since(ts).Seconds()))
		default:
			time.Sleep(time.Second * 5)
		}

	}

}

//get the file with the smallest filename order from the container
func getFirstFileName(files map[string]int) string {
	fileName:=""
	for k,v :=range files{
		if v != 0 {
			continue
		}
		if len(fileName) ==0 {
			fileName = k
		}
		if k < fileName{
			fileName = k
		}
	}
	if len(fileName) > 0 {
		files[fileName] = 1
	}
	return fileName
}
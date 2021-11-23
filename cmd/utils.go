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

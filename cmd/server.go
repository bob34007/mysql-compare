/**
 * @Author: guobob
 * @Description:
 * @File:  server.go
 * @Version: 1.0.0
 * @Date: 2021/11/12 09:12
 */

package cmd

import (
	"fmt"
	"github.com/bobguo/mysql-compare/stat"
	"go.uber.org/zap"
	"net/http"
)

/*
import "github.com/spf13/cobra"
func NewServerCommand() *cobra.Command {

	var (
		dataDir       string
		backDir       string
		maxGoroutines int32
		runTime       uint32
	)
	cmd := &cobra.Command{
		Use:   "",
		Short: "Compare the result sets of packet and replay Server",
		RunE: func(cmd *cobra.Command, args []string) error {


			return nil
		},
	}
	return cmd
}
*/

var exitChannel chan int
var logger = zap.L().Named("server")

var exitStatus = false

func generateListenStr(port uint16) string {

	return "0.0.0.0"+":"+fmt.Sprintf("%v",port)

}

func HandleExit(w http.ResponseWriter, r *http.Request) {

	logger.Info("request exit from " + r.Host )
	defer logger.Info("response exit to " + r.Host )
	if !exitStatus {
		exitStatus=true
		exitChannel <- 1
		_,err:=w.Write([]byte("ok!"))
		if err !=nil {
			logger.Warn("write response file," + err.Error())
		}
	} else {
		_,err:=w.Write([]byte("program is exited!"))
		if err !=nil {
			logger.Warn("write response file," + err.Error())
		}
	}

}



func HandleQueryStats(w http.ResponseWriter, r *http.Request) {

	logger.Info("request query stats from " + r.Host )
	defer logger.Info("response query stats to " + r.Host)
	err,str:=stat.QueryMapStr()
	if err !=nil{
		_,err=w.Write([]byte(err.Error()))
		if err !=nil {
			logger.Warn("write response file," + err.Error())
		}
	}
	_,err =w.Write([]byte(str))
	if err !=nil {
		logger.Warn("write response file," + err.Error())
	}

}

func AddPortListenAndServer(port uint16,c chan int){

	exitChannel =c

	http.HandleFunc("/stats", HandleQueryStats)
	http.HandleFunc("/exit", HandleExit)


	err := http.ListenAndServe(generateListenStr(port), nil)
	if err !=nil{
		logger.Warn(fmt.Sprintf("listen port:%v fail ,%v",port,err.Error()))
	}

}
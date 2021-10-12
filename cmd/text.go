package cmd

import (
	"context"
	"github.com/bobguo/mysql-compare/stat"
	"github.com/bobguo/mysql-compare/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)



func NewTextCompareCommand() *cobra.Command {

	var (
		dataDir string
	)
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare the result sets of packet and replay Server",
		RunE: func(cmd *cobra.Command, args []string) error {

			log:= zap.L().Named("compare")
			log.Info("process begin run at " + time.Now().String())
			ok,err:=utils.CheckDirExist(dataDir)
			if !ok{
				log.Error("param dataDir error , " + err.Error())
				return nil
			}
			m,err := utils.GetDataFile(dataDir)
			if err!=nil{
				log.Error("get file from dataDir fail , "+ err.Error())
				return nil
			}

			var wg sync.WaitGroup
			wg.Add(len(m))

			for k,v := range m {
				go func (k string ,v *os.File) {
					utils.DoCompare(k,v,&wg)
				}(k,v)
			}
			ctx ,cancel := context.WithCancel(context.Background())
			go stat.Statis.PrintStaticWithTimer(ctx,log)

			wg.Wait()
			cancel()
			defer utils.CloseFile(m)
			//wait 200ms before exit for goruntine done
			log.Info("process end run at " + time.Now().String())
			//time.Sleep(200 * time.Millisecond)
			<-time.After(200 * time.Millisecond)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dataDir,"data-dir","d","","directory used to read the result set")
	return cmd
}





func NewTextCommand() *cobra.Command {
	//add sub command test
	cmd := &cobra.Command{
		Use:   "text",
		Short: "Text format utilities",
	}
	cmd.AddCommand(NewTextCompareCommand())
	return cmd
}

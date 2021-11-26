package cmd

import (
	"context"
	"github.com/bobguo/mysql-compare/compare"
	"github.com/bobguo/mysql-compare/stat"
	"github.com/bobguo/mysql-compare/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

func NewTextCompareCommand() *cobra.Command {
	var (
		cfg = &utils.Config{}
	)
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare the result sets of packet and replay Server",
		RunE: func(cmd *cobra.Command, args []string) error {

			log := zap.L().Named("compare")
			log.Info("process begin run at " + time.Now().String())
			err := cfg.CheckConfig()
			if err !=nil {
				log.Error("param dataDir error , " + err.Error())
				return err
			}

			go printTime(log)

			ts := time.Now()
			t := time.NewTicker(3 * time.Second)
			defer t.Stop()

			mu := new(sync.Mutex)
			var ctGorountines int32 = 0
			files := make(map[string]int, 0)
			err = utils.GetDataFile(cfg.DataDir, files, mu)
			if err != nil {
				log.Error("get file from dataDir fail , " + err.Error())
				return nil
			}

			compare.BasePercent = uint64(cfg.BasePercent)
			//use for save server channel
			c := make(chan int, 1)
			go AddPortListenAndServer(cfg.ListenPort, c)

			wg := new(sync.WaitGroup)
			ctx, cancel := context.WithCancel(context.Background())
			go stat.Statis.PrintStaticWithTimer(ctx, log)
			go utils.WatchDirCreateFile(ctx, cfg.DataDir, files, mu, log)

			var fileName = ""
			var exit = false
			for {
				if atomic.LoadInt32(&ctGorountines) < cfg.MaxGoroutines {
					mu.Lock()
					fileName = getFirstFileName(files)
					if len(fileName) <=0{
						mu.Unlock()
						goto LOOP
					}
					atomic.AddInt32(&ctGorountines, 1)
					wg.Add(1)
					go func(fileName string) {
						compare.DoCompare(fileName, &ctGorountines, wg, cfg.DataDir, cfg.BackDir)
					}(fileName)
					if atomic.LoadInt32(&ctGorountines) >= cfg.MaxGoroutines {
						mu.Unlock()
						goto LOOP
					}
					mu.Unlock()
				}
LOOP:
				select {
				case <-t.C:
					if time.Now().Sub(ts).Seconds() > float64(cfg.RunTime*60) {
						exit = true
					}
				case <-c:
					exit = true
				default:
					//
				}
				if exit {
					break
				}
			}

			cancel()
			wg.Wait()
			if err = stat.PrintMap(log); err != nil {
				log.Error(err.Error())
			}

			//wait 200ms before exit for goruntine done
			log.Info("process end run at " + time.Now().String())

			<-time.After(200 * time.Millisecond)

			close(c)
			return nil
		},
	}

	cmd.Flags().StringVarP(&cfg.DataDir, "data-dir", "d", "", "directory used to read the result set")
	cmd.Flags().StringVarP(&cfg.BackDir, "back-dir", "b", "", "directory used to back up the result file")
	cmd.Flags().Int32VarP(&cfg.MaxGoroutines, "max-routines", "g", 10, "max goroutines to parse result files")
	cmd.Flags().Uint32VarP(&cfg.RunTime, "runtime", "t", 10,
		"program runtime, if zero is specified then all files in the current directory will be processed")
	cmd.Flags().Uint16VarP(&cfg.ListenPort, "listen-port", "P", 7001, "http server port , Provide query statistical (query) information and exit (exit) services")
	cmd.Flags().Uint16VarP(&cfg.BasePercent, "base-percent", "B", 100, "SQL execution time deterioration statistics benchmark %")

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

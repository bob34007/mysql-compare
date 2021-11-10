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
		dataDir       string
		backDir       string
		maxGoroutines int32
		runTime       uint32
	)
	cmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare the result sets of packet and replay Server",
		RunE: func(cmd *cobra.Command, args []string) error {

			log := zap.L().Named("compare")
			log.Info("process begin run at " + time.Now().String())
			ok, err := utils.CheckDirExist(dataDir)
			if !ok {
				log.Error("param dataDir error , " + err.Error())
				return nil
			}
			ts := time.Now()
			t := time.NewTicker(3 * time.Second)
			defer t.Stop()

			mu := new(sync.Mutex)
			var ctGorountines int32 = 0
			files := make(map[string]int, 0)
			err = utils.GetDataFile(dataDir, files, mu)
			if err != nil {
				log.Error("get file from dataDir fail , " + err.Error())
				return nil
			}

			wg := new(sync.WaitGroup)
			ctx, cancel := context.WithCancel(context.Background())
			go stat.Statis.PrintStaticWithTimer(ctx, log)
			go utils.WatchDirCreateFile(ctx, dataDir, files, mu, log)

			var exit = false
			for {
				if atomic.LoadInt32(&ctGorountines) < maxGoroutines {
					mu.Lock()
					for k, v := range files {
						if v > 0 {
							continue
						}
						files[k] = 1
						atomic.AddInt32(&ctGorountines, 1)
						wg.Add(1)
						go func(k string) {
							compare.DoCompare(k, &ctGorountines, wg,dataDir,backDir)
						}(k)
						if atomic.LoadInt32(&ctGorountines) >= maxGoroutines {
							break
						}
					}
					mu.Unlock()
				}

				select {
				case <-t.C:
					if time.Now().Sub(ts).Seconds() > float64(runTime*60) {
						exit = true
					}
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
			return nil
		},
	}

	cmd.Flags().StringVarP(&dataDir, "data-dir", "d", "", "directory used to read the result set")
	cmd.Flags().StringVarP(&backDir, "back-dir", "b", "", "directory used to back up the result file")
	cmd.Flags().Int32VarP(&maxGoroutines, "max-routines", "g", 10, "max goroutines to parse result files")
	cmd.Flags().Uint32VarP(&runTime, "runtime", "t", 10,
		"program runtime, if zero is specified then all files in the current directory will be processed")
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

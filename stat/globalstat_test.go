package stat

import (
	"context"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

var logger *zap.Logger

func init (){
	cfg := zap.NewDevelopmentConfig()
	//cfg.Level = zap.NewAtomicLevelAt()
	cfg.DisableStacktrace = !cfg.Level.Enabled(zap.DebugLevel)
	logger, _ = cfg.Build()
	zap.ReplaceGlobals(logger)
	logger = zap.L().With(zap.String("conn","stat_test.go"))
	logger = logger.Named("test")
}



func TestStat_AddStatic_Succ(t *testing.T) {
	item := "ReadSuccFiles"
	Statis.AddStatic(item,logger)
	ast:=assert.New(t)
	ast.Equal(Statis.MStat[item],uint64(1))
}

func TestStat_AddStatic_Fail(t *testing.T) {
	item := "AAAA"
	Statis.AddStatic(item,logger)
	_,ok:=Statis.MStat[item]
	ast:=assert.New(t)
	ast.Equal(ok,false)
}

func TestStat_AddStatica(t *testing.T) {
	item := make ([]string,0)
	str1 := "ReadFailFiles"
	str2 := "AAAA"
	item = append(item,str1,str2)

	Statis.AddStatica(item,logger)

	_,ok:=Statis.MStat[str2]

	ast:=assert.New(t)
	//fmt.Println(Statis.MStat[str1])
	ast.Equal(Statis.MStat[str1],uint64(1))
	ast.Equal(ok,false)
}

func TestStat_AddStaticm(t *testing.T) {
	item := make(map[string]uint64)
	str1 := "ResultFiles"
	str2 := "AAAA"
	item[str1] =10
	item[str2] =1
	Statis.AddStaticm(item,logger)
	_,ok:=Statis.MStat[str2]

	ast:=assert.New(t)
	//fmt.Println(Statis.MStat[str1])
	ast.Equal(Statis.MStat[str1],uint64(10))
	ast.Equal(ok,false)

}

func TestStat_SetBucketNum(t *testing.T) {

	var execTime uint64 = uint64(5* time.Millisecond)

	Statis.SetBucketNum(execTime,0)
	Statis.SetBucketNum(execTime,1)

	var execTime1 uint64 = uint64(15* time.Millisecond)
	Statis.SetBucketNum(execTime1,0)
	Statis.SetBucketNum(execTime1,1)

	var execTime2 uint64 = uint64(25* time.Millisecond)
	Statis.SetBucketNum(execTime2,0)
	Statis.SetBucketNum(execTime2,1)

	var execTime3 uint64 = uint64(35* time.Millisecond)
	Statis.SetBucketNum(execTime3,0)
	Statis.SetBucketNum(execTime3,1)

	var execTime4 uint64 = uint64(45* time.Millisecond)
	Statis.SetBucketNum(execTime4,0)
	Statis.SetBucketNum(execTime4,1)

	var execTime5 uint64 = uint64(55* time.Millisecond)
	Statis.SetBucketNum(execTime5,0)
	Statis.SetBucketNum(execTime5,1)

	var execTime6 uint64 = uint64(105* time.Millisecond)
	Statis.SetBucketNum(execTime6,0)
	Statis.SetBucketNum(execTime6,1)

	ast:= assert.New(t)
	ast.Equal(Statis.MStat["PrExecTimeIn10ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeIn20ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeIn30ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeIn40ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeIn50ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeIn100ms"],uint64(1))
	ast.Equal(Statis.MStat["PrExecTimeOut100ms"],uint64(1))

	ast.Equal(Statis.MStat["RrExecTimeIn10ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeIn20ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeIn30ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeIn40ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeIn50ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeIn100ms"],uint64(1))
	ast.Equal(Statis.MStat["RrExecTimeOut100ms"],uint64(1))
}

func TestStat_PrintStaticToLog(t *testing.T) {
	Statis.PrintStaticToLog(logger)
}

func TestStat_PrintStaticToConsole(t *testing.T) {
	Statis.PrintStaticToConsole()
}

func ControlTimer(cancel context.CancelFunc){
	var i =0
	ticker1 :=time.NewTicker(100*time.Millisecond)
	for{
		select {
		case <-ticker1.C:
			if i<1 {
				i++
			} else {
				cancel()
				return
			}
		}
	}
}

func TestStat_PrintStaticWithTimer(t *testing.T) {
	ticker := time.NewTicker(10*time.Millisecond)
	patch := gomonkey.ApplyFunc(time.NewTicker, func (d time.Duration) *time.Ticker{
		return ticker
	})
	defer patch.Reset()

	ctx1,cancel := context.WithCancel(context.Background())

	go ControlTimer(cancel)

	Statis.PrintStaticWithTimer(ctx1,logger)

}
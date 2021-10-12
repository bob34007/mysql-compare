package stat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

var Statis = new(Stat)

func init() {
	m := make(map[string]uint64)
	Statis.MStat = m
	Statis.Lock()
	defer Statis.Unlock()
	m["ResultFiles"] = 0
	m["ReadSuccFiles"] = 0
	m["ReadFailFiles"] = 0
	m["ExecSqlNum"] = 0
	m["ExecSuccNum"] = 0
	m["ExecFailNum"] = 0
	m["ExecErrNoNotEqual"] = 0
	m["ExecTimeNotEqual"] = 0
	m["RowCountNotequal"] = 0
	m["RowDetailNotEqual"] = 0
	m["PrExecRowCount"] = 0
	m["PrExecSuccCount"] = 0
	m["PrExecFailCount"] = 0
	m["PrExecTimeCount"] = 0
	m["PrMaxExecTime"] = 0
	m["PrMinExecTime"] = 0
	m["PrExecTimeIn10ms"] = 0
	m["PrExecTimeIn20ms"] = 0
	m["PrExecTimeIn30ms"] = 0
	m["PrExecTimeIn40ms"] = 0
	m["PrExecTimeIn50ms"] = 0
	m["PrExecTimeIn100ms"] = 0
	m["PrExecTimeOut100ms"] = 0
	m["RrExecTimeCount"] = 0
	m["RrExecRowCount"] = 0
	m["RrExecSuccCount"] = 0
	m["RrExecFailCount"] = 0
	m["RrMaxExecTime"] = 0
	m["RrMinExecTime"] = 0
	m["RrExecTimeIn10ms"] = 0
	m["RrExecTimeIn20ms"] = 0
	m["RrExecTimeIn30ms"] = 0
	m["RrExecTimeIn40ms"] = 0
	m["RrExecTimeIn50ms"] = 0
	m["RrExecTimeIn100ms"] = 0
	m["RrExecTimeOut100ms"] = 0
}

type Stat struct {
	sync.Mutex
	MStat map[string]uint64
}

func (s *Stat) SetBucketNum(execTime uint64, serverType int8) {
	execTimeMS := execTime / uint64(time.Millisecond)
	s.Lock()
	defer s.Unlock()
	switch true {
	case execTimeMS < 10:
		if serverType == 0 {
			s.MStat["PrExecTimeIn10ms"]++
		} else {
			s.MStat["RrExecTimeIn10ms"]++
		}
	case execTimeMS >= 10 && execTimeMS < 20:
		if serverType == 0 {
			s.MStat["PrExecTimeIn20ms"]++
		} else {
			s.MStat["RrExecTimeIn20ms"]++
		}
	case execTimeMS >= 20 && execTimeMS < 30:
		if serverType == 0 {
			s.MStat["PrExecTimeIn30ms"]++
		} else {
			s.MStat["RrExecTimeIn30ms"]++
		}
	case execTimeMS >= 30 && execTimeMS < 40:
		if serverType == 0 {
			s.MStat["PrExecTimeIn40ms"]++
		} else {
			s.MStat["RrExecTimeIn40ms"]++
		}
	case execTimeMS >= 40 && execTimeMS < 50:
		if serverType == 0 {
			s.MStat["PrExecTimeIn50ms"]++
		} else {
			s.MStat["RrExecTimeIn50ms"]++
		}
	case execTimeMS >= 50 && execTimeMS < 100:
		if serverType == 0 {
			s.MStat["PrExecTimeIn100ms"]++
		} else {
			s.MStat["RrExecTimeIn100ms"]++
		}
	default:
		if serverType == 0 {
			s.MStat["PrExecTimeOut100ms"]++
		} else {
			s.MStat["RrExecTimeOut100ms"]++
		}
	}
}

func (s *Stat) AddStatic(item string, log *zap.Logger) {
	s.Lock()
	defer s.Unlock()
	_, ok := s.MStat[item]
	if ok {
		s.MStat[item]++
	} else {
		log.Warn("can not find key in map , " + item)
	}
}

func (s *Stat) AddStatica(item []string, log *zap.Logger) {
	s.Lock()
	defer s.Unlock()
	for i := range item {
		_, ok := s.MStat[item[i]]
		if ok {
			s.MStat[item[i]]++
		} else {
			log.Warn("can not find key in map , " + item[i])
		}
	}
}

func (s *Stat) AddStaticm(item map[string]uint64, log *zap.Logger) {
	s.Lock()
	defer s.Unlock()
	for i, val := range item {
		_, ok := s.MStat[i]
		if ok {
			s.MStat[i] += val
		} else {
			log.Warn("can not find key in map , " + i)
		}
	}
}

func (s *Stat) PrintStaticWithTimer(ctx context.Context, log *zap.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.PrintStaticToLog(log)
		case <-ctx.Done():
			s.PrintStaticToConsole()
			return
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (s *Stat) PrintStaticToLog(log *zap.Logger) {
	var logStr string
	var strResultFiles = "ResultFiles"
	var strReadSuccFiles = "ReadSuccFiles"
	var strReadFailFiles = "ReadFailFiles"
	var strExecSqlNum = "ExecSqlNum"
	var strExecSuccNum = "ExecSuccNum"
	var strExecFailNum = "ExecFailNum"
	var strExecErrNoNotEqual = "ExecErrNoNotEqual"
	var strExecTimeNotEqual = "ExecTimeNotEqual"
	var strRowCountNotequal = "RowCountNotequal"
	var strRowDetailNotEqual = "RowDetailNotEqual"

	s.Lock()
	logStr += "-------------timer print begin-------------\n"
	logStr += fmt.Sprintf("%v : %v \n", strResultFiles, s.MStat[strResultFiles])
	logStr += fmt.Sprintf("%v : %v \n", strReadSuccFiles, s.MStat[strResultFiles])
	logStr += fmt.Sprintf("%v : %v \n", strReadFailFiles, s.MStat[strReadFailFiles])
	logStr += fmt.Sprintf("%v : %v \n", strExecSqlNum, s.MStat[strExecSqlNum])
	logStr += fmt.Sprintf("%v : %v \n", strExecSuccNum, s.MStat[strExecSuccNum])
	logStr += fmt.Sprintf("%v : %v \n", strExecFailNum, s.MStat[strExecFailNum])
	logStr += fmt.Sprintf("%v : %v \n", strExecErrNoNotEqual, s.MStat[strExecErrNoNotEqual])
	logStr += fmt.Sprintf("%v : %v \n", strExecTimeNotEqual, s.MStat[strExecTimeNotEqual])
	logStr += fmt.Sprintf("%v : %v \n", strRowCountNotequal, s.MStat[strRowCountNotequal])
	logStr += fmt.Sprintf("%v : %v \n", strRowDetailNotEqual, s.MStat[strRowDetailNotEqual])
	logStr += "-------------timer print end-------------\n"
	s.Unlock()

	log.Info(logStr)
	return
}

func (s *Stat) PrintStaticToConsole() {
	var logStr string
	var strResultFiles = "ResultFiles"
	var strReadSuccFiles = "ReadSuccFiles"
	var strReadFailFiles = "ReadFailFiles"
	var strExecSqlNum = "ExecSqlNum"
	var strExecSuccNum = "ExecSuccNum"
	var strExecFailNum = "ExecFailNum"
	var strExecErrNoNotEqual = "ExecErrNoNotEqual"
	var strExecTimeNotEqual = "ExecTimeNotEqual"
	var strRowCountNotequal = "RowCountNotequal"
	var strRowDetailNotEqual = "RowDetailNotEqual"
	var strPrExecRowCount = "PrExecRowCount"
	var strPrExecSuccCount = "PrExecSuccCount"
	var strPrExecFailCount = "PrExecFailCount"
	var strPrExecTimeCount = "PrExecTimeCount"
	var strPrMaxExecTime = "PrMaxExecTime"
	var strPrMinExecTime = "PrMinExecTime"
	var strPrExecTimeIn10ms = "PrExecTimeIn10ms"
	var strPrExecTimeIn20ms = "PrExecTimeIn20ms"
	var strPrExecTimeIn30ms = "PrExecTimeIn30ms"
	var strPrExecTimeIn40ms = "PrExecTimeIn40ms"
	var strPrExecTimeIn50ms = "PrExecTimeIn50ms"
	var strPrExecTimeIn100ms = "PrExecTimeIn100ms"
	var strPrExecTimeOut100ms = "PrExecTimeOut100ms"
	var strRrExecTimeCount = "RrExecTimeCount"
	var strRrExecRowCount = "RrExecRowCount"
	var strRrExecSuccCount = "RrExecSuccCount"
	var strRrExecFailCount = "RrExecFailCount"
	var strRrMaxExecTime = "RrMaxExecTime"
	var strRrMinExecTime = "RrMinExecTime"
	var strRrExecTimeIn10ms = "RrExecTimeIn10ms"
	var strRrExecTimeIn20ms = "RrExecTimeIn20ms"
	var strRrExecTimeIn30ms = "RrExecTimeIn30ms"
	var strRrExecTimeIn40ms = "RrExecTimeIn40ms"
	var strRrExecTimeIn50ms = "RrExecTimeIn50ms"
	var strRrExecTimeIn100ms = "RrExecTimeIn100ms"
	var strRrExecTimeOut100ms = "RrExecTimeOut100ms"

	var ExecSuccPre uint64
	var ReadSuccPre uint64
	s.Lock()
	logStr += "-------------compare result begin-------------\n"
	logStr += fmt.Sprintf("%v : %v \n", strResultFiles, s.MStat[strResultFiles])

	if s.MStat[strResultFiles] >0 {
		ReadSuccPre = uint64(float64( s.MStat[strReadSuccFiles]) / float64(s.MStat[strResultFiles]) * 100)
	} else {
		ReadSuccPre =100
	}
	logStr += fmt.Sprintf("%v : %v %v%% \n", strReadSuccFiles, s.MStat[strReadSuccFiles],ReadSuccPre)
	logStr += fmt.Sprintf("%v : %v \n", strReadFailFiles, s.MStat[strReadFailFiles])
	logStr += fmt.Sprintf("%v : %v \n", strExecSqlNum, s.MStat[strExecSqlNum])
	if s.MStat[strExecSqlNum] > 0 {
		ExecSuccPre = uint64(float64(s.MStat[strExecSuccNum]) / float64(s.MStat[strExecSqlNum]) * 100)
	} else {
		ExecSuccPre = 100
	}
	logStr += fmt.Sprintf("%v : %v %v%% \n", strExecSuccNum, s.MStat[strExecSuccNum], ExecSuccPre)
	logStr += fmt.Sprintf("%v : %v \n", strExecFailNum, s.MStat[strExecFailNum])
	logStr += fmt.Sprintf("%v : %v \n", strExecErrNoNotEqual, s.MStat[strExecErrNoNotEqual])
	logStr += fmt.Sprintf("%v : %v \n", strExecTimeNotEqual, s.MStat[strExecTimeNotEqual])
	logStr += fmt.Sprintf("%v : %v \n", strRowCountNotequal, s.MStat[strRowCountNotequal])
	logStr += fmt.Sprintf("%v : %v \n", strRowDetailNotEqual, s.MStat[strRowDetailNotEqual])
	logStr += "-------------result from packet-------------\n"
	logStr += fmt.Sprintf("%v : %v \n", strPrExecRowCount, s.MStat[strPrExecRowCount])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecSuccCount, s.MStat[strPrExecSuccCount])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecFailCount, s.MStat[strPrExecFailCount])
	logStr += fmt.Sprintf("%v(us) : %v \n", strPrExecTimeCount, s.MStat[strPrExecTimeCount])
	logStr += fmt.Sprintf("%v : %v \n", strPrMaxExecTime, s.MStat[strPrMaxExecTime])
	logStr += fmt.Sprintf("%v : %v \n", strPrMinExecTime, s.MStat[strPrMinExecTime])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn10ms, s.MStat[strPrExecTimeIn10ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn20ms, s.MStat[strPrExecTimeIn20ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn30ms, s.MStat[strPrExecTimeIn30ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn40ms, s.MStat[strPrExecTimeIn40ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn50ms, s.MStat[strPrExecTimeIn50ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeIn100ms, s.MStat[strPrExecTimeIn100ms])
	logStr += fmt.Sprintf("%v : %v \n", strPrExecTimeOut100ms, s.MStat[strPrExecTimeOut100ms])
	logStr += "-------------result from replay server-------------\n"
	logStr += fmt.Sprintf("%v : %v \n", strRrExecRowCount, s.MStat[strRrExecRowCount])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecSuccCount, s.MStat[strRrExecSuccCount])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecFailCount, s.MStat[strRrExecFailCount])
	logStr += fmt.Sprintf("%v(us) : %v\n", strRrExecTimeCount, s.MStat[strRrExecTimeCount])
	logStr += fmt.Sprintf("%v : %v \n", strRrMaxExecTime, s.MStat[strRrMaxExecTime])
	logStr += fmt.Sprintf("%v : %v \n", strRrMinExecTime, s.MStat[strRrMinExecTime])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn10ms, s.MStat[strRrExecTimeIn10ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn20ms, s.MStat[strRrExecTimeIn20ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn30ms, s.MStat[strRrExecTimeIn30ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn40ms, s.MStat[strRrExecTimeIn40ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn50ms, s.MStat[strRrExecTimeIn50ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeIn100ms, s.MStat[strRrExecTimeIn100ms])
	logStr += fmt.Sprintf("%v : %v \n", strRrExecTimeOut100ms, s.MStat[strRrExecTimeOut100ms])
	logStr += "-------------compare result end-------------\n"
	s.Unlock()
	fmt.Println(logStr)
	return
}

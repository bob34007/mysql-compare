package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bobguo/mysql-compare/stat"
	"github.com/pingcap/errors"
	"go.uber.org/zap"
)

type ResFromFile struct {
	Type   uint64        `json:"type"`
	StmtID uint64        `json:"stmtID,omitempty"`
	Params []interface{} `json:"params,omitempty"`
	DB     string        `json:"db,omitempty"`
	Query  string        `json:"query,omitempty"`
	//read from packet
	PrBeginTime uint64     `json:"pr-begin-time"`
	PrEndTime   uint64     `json:"pr-end-time"`
	PrErrorNo   uint16     `json:"pr-error-no"`
	PrErrorDesc string     `json:"pr-error-desc"`
	PrResult    [][]string `json:"pr-result"`
	//read from replay server
	RrBeginTime uint64     `json:"rr-begin-time"`
	RrEndTime   uint64     `json:"rr-end-time"`
	RrErrorNo   uint16     `json:"rr-error-no"`
	RrErrorDesc string     `json:"rr-error-desc"`
	RrResult    [][]string `json:"rr-result"`
	Logger      *zap.Logger
	File      *os.File
	Pos         int64
}

func NewResForWriteFile(file *os.File, log *zap.Logger) *ResFromFile {

	rs := new(ResFromFile)
	//common
	//TODO
	//rs.DB
	//rs.Type
	//rs.StmtID
	//rs.Params

	// packet result
	rs.File = file

	rs.Logger = log
	return rs
}

func (rs *ResFromFile) InitResFromFile() {
	rs.Type = 0
	rs.StmtID = 0
	rs.Params = rs.Params[:0]
	rs.DB = ""
	rs.Query = ""
	rs.PrBeginTime = 0
	rs.PrEndTime = 0
	rs.PrErrorNo = 0
	rs.PrErrorDesc = ""
	rs.PrResult = rs.PrResult[:0][:0]
	rs.RrBeginTime = 0
	rs.RrEndTime = 0
	rs.RrErrorNo = 0
	rs.RrErrorDesc = ""
	rs.RrResult = rs.RrResult[:0][:0]
}


func (rs *ResFromFile) GetResFromFile() ([]byte, error) {

	f := rs.File

	l := make([]byte,8)

	_, err := f.Read(l)
	if err != nil && err != io.EOF {
		rs.Logger.Warn("read data from file fail , " + err.Error())
		return nil, err
	}
	if err == io.EOF {
		rs.Logger.Info("read end , " + err.Error())
		return nil, err
	}

	dataLen :=  binary.BigEndian.Uint64(l)


	data := make([]byte,dataLen)

	_, err = f.Read(data)
	if err != nil  {
		rs.Logger.Warn("read data from file fail , " + err.Error())
		return nil, err
	}
	//fmt.Println(data,"------",dataLen,"-------",n)
	return data[:dataLen-1], nil

}

func (rs *ResFromFile) UnMarshalToStruct(s []byte) error {
	err := json.Unmarshal(s, rs)
	if err != nil {
		rs.Logger.Warn("Unmarshal json to struct fail , " + err.Error())
		return err
	}
	return nil
}

type SqlCompareRes struct {
	Type     string        `json:"sqltype"`
	Sql      string        `json:"sql"`
	Values   []interface{} `json:"values"`
	ErrCode  int           `json:"errcode"`
	ErrDesc  string        `json:"errdesc"`
	RrValues [][]string    `json:"rrValues"`
	PrValues [][]string    `json:"prValues"`
}

const (
	EventHandshake uint64 = iota
	EventQuit
	EventQuery
	EventStmtPrepare
	EventStmtExecute
	EventStmtClose
)

func (rs *ResFromFile) TypeString() string {
	switch rs.Type {
	case EventHandshake:
		return "EventHandshake"
	case EventQuit:
		return "EventQuit"
	case EventQuery:
		return "EventQuery"
	case EventStmtPrepare:
		return "EventStmtPrepare"
	case EventStmtExecute:
		return "EventStmtExecute"
	case EventStmtClose:
		return "EventStmtClose"
	default:
		return "UNKnownType"
	}
}

func (rs *ResFromFile) CompareRes() *SqlCompareRes {

	res := new(SqlCompareRes)
	res.Type = rs.TypeString()
	res.Sql = rs.Query
	res.Values = rs.Params

	m := make(map[string]uint64)

	defer stat.Statis.AddStaticm(m, rs.Logger)
	m["ExecSqlNum"] = 1

	// Calculate  SQL execution time in packet
	var prSqlExecTime uint64
	if rs.PrBeginTime < rs.PrEndTime {
		prSqlExecTime = rs.PrEndTime - rs.PrBeginTime
	} else {
		prSqlExecTime = 0
	}
	m["PrExecTimeCount"] = prSqlExecTime

	//Calculate  SQL execution time in replay server
	var rrSqlExecTime uint64
	if rs.RrBeginTime < rs.RrEndTime {
		rrSqlExecTime = rs.RrEndTime - rs.RrBeginTime
	} else {
		rrSqlExecTime = 0
	}
	m["RrExecTimeCount"] = rrSqlExecTime

	stat.Statis.SetBucketNum(rrSqlExecTime, 1)
	stat.Statis.SetBucketNum(prSqlExecTime, 0)

	var prlen = 0
	var rrlen = 0

	prlen = len(rs.PrResult)
	rrlen = len(rs.RrResult)

	m["PrExecRowCount"] = uint64(prlen)
	m["RrExecRowCount"] = uint64(rrlen)

	if rs.PrErrorNo != 0 {
		m["PrExecFailCount"] = 1
	} else {
		m["PrExecSuccCount"] = 1
	}

	if rs.RrErrorNo != 0 {
		m["RrExecFailCount"] = 1
	} else {
		m["RrExecSuccCount"] = 1
	}
	//compare errcode
	if rs.PrErrorNo != rs.RrErrorNo {
		res.ErrCode = 1
		res.ErrDesc = fmt.Sprintf("%v-%v", rs.PrErrorNo, rs.RrErrorNo)
		m["ExecErrNoNotEqual"] = 1
		m["ExecFailNum"] = 1
		return res
	}

	//compare exec time

/*
	if rrSqlExecTime > (10*prSqlExecTime) &&
		rrSqlExecTime > 150*1000000 {
		//From http://en.wikipedia.org/wiki/Order_of_magnitude: "We say two
		//numbers have the same order of magnitude of a number if the big
		//one divided by the little one is less than 10. For example, 23 and
		//82 have the same order of magnitude, but 23 and 820 do not."
		res.ErrCode = 2
		res.ErrDesc = fmt.Sprintf("%v us-%v us",
			prSqlExecTime,
			rrSqlExecTime)
		m["ExecFailNum"] = 1
		m["ExecTimeNotEqual"] = 1

		return res
	}
*/

	//compare  result row num

	if prlen != rrlen {
		res.ErrCode = 3
		res.ErrDesc = fmt.Sprintf("%v-%v", prlen, rrlen)
		m["ExecFailNum"] = 1
		m["RowCountNotequal"] = 1
		res.RrValues = rs.RrResult
		res.PrValues = rs.PrResult
		return res
	} else if prlen == 0 {
		res.ErrCode = 0
		m["ExecSuccNum"] = 1
		return res
	}

	//compare result row detail
	ok, err := rs.CompareResDetail()
	if !ok {
		rs.Logger.Error("Compare result data fail , " + err.Error())
		m["ExecFailNum"] = 1
		m["RowDetailNotEqual"] = 1
		res.ErrCode = 4
		res.RrValues = rs.RrResult
		res.PrValues = rs.PrResult
		return res
	}
	res.ErrCode = 0
	m["ExecSuccNum"] = 1

	return res
}

func (rs *ResFromFile) HashPrResDetail() ([][]interface{}, error) {
	//var rowStr string
	res := make([][]interface{}, 0)
	for i := range rs.PrResult {
		rowStr := strings.Join(rs.PrResult[i], "")
		rowv := make([]interface{}, 0)
		v, err := hashString(rowStr)
		if err != nil {
			return nil, err
		}
		rowv = append(rowv, v, uint64(i))
		res = append(res, rowv)
	}
	return res, nil
}

func (rs *ResFromFile) HashRrResDetail() ([][]interface{}, error) {
	//var rowStr string
	res := make([][]interface{}, 0)
	for i := range rs.RrResult {
		rowStr := strings.Join(rs.RrResult[i], "")
		rowv := make([]interface{}, 0)
		v, err := hashString(rowStr)
		if err != nil {
			return nil, err
		}
		rowv = append(rowv, v, uint64(i))
		res = append(res, rowv)
	}
	return res, nil
}

func (rs *ResFromFile) CompareData(a, b [][]interface{}) error {
	for i := range a {
		ok := CompareInterface(a[i][0], b[i][0])
		if !ok {
			err := errors.New("compare row string hash value not equal ")
			return err
		}
	}
	return nil
}

func (rs *ResFromFile) CompareResDetail() (bool, error) {

	a, err := rs.HashPrResDetail()
	if err != nil {
		rs.Logger.Error(" hash Packet result fail , " + err.Error())
		return false, err
	}
	Sort2DSlice(a)
	b, err := rs.HashRrResDetail()
	if err != nil {
		rs.Logger.Error(" hash Replay server  result fail , " + err.Error())
		return false, err
	}
	Sort2DSlice(b)

	err = rs.CompareData(a, b)
	if err != nil {
		rs.Logger.Warn("compare data fail , " + err.Error())
		return false, err
	}

	return true, nil
}

//read result from file ,and compare packet result and replay server result
func DoCompare(fn string, f *os.File, wg *sync.WaitGroup) {
	defer wg.Done()
	logger := fileName(fn).Logger()
	stat.Statis.AddStatic("ResultFiles", logger)


	rs := NewResForWriteFile(f, logger)
	for {
		rs.InitResFromFile()
		s, err := rs.GetResFromFile()
		if err != nil && err != io.EOF {
			stat.Statis.AddStatic("ReadFailFiles", logger)
			logger.Error("read result file fail , " + err.Error())
			return
		}
		if err == io.EOF {
			break
		}

		err = rs.UnMarshalToStruct(s)
		if err != nil {
			stat.Statis.AddStatic("ReadFailFiles", logger)
			logger.Error("UnMarshal read string to struct , " + err.Error()+string(s))
			return
		}

		res := rs.CompareRes()
		if res.ErrCode != 0 {
			jsons, err := json.Marshal(res)
			if err != nil {
				logger.Warn("Marshal compare result to json fail , " + err.Error())
			} else {
				logger.Warn(string(String(jsons)))
			}
		}

	}
	logger.Info("read data file success")
	stat.Statis.AddStatic("ReadSuccFiles", logger)
	return
}

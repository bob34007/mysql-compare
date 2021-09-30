package utils

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/bobguo/mysql-compare/stat"
	"github.com/pingcap/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"sync"
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
	Reader      *bufio.Reader
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
	rs.Reader = bufio.NewReader(file)

	rs.Logger = log
	return rs
}

func (rs *ResFromFile) CheckDataValid(s []byte) ([]byte, error) {

	l := binary.BigEndian.Uint64(s)
	s = s[8:]
	if l != uint64(len(s)) {
		logstr := fmt.Sprintf(" read res is invalid ,len not qeual , %v-%v ", l, len(s))
		rs.Logger.Warn(logstr)
		return nil, errors.New(logstr)
	}
	return s, nil
}

func (rs *ResFromFile) GetResFromFile() ([]byte, error) {

	reader := rs.Reader
	s, err := reader.ReadBytes('\n')
	if err != nil {
		rs.Logger.Warn("read data from file fail , " + err.Error())
		return nil, err
	}
	s = s[0 : len(s)-1]
	ss, err := rs.CheckDataValid(s)
	if err!=nil {
		rs.Logger.Warn(" read data is invalid ")
		return nil, err
	}
	return ss, nil

}


func (rs *ResFromFile)UnMarshalToStruct(s []byte) error{
	err := json.Unmarshal(s,rs)
	if err!=nil{
		rs.Logger.Warn("Unmarshal json to struct fail , "+err.Error())
		return err
	}
	return nil
}


type SqlCompareRes struct {
	Sql     string        `json:"sql"`
	Values  []interface{} `json:"values"`
	ErrCode int           `json:"errcode"`
	ErrDesc string        `json:"errdesc"`
}

func (rs *ResFromFile) CompareRes() *SqlCompareRes {

	res := new(SqlCompareRes)
	res.Sql = rs.Query
	res.Values = rs.Params

	m := make(map[string]uint64)

	defer stat.Statis.AddStaticm(m,rs.Logger)
	m["ExecSqlNum"] =1

	// Calculate  SQL execution time in packet
	var prSqlExecTime uint64
	if rs.PrBeginTime < rs.PrEndTime {
		prSqlExecTime = rs.PrEndTime - rs.PrBeginTime
	} else {
		prSqlExecTime = 0
	}
	m["PrExecTimeCount"]=prSqlExecTime

	//Calculate  SQL execution time in replay server
	var rrSqlExecTime uint64
	if rs.RrBeginTime < rs.RrEndTime {
		rrSqlExecTime = rs.RrEndTime - rs.RrBeginTime
	} else {
		rrSqlExecTime = 0
	}
	m["RrExecTimeCount"] = rrSqlExecTime

	stat.Statis.SetBucketNum(rrSqlExecTime,1)
	stat.Statis.SetBucketNum(prSqlExecTime,0)

	var prlen  = 0
	var rrlen  = 0

	prlen = len(rs.PrResult)
	rrlen = len(rs.RrResult)

	m["PrExecRowCount"] = uint64(prlen)
	m["RrExecRowCount"] = uint64(rrlen)

	if rs.PrErrorNo !=0{
		m["PrExecFailCount"] =1
	} else {
		m["PrExecSuccCount"] =1
	}

	if rs.RrErrorNo != 0 {
		m["RrExecFailCount"]=1
	} else {
		m["RrExecSuccCount"]=1
	}
	//compare errcode
	if rs.PrErrorNo != rs.RrErrorNo {
		res.ErrCode = 1
		res.ErrDesc = fmt.Sprintf("%v-%v", rs.PrErrorNo, rs.RrErrorNo)
		m["ExecErrNoNotEqual"]=1
		m["ExecFailNum"]=1
		return res
	}

	//compare exec time

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
		m["ExecFailNum"]=1
		m["ExecTimeNotEqual"]=1

		return res
	}

	//compare  result row num

	if prlen != rrlen {
		res.ErrCode = 3
		res.ErrDesc = fmt.Sprintf("%v-%v", prlen, rrlen)
		m["ExecFailNum"]=1
		m["RowCountNotequal"]=1

		return res
	} else if prlen == 0 {
		res.ErrCode = 0
		m["ExecSuccNum"]=1
		return res
	}

	//compare result row detail
	ok,err := rs.CompareResDetail()
	if !ok{
		rs.Logger.Error("Compare result data fail , "+err.Error())
		m["ExecFailNum"]=1
		m["RowDetailNotEqual"]=1
		res.ErrCode = 4
		return res
	}
	res.ErrCode = 0
	m["ExecSuccNum"]=1

	return res
}


func (rs *ResFromFile) HashPrResDetail() ([][]interface{},error) {
	//var rowStr string
	res := make([][]interface{},0)
	for i := range rs.PrResult {
		rowStr := strings.Join(rs.PrResult[i],"")
		rowv := make([]interface{},0)
		v ,err:= hashString(rowStr)
		if err!=nil{
			return nil ,err
		}
		rowv = append(rowv,v,uint64(i))
		res = append(res,rowv)
	}
	return res,nil
}

func (rs *ResFromFile) HashRrResDetail() ([][]interface{},error) {
	//var rowStr string
	res := make([][]interface{},0)
	for i := range rs.RrResult {
		rowStr := strings.Join(rs.RrResult[i],"")
		rowv := make([]interface{},0)
		v ,err:= hashString(rowStr)
		if err!=nil{
			return nil ,err
		}
		rowv = append(rowv,v,uint64(i))
		res = append(res,rowv)
	}
	return res,nil
}


func (rs *ResFromFile) CompareData(a,b [][]interface{}) error {
	for i := range a {
		ok := CompareInterface(a[i][0],b[i][0])
		if !ok {
			err := errors.New("compare row string hash value not equal ")
			return err
		}
	}
	return nil
}


func (rs *ResFromFile)CompareResDetail() (bool,error){

	a ,err:= rs.HashPrResDetail()
	if err!=nil{
		rs.Logger.Error(" hash Packet result fail , "+err.Error())
		return false,err
	}
	Sort2DSlice(a)
	b ,err:=rs.HashRrResDetail()
	if err!=nil{
		rs.Logger.Error(" hash Replay server  result fail , "+err.Error())
		return false,err
	}
	Sort2DSlice(b)

	err =rs.CompareData(a,b)
	if err!=nil{
		rs.Logger.Warn("compare data fail , " + err.Error())
		return false ,err
	}

	return true ,nil
}


//read result from file ,and compare packet result and replay server result
func DoCompare(fn string ,f *os.File,wg *sync.WaitGroup) {
	defer wg.Done()
	logger := fileName(fn).Logger()
	stat.Statis.AddStatic("ResultFiles",logger)

	rs := NewResForWriteFile(f,logger)
	for {
		s, err := rs.GetResFromFile()
		if err != nil && err!=io.EOF {
			stat.Statis.AddStatic("ReadFailFiles", logger)
			logger.Error("read result file fail , " + err.Error())
			return
		}
		if err == io.EOF{
			break
		}

		err = rs.UnMarshalToStruct(s)
		if err != nil {
			stat.Statis.AddStatic("ReadFailFiles", logger)
			logger.Error("UnMarshal read string to struct , " + err.Error())
			return
		}

		res := rs.CompareRes()
		if res.ErrCode != 0 {
			jsons, err := json.Marshal(res)
			if err != nil {
				logger.Warn("Marshal compare result to json fail , " + err.Error())
			} else {
				log.Warn(string(String(jsons)))
			}
		}
	}
	stat.Statis.AddStatic("ReadSuccFiles",logger)
	return
}

/*
func DoCompare(fn string ,f *os.File , wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Println(fn)
}
*/

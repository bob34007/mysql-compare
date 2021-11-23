package compare

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/bobguo/mysql-compare/parser"
	"github.com/bobguo/mysql-compare/stat"
	"github.com/bobguo/mysql-compare/utils"
	"github.com/pingcap/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)


var BasePercent uint64


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
	defer func() {
		if err := recover(); err != nil {
			//rs.Logger.Warn(err)
			rs.Logger.Warn(f.Name())
		}

	}()
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

func (rs *ResFromFile)DetermineNeedCompareResult() bool{
	if rs.Type == utils.EventQuery || rs.Type == utils.EventStmtExecute {
		isSelect, err := parser.CheckIsSelectStmt(rs.Query)
		if err != nil {
			rs.Logger.Error("determine if SQL is select occurred error" +
				err.Error())
			return false
		}
		if isSelect == false {
			rs.Logger.Debug(rs.Query + " type  is not select , we do not compare result ")
			return false
		}
	} else{
		return false
	}
	return true
}

func (rs *ResFromFile) CompareRes() *SqlCompareRes {

	res := new(SqlCompareRes)
	res.Type = utils.TypeString(rs.Type)
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

	m["PrExecRowCount"] = uint64(len(rs.PrResult))
	m["RrExecRowCount"] = uint64(len(rs.RrResult))

	if rs.PrErrorNo != 0 {
		m["PrExecFailCount"] = 1
		rs.Logger.Warn(fmt.Sprintf("sql exec on pr fail , %s-%v-%v-%s",res.Sql,res.Values,rs.PrErrorNo,rs.PrErrorDesc))
	} else {
		m["PrExecSuccCount"] = 1
	}

	if rs.RrErrorNo != 0 {
		m["RrExecFailCount"] = 1
		rs.Logger.Warn(fmt.Sprintf("sql exec on rr fail , %s-%v-%v-%s",res.Sql,res.Values,rs.RrErrorNo,rs.RrErrorDesc))
	} else {
		m["RrExecSuccCount"] = 1
	}

	//If it is a select statement then continue to the next step
	if !rs.DetermineNeedCompareResult(){
		return nil
	}

	defer res.AddOneSqlResultToSQLStat(rs)

	m["CompareSqlNum"]=1

	//compare errcode
	if rs.PrErrorNo != rs.RrErrorNo {
		res.ErrCode = 1
		res.ErrDesc = fmt.Sprintf("%v-%v", rs.PrErrorNo, rs.RrErrorNo)
		m["ExecErrNoNotEqual"] = 1
		m["CompareFailNum"] = 1
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
	if len(rs.PrResult) != len(rs.RrResult) {
		res.ErrCode = 3
		res.ErrDesc = fmt.Sprintf("%v-%v", len(rs.PrResult) , len(rs.RrResult))
		m["CompareFailNum"] = 1
		m["RowCountNotequal"] = 1
		res.RrValues = rs.RrResult
		res.PrValues = rs.PrResult
		return res
	} else if len(rs.PrResult) == 0 {
		res.ErrCode = 0
		m["CompareSuccNum"] = 1
		return res
	}

	//compare result row detail
	ok, err := rs.CompareResDetail()
	if !ok {
		rs.Logger.Error("Compare result data fail , " + err.Error())
		m["CompareFailNum"] = 1
		m["RowDetailNotEqual"] = 1
		res.ErrCode = 4
		res.RrValues = rs.RrResult
		res.PrValues = rs.PrResult
		return res
	}
	res.ErrCode = 0
	m["CompareSuccNum"] = 1



	return res
}

func (res *SqlCompareRes) AddOneSqlResultToSQLStat(rs *ResFromFile){
	var sqlExecFail ,sqlExecSucc , sqlErrNoNotEqual,sqlRowCountNotEqual,
		sqlRowDetailNotEqual,sqlExecTimePr,sqlExecTimeRr uint64
	if rs.PrBeginTime < rs.PrEndTime {
		sqlExecTimePr = rs.PrEndTime - rs.PrBeginTime
	} else {
		sqlExecTimePr = 0
	}
	if rs.RrBeginTime < rs.RrEndTime {
		sqlExecTimeRr = rs.RrEndTime - rs.RrBeginTime
	} else {
		sqlExecTimeRr = 0
	}
	if rs.PrErrorNo!=rs.RrErrorNo || rs.PrErrorNo!=0 {
		sqlExecFail = 1
	} else {
		sqlExecSucc=1
	}
	switch res.ErrCode {
	case 1 :
		sqlErrNoNotEqual=1
	case 3:
		sqlRowCountNotEqual=1
	case 4:
		sqlRowDetailNotEqual=1
	}

	osr := parser.NewOneSQLResult(res.Sql,rs.Type,sqlExecSucc,sqlExecFail,
		sqlErrNoNotEqual,sqlRowCountNotEqual,sqlRowDetailNotEqual,sqlExecTimePr,
		sqlExecTimeRr,rs.Logger)
	err := osr.OneSQLResultInit()
	if err!=nil{
		rs.Logger.Error(err.Error())
		return
	}
	prAvgTime,rrAvgtime := osr.AddResultToSQLStat()
	rs.PrintExecCostTimeAbnormal(prAvgTime,rrAvgtime)
}


func (rs *ResFromFile)PrintExecCostTimeAbnormal(prAvgTime,rrAvgTime uint64){
	if prAvgTime == 0 && rrAvgTime ==0{
		return
	}
	var logStr string
	if rs.PrEndTime>rs.PrBeginTime && (rs.PrEndTime - rs.PrBeginTime > prAvgTime * (100+BasePercent) / 100) {
		if rs.Type == utils.EventQuery {
			logStr = fmt.Sprintf("sql %s exec cost time abnormal %v-%v,sql begin exec at %v",
				rs.Query,rs.PrEndTime - rs.PrBeginTime , prAvgTime,rs.PrBeginTime)
		}
		if rs.Type == utils.EventStmtExecute {
			logStr = fmt.Sprintf("sql %v with param %v exec cost time abnormal %v-%v,sql begin exec at %v",
				rs.Query,rs.Params,rs.PrEndTime - rs.PrBeginTime , prAvgTime,rs.PrBeginTime)
		}
	}
	if len(logStr) > 0 {
		rs.Logger.Error(logStr)
	}
	if rs.RrEndTime>rs.RrBeginTime && (rs.RrEndTime - rs.RrBeginTime > rrAvgTime * (100+BasePercent)  / 100) {
		if rs.Type == utils.EventQuery {
			logStr = fmt.Sprintf("sql %s exec cost time abnormal %v-%v,sql begin exec at %v",
				rs.Query,rs.RrEndTime - rs.RrBeginTime , rrAvgTime,rs.RrBeginTime)
		}
		if rs.Type == utils.EventStmtExecute {
			logStr = fmt.Sprintf("sql %v with param %v exec cost time abnormal %v-%v,sql begin exec at %v",
				rs.Query,rs.Params,rs.RrEndTime - rs.RrBeginTime , rrAvgTime,rs.RrBeginTime)
		}
	}
	if len(logStr) > 0 {
		rs.Logger.Error(logStr)
	}
	if rs.RrEndTime>rs.RrBeginTime && (rs.RrEndTime - rs.RrBeginTime > prAvgTime * (100+BasePercent) / 100) {
		if rs.Type == utils.EventQuery {
			logStr = fmt.Sprintf("sql %s exec cost time lg than production exec cost avg time  %v-%v,sql begin exec at %v",
				rs.Query,rs.RrEndTime - rs.RrBeginTime , prAvgTime,rs.RrBeginTime)
		}
		if rs.Type == utils.EventStmtExecute {
			logStr = fmt.Sprintf("sql %v with param %v exec cost time lg than production exec cost avg time %v-%v,sql begin exec at %v",
				rs.Query,rs.Params,rs.RrEndTime - rs.RrBeginTime , prAvgTime,rs.RrBeginTime)
		}
	}
	if len(logStr) > 0 {
		rs.Logger.Error(logStr)
	}
}


func (rs *ResFromFile) HashPrResDetail() ([][]interface{}, error) {
	//var rowStr string
	res := make([][]interface{}, 0)
	for i := range rs.PrResult {
		rowStr := strings.Join(rs.PrResult[i], "")
		rowv := make([]interface{}, 0)
		v, err := utils.HashString(rowStr)
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
		v, err := utils.HashString(rowStr)
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
		ok := utils.CompareInterface(a[i][0], b[i][0])
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
	utils.Sort2DSlice(a)
	b, err := rs.HashRrResDetail()
	if err != nil {
		rs.Logger.Error(" hash Replay server  result fail , " + err.Error())
		return false, err
	}
	utils.Sort2DSlice(b)

	err = rs.CompareData(a, b)
	if err != nil {
		rs.Logger.Warn("compare data fail , " + err.Error())
		return false, err
	}

	return true, nil
}

func  DoComparePre(fn string,log *zap.Logger) (*os.File,error){
	f , err:=utils.OpenFile(fn)
	if err !=nil{
		log.Error("open file fail , "+ fn)
		return nil,err
	}
	return f,nil
}

func DoCompareFinish (f *os.File,log *zap.Logger,filePath,backDir,fileName string){
	err :=utils.CloseFile(f)
	if err!=nil {
		log.Error("close file fail , filename "+ fileName +" "+err.Error())
	}
	if len(backDir) ==0{
		return
	}
	err = utils.MoveFileToBackupDir(filePath,fileName,backDir)
	if err!=nil{
		log.Error("back file fail , filename "+ f.Name()+" "+err.Error())
	}
	return
}

//read result from file ,and compare packet result and replay server result
func DoCompare(fileName string,ct *int32,wg *sync.WaitGroup,filePath,backDir string ) {

	defer atomic.AddInt32(ct,-1)
	defer wg.Done()
	fn := filePath+"/"+fileName
	logger := utils.FileName(fn).Logger()
	logger.Info("begin to  process file " + fn)
	defer logger.Info("end to  process file " + fn)

	f ,err := DoComparePre(fn,logger)
	if err !=nil{
		logger.Error("open result data file fail , " + err.Error())
		return
	}
	defer DoCompareFinish(f,logger,filePath,backDir,fn)


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
		if res!=nil&&res.ErrCode != 0 {
			jsons, err := json.Marshal(res)
			if err != nil {
				logger.Warn("Marshal compare result to json fail , " + err.Error())
			} else {
				logger.Warn(string(utils.String(jsons)))
			}
		}

	}
	logger.Info("read data file success" + fn)
	stat.Statis.AddStatic("ReadSuccFiles", logger)
	return
}

/*******************************************************************************
 * Copyright (c)  2021 PingCAP, Inc.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 ******************************************************************************/

/**
 * @Author: guobob
 * @Description:
 * @File:  sql.go
 * @Version: 1.0.0
 * @Date: 2021/11/6 10:31
 */

package parser

import (
	"github.com/bobguo/mysql-compare/stat"
	"github.com/bobguo/mysql-compare/utils"
	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
	"go.uber.org/zap"
	"strings"
)


func (osr *OneSQLResult)GenerateSQLTemplate()  {
	str := strings.Trim(osr.SQL, " ")
	str = strings.ToLower(str)

	osr.Log.Debug("before generate template SQL is : " + str)
	osr.SQLTemplate= parser.Normalize(str)
	osr.Log.Debug("after generate template SQL is : " + osr.SQLTemplate)
}


type CheckIsSelectOrNot struct {
	t bool
}

func (v *CheckIsSelectOrNot) Enter(in ast.Node) (ast.Node, bool) {

	if _, ok := in.(*ast.SelectStmt); ok {
		v.t = true
	}
	return in, false
}

func (v *CheckIsSelectOrNot) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func extract(rootNode *ast.StmtNode) bool {

	v := &CheckIsSelectOrNot{t :false}

	(*rootNode).Accept(v)

	return  v.t
}

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	if len(stmtNodes) ==0{
		return nil,errors.New("parse sql result is nil , " +sql)
	}

	return &stmtNodes[0], nil
}

func CheckIsSelectStmt( sql string) (bool,error){

	astNode, err := parse(sql)
	if err != nil {
		return false,err
	}

	return extract(astNode),nil
}


type   OneSQLResult struct {
	SQLType uint64
	SQL string
	SQLHashKey uint64
	SQLTemplate string
	SQLExecSucc uint64
	SQLExecFail uint64
	SQLErrNoNotEqual    uint64
	SQLRowCountNotEqual uint64
	SQLRowDetailNotEqual uint64
	SQLExecTimePr  uint64
	SQLExecTimeRr  uint64
	Log *zap.Logger
}

func NewOneSQLResult (sql string,sqlType uint64 ,sqlExecSucc,sqlExecFail,sqlErrNoNotEqual,sqlRowCountNotEqual,
	sqlRowDetailNotEqual,sqlExecTimePr,sqlExecTimeRr uint64,log *zap.Logger) *OneSQLResult{

	return  &OneSQLResult{
		SQL:sql,
		SQLType: sqlType,
		SQLExecSucc: sqlExecSucc,
		SQLExecFail: sqlExecFail,
		SQLErrNoNotEqual: sqlErrNoNotEqual,
		SQLRowCountNotEqual: sqlRowCountNotEqual,
		SQLRowDetailNotEqual: sqlRowDetailNotEqual,
		SQLExecTimePr: sqlExecTimePr,
		SQLExecTimeRr: sqlExecTimeRr,
		Log:log,
	}

}


func (osr *OneSQLResult)OneSQLResultInit() error {
	var err error

	if  osr.SQLType ==utils.EventQuery {
		osr.GenerateSQLTemplate()
		osr.SQLHashKey,err  = utils.HashString(osr.SQLTemplate)
		if err != nil{
			return err
		}
	}

	if  osr.SQLType == utils.EventStmtExecute{
		osr.SQLHashKey,err  = utils.HashString(osr.SQL)
		if err != nil{
			return err
		}
	}
	return nil
}

func (osr *OneSQLResult)AddResultToSQLStat()  {
	sqlCompareSucc := 1-osr.SQLErrNoNotEqual - osr.SQLRowCountNotEqual -osr.SQLRowDetailNotEqual
	sqlCompareFail := 1- sqlCompareSucc
	if osr.SQLType ==utils.EventQuery  {
		stat.AddKey(osr.SQLHashKey, osr.SQLTemplate, 1, osr.SQLExecSucc, osr.SQLExecFail, sqlCompareSucc,
			sqlCompareFail, osr.SQLErrNoNotEqual, osr.SQLRowCountNotEqual, osr.SQLRowDetailNotEqual,
			osr.SQLExecTimePr, osr.SQLExecTimeRr, osr.Log)
	}
	if osr.SQLType == utils.EventStmtExecute  {
		stat.AddKey(osr.SQLHashKey, osr.SQL, 1, osr.SQLExecSucc, osr.SQLExecFail, sqlCompareSucc,
			sqlCompareFail, osr.SQLErrNoNotEqual, osr.SQLRowCountNotEqual, osr.SQLRowDetailNotEqual,
			osr.SQLExecTimePr, osr.SQLExecTimeRr, osr.Log)
	}
}
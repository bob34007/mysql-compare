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
 * @File:  sqlstat.go
 * @Version: 1.0.0
 * @Date: 2021/11/8 09:37
 */

package stat

import (
	"encoding/json"
	"github.com/bobguo/mysql-compare/utils"
	"go.uber.org/zap"
	"sync"
)

// Store the results of a sql template
type SQLResult struct {
	//sql template string
	SQLTemplate string `json:"sql"`
	//current sql template execution times
	SQLExecCount uint64 `json:"sql_exec_count"`
	//current sql template execution success times
	SQLExecSuccCount uint64 `json:"sql_exec_succ_count"`
	//current sql template execution fail times
	SQLExecFailCount uint64 `json:"sql_exec_fail_count"`
	//current sql result compare success count
	SQLCompareSuccCount uint64 `json:"sql_compare_succ_count"`
	//current sql result compare fail count
	SQLCompareFailCount uint64 `json:"sql_compare_fail_count"`
	//current sql result compare error no fail count
	SQLCompareErrNoFailCount uint64 `json:"sql_compare_errno_fail_count"`
	//current sql result compare row count fail count
	SQLCompareRowCountFailCount uint64 `json:"sql_compare_rowcount_fail_count"`
	//current sql result compare row detail fail count
	SQLCompareRowDetailFailCount uint64 `json:"sql_compare_rowdetail_fail_count"`
	//sum of current sql execution time on production environment
	SQLExecTimeCountPr uint64 `json:"sql_exec_time_count_pr"`
	//sum of current sql execution time on simulation environment
	SQLExecTimeCountRr uint64 `json:"sql_exec_time_count_rr"`
	//statistics on the number of inconsistent execution times of sql statements
	//in the production and simulation environments
	SQLExecTimeCompareFail uint64 `json:"sql_exec_time_compare_fail"`
	//current sql execution time deterioration statistics min
	SQLExecTimeStandard uint64 `json:"sql_exec_time_standard"`
}


// key : SQLTemplates hash value
//value : slice SQLResult, preventing hash collisions
//used to cache all SQL template execution results
var SQLResults map[uint64][]*SQLResult
//for protecting SQL template storage structures
var Mu *sync.RWMutex
//var Log *zap.Logger


func init(){
	SQLResults = make(map[uint64][]*SQLResult)
	Mu  = new(sync.RWMutex)
}



//add key to map
func AddKey( key uint64,SQL string, ExecSQL, ExecSQLSucc, ExecSQLFail, SQLCompareSucc,
	SQLCompareFail, SQLCompareErrNoFail, SQLCompareRowCountFail,
	SQLCompareRowDetailFail, SQLExecTimePr, SQLExecTimeRr uint64,log *zap.Logger)  {
	Mu.Lock()
	defer Mu.Unlock()
	if v, ok := SQLResults[key]; ok {
		var found = false
		for i, _ := range v {
			if SQL == (*SQLResults[key][i]).SQLTemplate {
				(*SQLResults[key][i]).SQLExecCount += ExecSQL
				(*SQLResults[key][i]).SQLExecSuccCount += ExecSQLSucc
				(*SQLResults[key][i]).SQLExecFailCount += ExecSQLFail
				(*SQLResults[key][i]).SQLCompareSuccCount += SQLCompareSucc
				(*SQLResults[key][i]).SQLCompareFailCount += SQLCompareFail
				(*SQLResults[key][i]).SQLCompareErrNoFailCount += SQLCompareErrNoFail
				(*SQLResults[key][i]).SQLCompareRowCountFailCount += SQLCompareRowCountFail
				(*SQLResults[key][i]).SQLCompareRowDetailFailCount += SQLCompareRowDetailFail
				(*SQLResults[key][i]).SQLExecTimeCountPr += SQLExecTimePr
				(*SQLResults[key][i]).SQLExecTimeCountRr += SQLExecTimeRr
				found = true
				break
			}
		}
		if !found {
			SQLResults[key] = append(SQLResults[key], &SQLResult{
				SQLTemplate:                 SQL,
				SQLExecCount:                 ExecSQL,
				SQLExecSuccCount:             ExecSQLSucc,
				SQLExecFailCount:             ExecSQLFail,
				SQLCompareSuccCount:          SQLCompareSucc,
				SQLCompareFailCount:          SQLCompareFail,
				SQLCompareErrNoFailCount:     SQLCompareErrNoFail,
				SQLCompareRowCountFailCount:  SQLCompareRowCountFail,
				SQLCompareRowDetailFailCount: SQLCompareRowDetailFail,
				SQLExecTimeCountPr:           SQLExecTimePr,
				SQLExecTimeCountRr:           SQLExecTimeRr,
			})
		}
	} else {
		sliceSQLResult := make([]*SQLResult, 0)
		sliceSQLResult = append(sliceSQLResult, &SQLResult{
			SQLTemplate:                 SQL,
			SQLExecCount:                 ExecSQL,
			SQLExecSuccCount:             ExecSQLSucc,
			SQLExecFailCount:             ExecSQLFail,
			SQLCompareSuccCount:          SQLCompareSucc,
			SQLCompareFailCount:          SQLCompareFail,
			SQLCompareErrNoFailCount:     SQLCompareErrNoFail,
			SQLCompareRowCountFailCount:  SQLCompareRowCountFail,
			SQLCompareRowDetailFailCount: SQLCompareRowDetailFail,
			SQLExecTimeCountPr:           SQLExecTimePr,
			SQLExecTimeCountRr:           SQLExecTimeRr,
		})
		SQLResults[key] = sliceSQLResult
	}
	log.Debug(" add key " + SQL + "to map success ")
	return
}

//generate SQLResult struct to string
func (rs *SQLResult) String() (string, error) {

	s, err := json.Marshal(rs)
	if err != nil {
		return "", err
	} else {
		return string(utils.String(s)), nil
	}
}

//write sql template result to file
func PrintMap(log *zap.Logger) error {
	Mu.RLock()
	defer Mu.RUnlock()
	//fmt.Println(SQLResults)
	for _, v := range SQLResults {
		for _, rs := range v {
			str, err := rs.String()
			if err != nil {
				return err
			}
			log.Info(str)
		}
	}
	return nil
}


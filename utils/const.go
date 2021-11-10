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
 * @File:  const.go
 * @Version: 1.0.0
 * @Date: 2021/11/9 11:26
 */

package utils


const (
	EventHandshake uint64 = iota
	EventQuit
	EventQuery
	EventStmtPrepare
	EventStmtExecute
	EventStmtClose
)

func TypeString(t uint64) string {
	switch t {
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
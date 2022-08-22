/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pagehelper

import (
	"fmt"
	"strings"
)

var MysqlModifier = Modifier{
	OrderBy: MysqlModifyOrderBy,
	Page:    MysqlModifyPage,
	Count:   MysqlModifyCount,
}

func MysqlModifyOrderBy(sql string, p *OrderByInfo) string {
	if p.Field == "" {
		return sql
	}
	b := strings.Builder{}
	b.WriteString(strings.TrimSpace(sql))
	b.WriteString(fmt.Sprintf(" ORDER BY `%s` %s ", p.Field, p.Order))
	return b.String()
}

func MysqlModifyPage(sql string, p *PageInfo) string {
	b := strings.Builder{}
	b.WriteString(strings.TrimSpace(sql))
	b.WriteString(fmt.Sprintf(" LIMIT %d, %d ", p.Page*p.PageSize, p.PageSize))
	return b.String()
}

func MysqlModifyCount(sql, countColumn string) string {
	if countColumn == "" {
		countColumn = "0"
	} else {
		countColumn = "`" + countColumn + "`"
	}
	b := strings.Builder{}
	b.WriteString("SELECT COUNT(")
	b.WriteString(countColumn)
	b.WriteString(") FROM (")
	b.WriteString(strings.TrimSpace(sql))
	b.WriteString(") AS __hp_tempCountTl")
	return b.String()
}

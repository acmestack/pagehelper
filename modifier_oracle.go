// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

import (
	"fmt"
	"strings"
)

var OracleModifier = Modifier{
	OrderBy: OracleModifyOrderBy,
	Page:    OracleModifyPage,
	Count:   OracleModifyCount,
}

func OracleModifyOrderBy(sql string, p *OrderByInfo) string {
	if p.Field == "" {
		return sql
	}
	b := strings.Builder{}
	b.WriteString(strings.TrimSpace(sql))
	b.WriteString(fmt.Sprintf(" ORDER BY `%s` %s ", p.Field, p.Order))
	return b.String()
}

func OracleModifyPage(sql string, p *PageInfo) string {
	b := strings.Builder{}
	b.WriteString("SELECT * FROM ( ")
	b.WriteString(" SELECT __hp_tempPageTl.*, ROWNUM ROW_ID FROM ( ")
	b.WriteString(strings.TrimSpace(sql))
	b.WriteString(" ) __hp_tempPageTl)")
	b.WriteString(fmt.Sprintf(" WHERE ROW_ID <= %d AND ROW_ID > %d ", p.Page*p.PageSize, p.PageSize))
	return b.String()
}

func OracleModifyCount(sql, countColumn string) string {
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

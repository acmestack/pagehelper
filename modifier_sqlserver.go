// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

import (
    "fmt"
    "strings"
)

var SqlServerModifier = Modifier{
    OrderBy: SqlServerModifyOrderBy,
    Page:    SqlServerModifyPage,
    Count:   SqlServerModifyCount,
}

func SqlServerModifyOrderBy(sql string, p *OrderByInfo) string {
    if p.Field == "" {
        return sql
    }
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" ORDER BY `%s` %s ", p.Field, p.Order))
    return b.String()
}

func SqlServerModifyPage(sql string, p *PageInfo) string {
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY ", p.Page*p.PageSize, p.PageSize))
    return b.String()
}

func SqlServerModifyCount(sql, countColumn string) string {
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

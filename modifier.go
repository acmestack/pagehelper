// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

const (
    DriverDummy     = "default"
    DriverMysql     = "mysql"
    DriverPostgre   = "postgre"
    DriverOracle    = "oracle"
    DriverSqlServer = "sqlserver"
)

type Modifier struct {
    OrderBy func(sql string, p *OrderByInfo) string
    Page    func(sql string, p *PageInfo) string
    Count   func(sql, countColumn string) string
}

var DummyModifier = Modifier{
    OrderBy: DummyModifyOrderBy,
    Page:    DummyModifyPage,
    Count:   DummyModifyCount,
}

func DummyModifyOrderBy(sql string, p *OrderByInfo) string {
    return sql
}

func DummyModifyPage(sql string, p *PageInfo) string {
    return sql
}

func DummyModifyCount(sql, countColumn string) string {
    return sql
}

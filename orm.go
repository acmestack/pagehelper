// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

import (
    "context"
    "fmt"
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/common"
    "github.com/xfali/gobatis/executor"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
    "github.com/xfali/gobatis/reflection"
    "github.com/xfali/gobatis/session"
    "github.com/xfali/gobatis/transaction"
    "strings"
)

const (
    pageHelperValue  = "_page_helper_value"
    orderHelperValue = "_order_helper_value"

    ASC  = "ASC"
    DESC = "DESC"
)

var (
    OrderByModifier = modifyOrderSql
    PageModifier    = modifyPageSql
    CountModifier   = modifyCountSql
)

type Factory struct {
    fac factory.Factory
}

type Executor struct {
    exec executor.Executor
    log  logging.LogFunc
}

func New(f factory.Factory) *Factory {
    return &Factory{f}
}

func (exec *Executor) Close(rollback bool) {
    exec.exec.Close(rollback)
}

func (exec *Executor) Exec(ctx context.Context, sql string, params ...interface{}) (common.Result, error) {
    return exec.exec.Exec(ctx, sql, params...)
}

func (exec *Executor) Begin() error {
    return exec.exec.Begin()
}

func (exec *Executor) Commit(require bool) error {
    return exec.exec.Commit(require)
}

func (exec *Executor) Rollback(require bool) error {
    return exec.exec.Rollback(require)
}

func (exec *Executor) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
    originSql := sql
    o := ctx.Value(orderHelperValue)
    if o != nil {
        if param, ok := o.(*OrderByInfo); ok {
            sql = OrderByModifier(sql, param)
        }
    }

    p := ctx.Value(pageHelperValue)
    if p != nil {
        if param, ok := p.(*PageInfo); ok {
            sql = PageModifier(sql, param)
            if param.total == -1 {
                param.total = exec.getTotal(ctx, originSql, param.countColumn, params...)
            }
        }
    }
    exec.log(logging.DEBUG, "PageHelper Query: [%s], params: %s\n", sql, fmt.Sprint(params))
    return exec.exec.Query(ctx, result, sql, params...)
}

func (f *Factory) InitDB() error {
    return f.fac.InitDB()
}

func (f *Factory) CreateTransaction() transaction.Transaction {
    return f.fac.CreateTransaction()
}

func (f *Factory) CreateSession() session.SqlSession {
    tx := f.CreateTransaction()
    return session.NewDefaultSqlSession(f.LogFunc(), tx, f.CreateExecutor(tx), false)
}

func (f *Factory) LogFunc() logging.LogFunc {
    return f.fac.LogFunc()
}

func (f *Factory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
    return &Executor{
        exec: executor.NewSimpleExecutor(transaction),
        log:  f.LogFunc(),
    }
}

func modifyOrderSql(sql string, p *OrderByInfo) string {
    if p.Field == "" {
        return sql
    }
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" ORDER BY `%s` %s ", p.Field, p.Order))
    return b.String()
}

func modifyPageSql(sql string, p *PageInfo) string {
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" LIMIT %d, %d ", p.Page*p.PageSize, p.PageSize))
    return b.String()
}

func (exec *Executor) getTotal(ctx context.Context, sql, countColumn string, params ...interface{}) int64 {
    totalSql := CountModifier(sql, countColumn)
    var total int64
    obj, err := gobatis.ParseObject(&total)
    if err == nil {
        exec.exec.Query(ctx, obj, totalSql, params...)
        return total
    }
    return 0
}

func modifyCountSql(sql, countColumn string) string {
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

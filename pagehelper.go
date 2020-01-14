/**
 * Copyright (C) 2019, Xiongfa Li.
 * All right reserved.
 * @author xiongfa.li
 * @version V1.0
 * Description:
 */

package pagehelper

import (
    "context"
    "fmt"
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

type OrderParam struct {
    Field string
    Order string
}

type PageParam struct {
    Page     int
    PageSize int
}

func New(f factory.Factory) *Factory {
    return &Factory{f}
}

type Factory struct {
    fac factory.Factory
}

type Executor struct {
    exec executor.Executor
    log  logging.LogFunc
}

func StartPage(page, pageSize int, ctx context.Context) context.Context {
    return context.WithValue(ctx, pageHelperValue, &PageParam{Page: page, PageSize: pageSize})
}

func OrderBy(field, order string, ctx context.Context) context.Context {
    return context.WithValue(ctx, orderHelperValue, &OrderParam{Field: field, Order: order})
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
    o := ctx.Value(orderHelperValue)
    if o != nil {
        if param, ok := o.(*OrderParam); ok {
            sql, params = modifySqlOrder(sql, param, params)
        }
    }

    p := ctx.Value(pageHelperValue)
    if p != nil {
        if param, ok := p.(*PageParam); ok {
            sql = modifySql(sql, param)
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

func modifySqlOrder(sql string, p *OrderParam, params []interface{}) (string, []interface{}) {
    if p.Field == "" {
        return sql, params
    }
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" ORDER BY ? %s ", p.Order))
    params = append(params, p.Field)
    return b.String(), params
}

func modifySql(sql string, p *PageParam) string {
    b := strings.Builder{}
    b.WriteString(strings.TrimSpace(sql))
    b.WriteString(fmt.Sprintf(" LIMIT %d, %d ", p.Page*p.PageSize, p.PageSize))
    return b.String()
}

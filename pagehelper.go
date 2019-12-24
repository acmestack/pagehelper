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
    "github.com/xfali/gobatis/executor"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/reflection"
    "github.com/xfali/gobatis/transaction"
    "strings"
)

const (
    pageHelperValue = "_page_helper_value"
)

type PageParam struct {
    Page     int
    PageSize int
}

func New(f *factory.DefaultFactory) *PageHelperFactory {
    return &PageHelperFactory{*f}
}

type PageHelperFactory struct {
    factory.DefaultFactory
}

type PageHelperExecutor struct {
    executor.SimpleExecutor
}

func StartPage(page, pageSize int, ctx context.Context) context.Context {
    return context.WithValue(ctx, pageHelperValue, &PageParam{Page: page, PageSize: pageSize})
}

func (exec *PageHelperExecutor) Query(ctx context.Context, result reflection.Object, sql string, params ...interface{}) error {
    p := ctx.Value(pageHelperValue)
    if p != nil {
        if param, ok := p.(*PageParam); ok {
            sql = modifySql(sql, param)
        }
    }

    return exec.SimpleExecutor.Query(ctx, result, sql, params...)
}

func (f *PageHelperFactory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
    return &PageHelperExecutor{*executor.NewSimpleExecutor(transaction)}
}

func modifySql(sql string, p *PageParam) string {
    b := strings.Builder{}
    b.WriteString(sql)
    b.WriteString(fmt.Sprintf(" OFFSET %d LIMIT %d ", p.Page*p.PageSize, p.PageSize))
    return b.String()
}

/*
 * Copyright (c) 2022, AcmeStack
 * All rights reserved.
 *
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
	"context"
	"fmt"

	"github.com/acmestack/gobatis"
	"github.com/acmestack/gobatis/common"
	"github.com/acmestack/gobatis/datasource"
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/factory"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/reflection"
	"github.com/acmestack/gobatis/session"
	"github.com/acmestack/gobatis/transaction"
)

const (
	pageHelperValue  = "_page_helper_value"
	orderHelperValue = "_order_helper_value"

	ASC  = "ASC"
	DESC = "DESC"
)

type Factory struct {
	fac factory.Factory
}

type Executor struct {
	exec executor.Executor
	log  logging.LogFunc

	modifier Modifier
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
			sql = exec.modifier.OrderBy(sql, param)
		}
	}

	p := ctx.Value(pageHelperValue)
	if p != nil {
		if param, ok := p.(*PageInfo); ok {
			sql = exec.modifier.Page(sql, param)
			if param.total == -1 {
				param.total = exec.getTotal(ctx, originSql, param.countColumn, params...)
			}
		}
	}
	exec.log(logging.DEBUG, "PageHelper Query: [%s], params: %s\n", sql, fmt.Sprint(params))
	return exec.exec.Query(ctx, result, sql, params...)
}

func (f *Factory) Open(source datasource.DataSource) error {
	return f.fac.Open(source)
}

func (f *Factory) Close() error {
	return f.fac.Close()
}

func (f *Factory) GetDataSource() datasource.DataSource {
	return f.fac.GetDataSource()
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
	driver := f.fac.GetDataSource().DriverName()

	return &Executor{
		exec:     f.fac.CreateExecutor(transaction),
		log:      f.LogFunc(),
		modifier: SelectModifier(driver),
	}
}

func (exec *Executor) getTotal(ctx context.Context, sql, countColumn string, params ...interface{}) int64 {
	totalSql := exec.modifier.Count(sql, countColumn)
	var total int64
	obj, err := gobatis.ParseObject(&total)
	if err == nil {
		exec.exec.Query(ctx, obj, totalSql, params...)
		return total
	}
	return 0
}

var modifierMap = map[string]Modifier{
	DriverDummy:     DummyModifier,
	DriverMysql:     MysqlModifier,
	DriverOracle:    OracleModifier,
	DriverPostgre:   PostgreModifier,
	DriverSqlServer: SqlServerModifier,
}

func RegisterModifier(driver string, m Modifier) {
	if driver != "" {
		modifierMap[driver] = m
	}
}

func SelectModifier(driver string) Modifier {
	if m, ok := modifierMap[driver]; ok {
		return m
	}
	return modifierMap[DriverDummy]
}

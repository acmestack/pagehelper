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
	"github.com/acmestack/gobatis/executor"
	"github.com/acmestack/gobatis/logging"
	"github.com/acmestack/gobatis/session"
	"github.com/acmestack/gobatis/transaction"
)

type IFactory struct {
	InitDBFunc            func() error
	CreateTransactionFunc func() transaction.Transaction
	CreateExecutorFunc    func(transaction.Transaction) executor.Executor
	CreateSessionFunc     func() session.SqlSession
	LogFuncFunc           func() logging.LogFunc
}

func (f *IFactory) InitDB() error {
	return f.InitDBFunc()
}

func (f *IFactory) CreateTransaction() transaction.Transaction {
	return f.CreateTransactionFunc()
}

func (f *IFactory) CreateSession() session.SqlSession {
	tx := f.CreateTransactionFunc()
	return session.NewDefaultSqlSession(f.LogFuncFunc(), tx, f.CreateExecutorFunc(tx), false)
}

func (f *IFactory) LogFunc() logging.LogFunc {
	return f.LogFuncFunc()
}

func (f *IFactory) CreateExecutor(transaction transaction.Transaction) executor.Executor {
	return f.CreateExecutorFunc(transaction)
}

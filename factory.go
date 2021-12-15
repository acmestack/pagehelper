// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

import (
	"github.com/xfali/gobatis/executor"
	"github.com/xfali/gobatis/logging"
	"github.com/xfali/gobatis/session"
	"github.com/xfali/gobatis/transaction"
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

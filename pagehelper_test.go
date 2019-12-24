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
    "github.com/xfali/gobatis"
    "github.com/xfali/gobatis/factory"
    "github.com/xfali/gobatis/logging"
    "testing"
    "time"
)

type TestTable struct {
    TestTable gobatis.ModelName "test_table"
    Id        int64             `xfield:"id"`
    Username  string            `xfield:"username"`
    Password  string            `xfield:"password"`
}

func TestPageHelper(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = StartPage(1, 2, ctx)

    p := ctx.Value(pageHelperValue)
    printV(t, p)

    select {
    case <-ctx.Done():
        break
    }
    printV(t, p)
}

func TestPageHelper2(t *testing.T) {
    fac := factory.DefaultFactory{
        Host:     "localhost",
        Port:     3306,
        DBName:   "test",
        Username: "root",
        Password: "123",
        Charset:  "utf8",

        MaxConn:     1000,
        MaxIdleConn: 500,

        Log: logging.DefaultLogf,
    }
    sessMgr := gobatis.NewSessionManager(New(&fac))
    session := sessMgr.NewSession()
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = StartPage(1, 2, ctx)

    session.SetContext(ctx)

    var ret []TestTable
    session.Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
}

func TestModify(t *testing.T) {
    sql := modifySql("select * from x", &PageParam{1, 2})
    t.Log(sql)
}

func printV(t *testing.T, p interface{}) {
    if p, ok := p.(*PageParam); ok {
        t.Logf("param: %d %d", p.Page, p.PageSize)
    } else {
        t.Fail()
    }
}

func TestContext(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = context.WithValue(ctx, "1", "a")

    t.Log(ctx.Value("1"))
    select {
    case <-ctx.Done():
        break
    }
    t.Log(ctx.Value("1"))
}

type A struct{ I int }
type B struct{ A }

func TestStruct(t *testing.T) {
    a := &A{10}
    b := B{*a}
    t.Logf("b:%d\n", b.I)
    if b.I != 10 {
        t.Fail()
    }
}

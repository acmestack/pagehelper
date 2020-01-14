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
    "strings"
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
    t.Run("StartPage", func(t *testing.T) {
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        ctx = StartPage(1, 2, ctx)

        p := ctx.Value(pageHelperValue)
        printPage(t, p)

        select {
        case <-ctx.Done():
            break
        }
        printPage(t, p)
    })

    t.Run("OrderBy", func(t *testing.T) {
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        ctx = OrderBy("test", ASC, ctx)

        p := ctx.Value(orderHelperValue)
        printOrder(t, p)

        select {
        case <-ctx.Done():
            break
        }
        printOrder(t, p)
    })

    t.Run("PageHelper and OrderBy", func(t *testing.T) {
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        ctx = OrderBy("test", ASC, ctx)
        ctx = StartPage(1, 2, ctx)

        o := ctx.Value(orderHelperValue)
        printOrder(t, o)

        p := ctx.Value(pageHelperValue)
        printPage(t, p)

        select {
        case <-ctx.Done():
            break
        }
        printPage(t, p)
        printOrder(t, o)
    })

    t.Run("complex", func(t *testing.T) {
        ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
        ctx = OrderBy("test", ASC, ctx)
        ctx = StartPage(1, 2, ctx)
        ctx = StartPage(3, 10, ctx)
        ctx = OrderBy("tat", DESC, ctx)
        ctx, _ = context.WithTimeout(ctx, time.Second)

        now := time.Now()
        o := ctx.Value(orderHelperValue)
        printOrder(t, o)
        t.Logf("time :%d ms \n", time.Since(now)/time.Millisecond)

        p := ctx.Value(pageHelperValue)
        printPage(t, p)

        select {
        case <-ctx.Done():
            break
        }
        printPage(t, p)
        printOrder(t, o)
    })
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

func TestModifyPage(t *testing.T) {
    sql := modifySql("select * from x", &PageParam{1, 2})
    t.Log(sql)
}

func order(sql string, params ...interface{}) (string, []interface{}) {
    return modifySqlOrder(sql, &OrderParam{"test", ASC}, params)
}

func TestModifyOrder(t *testing.T) {
    sql, p := order("select ? from x", "field1")
    t.Log(sql)
    if len(p) != 2 {
        t.Fatal()
    }
    for _, v := range p {
        t.Log(v)
    }
}

func TestModifyOrderAndPage(t *testing.T) {
    sql, p := order("select ? from x", "field1")
    t.Log(sql)
    if len(p) != 2 {
        t.Fatal()
    }

    sql = modifySql(sql, &PageParam{1, 2})

    t.Log(sql)
    for _, v := range p {
        t.Log(v)
    }

    if strings.TrimSpace(sql) != "select ? from x ORDER BY ? ASC LIMIT 2, 2" {
        t.Fail()
    }
}

func printPage(t *testing.T, p interface{}) {
    if p, ok := p.(*PageParam); ok {
        t.Logf("page param: %d %d", p.Page, p.PageSize)
    } else {
        t.Fail()
    }
}

func printOrder(t *testing.T, p interface{}) {
    if p, ok := p.(*OrderParam); ok {
        t.Logf("order param: %s %s", p.Field, p.Order)
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

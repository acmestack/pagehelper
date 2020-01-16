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
    "github.com/xfali/gobatis/datasource"
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
        ctx = StartPage(ctx, 1, 2)

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
        ctx = OrderBy(ctx, "test", ASC)

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
        ctx = OrderBy(ctx, "test", ASC)
        ctx = StartPage(ctx, 1, 2)

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
        ctx = OrderBy(ctx, "test", ASC)
        ctx = StartPage(ctx, 1, 2)
        ctx = StartPage(ctx, 3, 10)
        ctx = OrderBy(ctx, "tat", DESC)
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
    pFac := New(gobatis.NewFactory(
        gobatis.SetMaxConn(100),
        gobatis.SetMaxIdleConn(50),
        gobatis.SetDataSource(&datasource.MysqlDataSource{
            Host:     "localhost",
            Port:     3306,
            DBName:   "test",
            Username: "root",
            Password: "123",
            Charset:  "utf8",
        })))
    sessMgr := gobatis.NewSessionManager(pFac)
    session := sessMgr.NewSession()
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = StartPage(ctx, 1, 2)

    session.SetContext(ctx)

    var ret []TestTable
    session.Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
}

func TestModifyPage(t *testing.T) {
    sql := PageModifier("select * from x", &PageInfo{Page: 1, PageSize: 2,})
    t.Log(sql)
    if strings.TrimSpace(sql) != `select * from x LIMIT 2, 2` {
        t.Fail()
    }
}

func order(sql string, params ...interface{}) (string, []interface{}) {
    return OrderByModifier(sql, &OrderByInfo{"test", ASC}), params
}

func TestModifyOrder(t *testing.T) {
    sql, p := order("select ? from x", "field1")
    t.Log(sql)
    for _, v := range p {
        t.Log(v)
    }

    if strings.TrimSpace(sql) != "select ? from x ORDER BY `test` ASC" {
        t.Fail()
    }
}

func TestModifyCount(t *testing.T) {
    t.Run("empty", func(t *testing.T) {
        sql := CountModifier("select ? from x", "")
        t.Log(sql)

        if strings.TrimSpace(sql) != "SELECT COUNT(0) FROM (select ? from x) AS __hp_tempCountTl" {
            t.Fail()
        }
    })

    t.Run("test", func(t *testing.T) {
        sql := CountModifier("select ? from x", "test")
        t.Log(sql)

        if strings.TrimSpace(sql) != "SELECT COUNT(`test`) FROM (select ? from x) AS __hp_tempCountTl" {
            t.Fail()
        }
    })
}

func TestChangeModifyCount(t *testing.T) {
    CountModifier = func(sql, c string) string {
        return "test " + sql
    }
    sql := CountModifier("select ? from x", "")
    t.Log(sql)

    if strings.TrimSpace(sql) != "test select ? from x" {
        t.Fail()
    }
}

func TestGetTotal(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = OrderBy(ctx, "test", ASC)
    ctx = StartPage(ctx, 1, 2)

    total := GetTotal(ctx)
    t.Log(total)
    if total != 0 {
        t.Fail()
    }
}

func TestGetPageInfo(t *testing.T) {
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = OrderBy(ctx, "test", ASC)
    ctx = StartPage(ctx, 1, 10)

    pageInfo := GetPageInfo(ctx)
    pageInfo.total = 1001
    t.Log(
        "pageNum: ", pageInfo.GetPageNum(),
        "totalPage: ", pageInfo.GetTotalPage(),
        "pageSize: ", pageInfo.GetPageSize(),
        "total: ", pageInfo.GetTotal())
    if pageInfo.GetTotalPage() != 101 {
        t.Fail()
    }
}

func TestModifyOrderAndPage(t *testing.T) {
    sql, p := order("select ? from x", "field1")
    t.Log(sql)

    sql = PageModifier(sql, &PageInfo{Page: 1, PageSize: 2,})

    t.Log(sql)
    for _, v := range p {
        t.Log(v)
    }

    if strings.TrimSpace(sql) != "select ? from x ORDER BY `test` ASC LIMIT 2, 2" {
        t.Fail()
    }
}

func printPage(t *testing.T, p interface{}) {
    if p, ok := p.(*PageInfo); ok {
        t.Logf("page param: %d %d", p.Page, p.PageSize)
    } else {
        t.Fail()
    }
}

func printOrder(t *testing.T, p interface{}) {
    if p, ok := p.(*OrderByInfo); ok {
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

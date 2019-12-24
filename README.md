# pagehelper

## 介绍

pagehelper是与[gobatis](https://github.com/xfali/gobatis)配套的分页工具

## 待完成项

select查询返回的result还未能携带page、pageSize、total信息，需要自己组装。

（此功能需gobatis支持，暂未实现）

## 使用

### 1、使用pagehelper的Factory生成SessionManager

```cassandraql
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
    sessMgr := gobatis.NewSessionManager(pagehelper.New(&fac))
```

### 2、配置分页参数
```cassandraql
    session := sessMgr.NewSession()
    ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
    ctx = pagehelper.StartPage(1, 2, ctx)

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```
# pagehelper

## 介绍

pagehelper是与[gobatis](https://github.com/xfali/gobatis)配套的分页工具

## 待完成项

select查询返回的result还未能携带page、pageSize、total信息，需要自己组装。

（此功能需gobatis支持，暂未实现）

## 使用

### 1、使用pagehelper的Factory生成SessionManager

```cassandraql
    pFac := pagehelper.New(&factory.DefaultFactory{
                                    Host:     "localhost",
                                    Port:     3306,
                                    DBName:   "test",
                                    Username: "root",
                                    Password: "123",
                                    Charset:  "utf8",
                            
                                    MaxConn:     1000,
                                    MaxIdleConn: 500,
                            
                                    Log: logging.DefaultLogf,
                                })
    pFac.InitDB()
    sessMgr := gobatis.NewSessionManager(pFac)
```

### 2、配置分页参数
```cassandraql
    session := sessMgr.NewSession()
    ctx := pagehelper.StartPage(1, 10, context.Background())

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```

### 3、配置排序参数
```cassandraql
    session := sessMgr.NewSession()
    ctx := pagehelper.OrderBy("myfield", pagehelper.DESC, context.Background())

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```
*注意:*
  
由于golang对order by不能使用placeholder的方式，所以存在注入风险，请谨慎使用排序功能，如果使用，则需要自己做防注入的工作。

举例：

在获得order by参数时做参数校验
```$xslt
valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
if !valid.MatchString(ordCol) {
    // invalid column name, do not proceed in order to prevent SQL injection
}
```

### 4、使用builder
```$xslt
session := sessMgr.NewSession()
ctx := pagehelper.C(context.Background()).Page(1, 3).Order("test", pagehelper.ASC).Build()
var ret []TestTable
session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```
# pagehelper

## 介绍

pagehelper是与[gobatis](https://github.com/xfali/gobatis)配套的分页工具

## 待完成项

v0.1.0已添加此特性 ：返回总记录数

~~select查询返回的result还未能携带page、pageSize、total信息，需要自己组装。~~

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
    ctx := pagehelper.StartPage(context.Background(), 1, 10)

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```

### 3、配置分页参数（带自动统计总记录数功能）
```$xslt
    session := sessMgr.NewSession()
    ctx := pagehelper.StartPageWithCount(context.Background(), 1, 10, "")

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```
获得分页信息（以1001条记录为例）
```$xslt
    pageInfo := pagehelper.GetPageInfo(ctx)
    t.Log(
        "pageNum: ", pageInfo.GetPageNum(),
        "totalPage: ", pageInfo.GetTotalPage(),
        "pageSize: ", pageInfo.GetPageSize(),
        "total: ", pageInfo.GetTotal())
```
输出：
```$xslt
pageNum:        1 
totalPage:      101 
pageSize:       10 
total:          1001
```
*注意：*

会自动生成和执行带子查询的countSQL，请自行评估是否使用此功能，转而使用自定义SQL获取总记录数。

### 4、配置排序参数
```cassandraql
    session := sessMgr.NewSession()
    ctx := pagehelper.OrderBy(context.Background(), "myfield", pagehelper.DESC)

    var ret []TestTable
    session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```
*注意：*
  
由于golang对order by不能使用placeholder的方式，所以存在注入风险，请谨慎使用排序功能，如果使用，则需要自己做防注入的工作。

举例：

在获得order by参数时做参数校验
```$xslt
valid := regexp.MustCompile("^[A-Za-z0-9_]+$")
if !valid.MatchString(ordCol) {
    // invalid column name, do not proceed in order to prevent SQL injection
}
```

### 5、使用builder
```$xslt
session := sessMgr.NewSession()
ctx := pagehelper.C(context.Background()).Page(1, 3).OrderBy("test", pagehelper.ASC).Build()
var ret []TestTable
session.SetContext(ctx).Select("SELECT * FROM TBL_TEST").Param().Result(&ret)
```

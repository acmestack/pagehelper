// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package pagehelper

import "context"

type builder struct {
    ctx context.Context
}

//创建builder
//ctx 初始context
func C(ctx context.Context) *builder {
    return &builder{ctx: ctx}
}

//分页
//page 页码
//pageSize 分页大小
func (b *builder) Page(page, pageSize int) *builder {
    b.ctx = StartPage(page, pageSize, b.ctx)
    return b
}

//分页
//page 页码
//pageSize 分页大小
func (b *builder) PageWithTotal(page, pageSize int) *builder {
    b.ctx = StartPageWithTotal(page, pageSize, b.ctx)
    return b
}

//手动指定字段和排序
//field 字段
//order 排序 [ASC | DESC]
func (b *builder) Order(field, order string) *builder {
    b.ctx = OrderBy(field, order, b.ctx)
    return b
}

//升序（默认）
//field 字段
func (b *builder) ASC(field string) *builder {
    b.ctx = OrderBy(field, ASC, b.ctx)
    return b
}

//降序
//field 字段
func (b *builder) DESC(field string) *builder {
    b.ctx = OrderBy(field, DESC, b.ctx)
    return b
}

//获得含分页/排序信息的context
func (b *builder) Build() context.Context {
    return b.ctx
}

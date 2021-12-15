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
	"fmt"
)

type OrderByInfo struct {
	Field string
	Order string
}

type PageInfo struct {
	Page     int64
	PageSize int64

	countColumn string
	total       int64
}

func GetPageInfo(ctx context.Context) *PageInfo {
	if ctx == nil {
		return nil
	}
	p := ctx.Value(pageHelperValue)
	if p != nil {
		if param, ok := p.(*PageInfo); ok {
			return param
		}
	}
	return nil
}

func GetTotal(ctx context.Context) int64 {
	p := GetPageInfo(ctx)
	if p != nil {
		return p.total
	}
	return 0
}

// 分页
// page 页码
// pageSize 分页大小
// ctx 初始context
func StartPage(ctx context.Context, page, pageSize int64) context.Context {
	return context.WithValue(ctx, pageHelperValue, &PageInfo{Page: page, PageSize: pageSize, total: 0})
}

// 分页(包含total信息)
// page 页码
// pageSize 分页大小
// 用于统计的列名称，如果为空，则默认使用COUNT(0)统计
// ctx 初始context
func StartPageWithCount(ctx context.Context, page, pageSize int64, countColumn string) context.Context {
	return context.WithValue(ctx, pageHelperValue, &PageInfo{Page: page, PageSize: pageSize, total: -1, countColumn: countColumn})
}

//排序
//field 字段
//order 排序 [ASC | DESC]
//ctx 初始context
func OrderBy(ctx context.Context, field, order string) context.Context {
	return context.WithValue(ctx, orderHelperValue, &OrderByInfo{Field: field, Order: order})
}

//获得总记录数
func (p *PageInfo) GetTotal() int64 {
	return p.total
}

//设置总记录数，用户手动查询总记录数时使用
func (p *PageInfo) SetTotal(total int64) *PageInfo {
	p.total = total
	return p
}

//获得当前页码
func (p *PageInfo) GetPageNum() int64 {
	return p.Page
}

//获得每页的记录数
func (p *PageInfo) GetPageSize() int64 {
	return p.PageSize
}

//获得总页码
func (p *PageInfo) GetTotalPage() int64 {
	if p.PageSize <= 0 {
		return 0
	}

	return (p.total + int64(p.PageSize) - 1) / int64(p.PageSize)
}

func (p *PageInfo) String() string {
	return fmt.Sprintf("pageNum: %d , pageSize: %d , Total: %d , TotalPage: %d", p.GetPageNum(), p.GetPageSize(), p.GetTotal(), p.GetTotalPage())
}

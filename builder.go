/*
 * Licensed to the AcmeStack under one or more contributor license
 * agreements. See the NOTICE file distributed with this work for
 * additional information regarding copyright ownership.
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

import "context"

type builder struct {
	page  PageInfo
	order OrderByInfo

	ctx context.Context
}

// 创建builder
// ctx 初始context
func C(ctx context.Context) *builder {
	return &builder{ctx: ctx}
}

// 分页
// page 页码
// pageSize 分页大小
func (b *builder) Page(page, pageSize int64) *builder {
	b.page.Page = page
	b.page.PageSize = pageSize
	b.page.total = 0
	return b
}

// 分页
// page 页码
// pageSize 分页大小
func (b *builder) PageWithCount(page, pageSize int64, countColumn string) *builder {
	b.page.Page = page
	b.page.PageSize = pageSize
	b.page.countColumn = countColumn
	b.page.total = -1
	return b
}

func (b *builder) Count(countColumn string) *builder {
	b.page.countColumn = countColumn
	b.page.total = -1
	return b
}

// 手动指定字段和排序
// field 字段
// order 排序 [ASC | DESC]
func (b *builder) OrderBy(field, order string) *builder {
	b.order.Field = field
	b.order.Order = order
	return b
}

// 升序（默认）
// field 字段
func (b *builder) ASC(field string) *builder {
	return b.OrderBy(field, ASC)
}

// 降序
// field 字段
func (b *builder) DESC(field string) *builder {
	return b.OrderBy(field, DESC)
}

// 获得含分页/排序信息的context
func (b *builder) Build() context.Context {
	if b.page.PageSize > 0 {
		if b.page.total != -1 {
			b.ctx = StartPage(b.ctx, b.page.Page, b.page.PageSize)
		} else {
			b.ctx = StartPageWithCount(b.ctx, b.page.Page, b.page.PageSize, b.page.countColumn)
		}
	}

	if b.order.Field != "" {
		if b.order.Order == "" {
			b.order.Order = ASC
		}

		b.ctx = OrderBy(b.ctx, b.order.Field, b.order.Order)
	}

	return b.ctx
}

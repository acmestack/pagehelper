// Copyright (C) 2019, Xiongfa Li.
// All right reserved.
// @author xiongfa.li
// @version V1.0
// Description: 

package test

import (
	"context"
	"testing"
	"time"

	"github.com/acmestack/pagehelper"
)

const (
	pageHelperValue  = "_page_helper_value"
	orderHelperValue = "_order_helper_value"
)

func TestBuilder(t *testing.T) {
	ctx := pagehelper.C(context.Background()).Page(1, 3).OrderBy("test", pagehelper.ASC).Build()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	p := ctx.Value(pageHelperValue)
	o := ctx.Value(orderHelperValue)

	printOrder(t, o)
	printPage(t, p)
}

func TestBuilder2(t *testing.T) {
	ctx := pagehelper.C(context.Background()).PageWithCount(1, 3, "").OrderBy("test", pagehelper.ASC).Build()
	ctx = pagehelper.C(ctx).DESC("new_field").Build()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	p := ctx.Value(pageHelperValue)
	o := ctx.Value(orderHelperValue)

	printOrder(t, o)
	printPage(t, p)
}

func TestBuilder3(t *testing.T) {
	ctx := pagehelper.C(context.Background()).Page(1, 3).OrderBy("test", pagehelper.ASC).Build()
	ctx = pagehelper.C(ctx).DESC("new_field").ASC("new_field2").Count("test").Page(2, 100).Build()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	p := ctx.Value(pageHelperValue)
	o := ctx.Value(orderHelperValue)

	printOrder(t, o)
	printPage(t, p)
}

func printPage(t *testing.T, p interface{}) {
	if p, ok := p.(*pagehelper.PageInfo); ok {
		t.Logf("page param: %d %d", p.Page, p.PageSize)
	} else {
		t.Fail()
	}
}

func printOrder(t *testing.T, p interface{}) {
	if p, ok := p.(*pagehelper.OrderByInfo); ok {
		t.Logf("order param: %s %s", p.Field, p.Order)
	} else {
		t.Fail()
	}
}

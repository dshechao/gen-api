// file: main_test.go
package main

import (
	"github.com/kataras/iris/v12"
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
)

func TestIrisGin(t *testing.T) {

	cc := iris.Map{
		"a":        "123",
		"b":        456.12,
		"c":        "dddd",
		"d":        "ewrwerewqr",
		"username": "这是一个测试案例",
		"name":     123,
	}
	comment := iris.Map{
		"d":        "ewrwerewqr",
		"username": "这是一个测试案例`用户名`",
		"name":     123,
		"coment":   cc,
	}
	app := newApp()
	e := httptest.New(t, app)
	e.GET("/json").WithHeader("CurrentApiName", "测试内容接口111--测试").
		WithHeader("CurrentApiComment", "备注内容").Expect().Status(httptest.StatusOK)

	e.POST("/hello").WithHeader("CurrentApiName", "测试内容接口222--测试").
		WithHeader("CurrentApiComment", "备注内容").
		WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTU3LCJwaG9uZSI6IjE4OTEyMzQ1Njc4IiwiZW5hYmxlIjoxLCJleHAiOjE1ODkzNTkyOTQsImlzcyI6ImJzdC1jb21tdW5pdHktaWp3dCJ9.m6z1rt_vYIxSlYwMMQBnTpjWNjjhzqiXBT7yxp_E7tc").
		WithFormField("username999", "kataras").Expect().Status(httptest.StatusOK)

	e.POST("/reqbody").WithHeader("CurrentApiName", "测试内容接口333--测试").
		WithHeader("CurrentApiComment", "ooo").
		WithJSON(comment).
		Expect().Status(httptest.StatusOK)
	// give time to "gen" to generate the doc, 5 seconds are more than enough
	time.Sleep(5 * time.Second)
}

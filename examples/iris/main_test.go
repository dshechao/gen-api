// file: main_test.go
package main

import (
	"github.com/kataras/iris/v12"
	"testing"
	"time"

	"github.com/kataras/iris/v12/httptest"
)

func TestIrisGin(t *testing.T) {
	comment := iris.Map{
		"a": "123`测试参数收代理费`", "b": "456`不要的参数`", "c": "dddd`滴啊点那个`", "d": "ewrwerewqr`最后一个参数`", "username": "这是一个测试案例`用户名`",
	}
	app := newApp()
	e := httptest.New(t, app)
	e.GET("/json").WithHeader("CurrentApiName", "测试内容接口111--测试").
		WithHeader("CurrentApiComment", "备注内容").Expect().Status(httptest.StatusOK)

	e.POST("/hello").WithHeader("CurrentApiName", "测试内容接口222--测试").
		WithHeader("CurrentApiComment", "备注内容").
		WithFormField("username999", "kataras").Expect().Status(httptest.StatusOK)

	e.POST("/reqbody").WithHeader("CurrentApiName", "测试内容接口333--测试").
		WithHeader("CurrentApiComment", "ooo").
		WithJSON(comment).
		//WithJSON(myModel{Username: "kataras`8989好`",Gender: "lsdjfls;adj`大连市附近的数量`"}).
		//WithJSON(myModel{Username: "kataras",Gender: "lsdjfls;adj"}).
		Expect().Status(httptest.StatusOK)
	//Body().Equal("kataras")

	// give time to "gen" to generate the doc, 5 seconds are more than enough
	time.Sleep(5 * time.Second)
}

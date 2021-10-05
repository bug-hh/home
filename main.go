package main

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	_ "home/routers"
	_ "home/models"
	"net/http"
	"strings"
)

func main() {
	ignoreStaticPath()
	beego.Run()
}

func ignoreStaticPath() {
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)

}

func TransparentStatic(ctx *context.Context) {
	orPath := ctx.Request.URL.Path
	fmt.Println("request url: ", orPath)
	// 如果路径中包含 api 字段，说明是指令，应该取消静态资源重定向
	if strings.Index(orPath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/" + ctx.Request.URL.Path)
}


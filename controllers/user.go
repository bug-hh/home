package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"home/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) RetData(resp map[string]interface{})  {
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *UserController) Reg() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	json.Unmarshal(this.Ctx.Input.RequestBody, &resp)

	fmt.Println("mobile: ", resp["mobile"])
	fmt.Println("password: ", resp["password"])
	fmt.Println("sms_code", resp["sms_code"])

	o := orm.NewOrm()
	user := models.User{}

	user.Name = resp["mobile"].(string) // 后面的 .(string) 是断言
	user.Password_hash = resp["password"].(string)
	user.Mobile = resp["mobile"].(string)

	id, err := o.Insert(&user)
	if err != nil {
		resp["errno"] = 4002
		resp["errmsg"] = "注册失败"
		return
	}

	this.SetSession("name", user.Name)
	fmt.Println("注册成功 id: ", id)
	resp["errno"] = 0
	resp["msg"] = "注册成功"
}

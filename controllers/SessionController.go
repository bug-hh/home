package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"home/models"
	"home/utils"
)

type SessionController struct {
	beego.Controller
}

func (this *SessionController) RetData(resp map[string]interface{})  {
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *SessionController) GetSessionData() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	user := &models.User{}
	//resp["errno"] = utils.RECODE_DATAERR
	//resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
	name := this.GetSession("name")
	if name != nil {
		user.Name = name.(string)
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		resp["data"] = user
	}
}

func (this *SessionController) DeleteSessionData() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	this.DelSession("name")
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
}

func (this *SessionController) Login() {
	// 获取用户信息
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	json_err := json.Unmarshal(this.Ctx.Input.RequestBody, &resp)
	if json_err != nil{
		fmt.Println("json_err: json_err")
		return
	}

	// 判断是否合法
	if resp["mobile"] == nil || resp["password"] == nil {
		fmt.Println("判断是否合法")
		fmt.Println(resp)
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}
	// 与数据库匹配判断密码是否正确
	var query_user models.User
	o := orm.NewOrm()
	err := o.QueryTable("user").Filter("mobile", resp["mobile"]).One(&query_user)
	if err == orm.ErrMultiRows {
		// 多条的时候报错
		fmt.Printf("Returned Multi Rows Not One")
	}
	if err == orm.ErrNoRows {
		// 没有找到记录
		fmt.Printf("Not row found")
	}

	if query_user.Password_hash != resp["password"] {
		fmt.Println("密码错误")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = "密码错误"
		return
	}

	// 添加 session
	this.SetSession("name", resp["mobile"])
	this.SetSession("mobile", resp["mobile"])
	this.SetSession("user_id", query_user.Id)

	// 返回 json 数据给前端
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
}
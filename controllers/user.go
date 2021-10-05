package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"home/models"
	"home/utils"
	"path"
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

func (this *UserController) UploadAvatar() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	fileData, header, err := this.GetFile("avatar")
	if err != nil {
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}
	suffix := path.Ext(header.Filename)
	upload_bytes := make([]byte, header.Size)
	if _, err := fileData.Read(upload_bytes); err != nil {
		logs.Error(err.Error())
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}

	file_id, upload_err := models.UploadBinary(upload_bytes, suffix[1:])
	if upload_err != nil {
		logs.Error(upload_err.Error())
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}

	// 从 session 中获取 name
	user_name := this.GetSession("name")
	logs.Info("avatar user_name: ", user_name)

	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	qs.Filter("name", user_name).One(&user)
	// 更新 avatar_url, 然后更新数据库
	user.Avatar_url = file_id

	_, update_err := o.Update(&user)
	if update_err != nil {
		logs.Error(update_err.Error())
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}
	logs.Info("数据库更新成功")
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = fmt.Sprintf("http://%s/%s", this.Ctx.Request.Host, file_id)
}

func (this *UserController) GetUserData() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)

	user_name := this.GetSession("name")
	var user models.User
	orm.NewOrm().QueryTable("user").Filter("name", user_name).One(&user)

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = &user
}

func (this *UserController) UpdateUserName() {
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	user_name := this.GetSession("name")
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	qs.Filter("name", user_name).One(&user)
	temp := make(map[string]interface{})
	if json_err := json.Unmarshal(this.Ctx.Input.RequestBody, &temp); json_err != nil {
		logs.Error(json_err.Error())
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}
	logs.Info("user name: ", temp["name"])
	user.Name = temp["name"].(string)
	logs.Info("user.Name", user.Name)

	_, err := o.QueryTable("user").Filter("name", user_name).Update(orm.Params{
		"name": temp["name"].(string),
	})
	if err != nil {
		logs.Error("更新用户名失败")
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	this.SetSession("name", temp["name"].(string))
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
}
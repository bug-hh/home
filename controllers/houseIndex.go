package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type HouseIndexController struct {
	beego.Controller
}

func (this *HouseIndexController) RetData(resp map[string]interface{})  {
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *HouseIndexController) GetHouseIndex() {
	resp := make(map[string]interface{})

	resp["errno"] = 4001
	resp["errmsg"] = "查询失败"
	this.RetData(resp)

}
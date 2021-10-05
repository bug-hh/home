package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"home/models"
	"home/utils"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) RetData(resp map[string]interface{})  {
	this.Data["json"] = &resp
	this.ServeJSON()
}

func (this *AreaController) GetArea() {
	fmt.Println("connect success")

	resp := make(map[string]interface{})
	defer this.RetData(resp)

	// 从缓存中拿数据


	var areas []models.Area

	o := orm.NewOrm()
	num, err := o.QueryTable("area").All(&areas)

	if err != nil {
		fmt.Printf("出错了")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}

	if num == 0 {
		fmt.Printf("没有数据")
		resp["errno"] = utils.RECODE_NODATA
		resp["errmsg"] = utils.RecodeText(utils.RECODE_NODATA)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas

	fmt.Println("query data success: ", resp, "num = ", num)
}

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/cache"
	_ "github.com/beego/beego/v2/client/cache/redis"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"home/models"
	"home/utils"
	"time"
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
	redis_conn, err := cache.NewCache("redis", `{"key":"home","conn":":6379","dbNum":"0"}`)
	if err != nil {
		logs.Error(err)
	}

	result, _ := redis_conn.Get(context.Background(), "areas")

	var areas []models.Area

	if result != nil {
		logs.Info("获取缓存成功 result: ", string(result.([]byte)))
		//logs.Info(areas)
		json.Unmarshal(result.([]byte), &areas)
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
		resp["data"] = areas
		logs.Info("resp: ", resp)
		return
	}



	logs.Info("尝试从数据库中获取")
	o := orm.NewOrm()
	num, err := o.QueryTable("area").All(&areas)

	if err != nil {
		fmt.Printf("出错了")
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		return
	}

	if num == 0 {
		logs.Info("没有数据")
		resp["errno"] = utils.RECODE_NODATA
		resp["errmsg"] = utils.RecodeText(utils.RECODE_NODATA)
		return
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas

	// 将数据存入缓存
	logs.Info("将 areas 数据存入缓存")
	json_str, json_err := json.Marshal(areas)
	if json_err != nil {
		logs.Error(json_err)
	}
	redis_conn.Put(context.Background(), "areas", json_str, time.Second * 3600)

	fmt.Println("query data success: ", resp, "num = ", num)
}

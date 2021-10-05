package routers

import (
	"home/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("api/v1.0/areas", &controllers.AreaController{}, "get:GetArea")
    beego.Router("api/v1.0/house/index", &controllers.HouseIndexController{}, "get:GetHouseIndex")
    beego.Router("api/v1.0/session", &controllers.SessionController{}, "get:GetSessionData;delete:DeleteSessionData")
    beego.Router("api/v1.0/users", &controllers.UserController{}, "post:Reg")
    beego.Router("/api/v1.0/sessions", &controllers.SessionController{}, "post:Login")
}

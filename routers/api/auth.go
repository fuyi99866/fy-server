package api

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"go_server/pkg/logger"
	"go_server/pkg/setting"
	"go_server/pkg/util"
	"go_server/service/robot_service"
	"net/http"
)

type ConnectProfile struct {
	MqttAddr      string
	MqttPort      int
	RobotUser     string
	RobotPassword string
	Tls           bool
	RobotClientID string
	RequestTopic  string
	ResponseTopic string
	NotifyTopic   string
}

// @Summary  检查token是否过期
// @Tags 鉴权
// @Accept json
// @Produce json
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /alive  [get]
// @Security ApiKeyAuth
func TokenAlive(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary   登录获取登录token 信息
// @Tags 鉴权
// @Accept json
// @Produce  json
// @Param   body  body   models.UserLogin   true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /auth  [POST]
func Auth(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.UserLogin
	if err := c.ShouldBindJSON(&reqInfo); err != nil {
		body, _ := c.GetRawData()
		logger.Info("Auth request: ", string(body))
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	logger.Info("Auth: ", reqInfo.Username, reqInfo.Password)

	//有效性验证
	valid := validation.Validation{}
	valid.AlphaNumeric(reqInfo.Username, "username").Message("非法用户名")
	valid.Required(reqInfo.Username, "User").Message("用户名不能为空")
	valid.MinSize(reqInfo.Username, 1, "User").Message("用户名最短为1位")
	valid.MaxSize(reqInfo.Password, 100, "username").Message("最长为100字符")
	valid.MaxSize(reqInfo.Password, 100, "password").Message("最长为100字符")
	valid.Required(reqInfo.Password, "Pwd").Message("密码不能为空")
	valid.MinSize(reqInfo.Password, 6, "Pwd").Message("密码最短为6位")
	//检测用户的账号密码，看能否正常登录
	isOk, err := models.CheckUser(reqInfo.Username, reqInfo.Password)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if !isOk {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	token, err := util.GenerateToken(reqInfo.Username, reqInfo.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

// @Summary   获取机器人连接服务器需要的信息
// @Tags   连接
// @Accept json
// @Produce  json
// @Param sn query string true "机器人SN"
// @Success 200 {string} json ConnectProfile
// @Failure 400 {object} app.Response
// @Router /profile [GET]
func GetConnectProfile(c *gin.Context) {
	appG := app.Gin{C: c}
	sn := c.Query("sn")
	//获取机器人的企业号
	logger.Info("sn: ", sn)
	robot, err := models.GetRobotInfoBySn(sn)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, ConnectProfile{
		MqttAddr:      setting.Mqtt.Addr,
		MqttPort:      setting.Mqtt.Port,
		RobotUser:     setting.Mqtt.UserName,
		RobotPassword: setting.Mqtt.Password,
		Tls:           setting.Mqtt.Tls,
		RobotClientID: sn,
		RequestTopic:  robot_service.MakeTopic(robot.Company,robot.SN,robot_service.ROBOT_REQUEST),
		ResponseTopic: robot_service.MakeTopic(robot.Company,robot.SN,robot_service.ROBOT_RESPONSE),
		NotifyTopic:   robot_service.MakeTopic(robot.Company,robot.SN,robot_service.ROBOT_NOTIFY),
	})
}

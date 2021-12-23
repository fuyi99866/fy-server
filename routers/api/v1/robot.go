package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_server/models"
	"go_server/pkg/app"
	"go_server/pkg/e"
	"net/http"
)

// @Summary   添加机器人
// @Tags   坝光酒店部署
// @Accept json
// @Produce  json
// @Param body body  models.RobotStatusEdit true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /robot/add [POST]
// @Security ApiKeyAuth
func AddRobotTech(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.RobotStatusEdit
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		return
	}
	robot := models.RobotStatusTech{
		SN:        reqInfo.SN,
		NickName:  reqInfo.RobotNo,
		Company:   reqInfo.RobotBranchCode,
		RobotType: reqInfo.RobotType,
		PointId:   reqInfo.PointId,
	}

	if err := models.EditRobotStatusTech(robot); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary   获取所有机器人的状态统计信息
// @Tags   科技化云平台
// @Accept json
// @Produce  json
// @Param body body  models.RobotTypeInfo true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /robot/statistics [POST]
// @Security ApiKeyAuth
func GetRobotsStat(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.RobotTypeInfo
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		logrus.Error("err = ", err, )
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	list, err := models.GetAllRobotStatistics(reqInfo.RobotType)
	data := make([]struct {
		RobotType     int `json:"type"`
		Total         int `json:"total"`
		OnlineCount   int `json:"onlineCount"`
		OfflineCount  int `json:"offlineCount"`
		WorkingCount  int `json:"workingCount"`
		ChargingCount int `json:"chargingCount"`
	}, 0)
	data = append(data, struct {
		RobotType     int `json:"type"`
		Total         int `json:"total"`
		OnlineCount   int `json:"onlineCount"`
		OfflineCount  int `json:"offlineCount"`
		WorkingCount  int `json:"workingCount"`
		ChargingCount int `json:"chargingCount"`
	}{
		RobotType:     reqInfo.RobotType,
		Total:         list[0],
		OnlineCount:   list[1],
		OfflineCount:  list[2],
		WorkingCount:  list[3],
		ChargingCount: list[4],
	})
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// @Summary   给楼层部署机器人，并绑定位置和地图
// @Tags  坝光酒店部署
// @Accept json
// @Produce  json
// @Param body body  models.RobotRoomTechInfo true "body"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /robot/point [POST]
// @Security ApiKeyAuth
func SetRobotPoint(c *gin.Context) {
	appG := app.Gin{C: c}
	var reqInfo models.RobotRoomTechInfo
	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		return
	}
	rp := models.RobotRoomTech{
		SN:         reqInfo.SN,
		BuildingNo: reqInfo.BuildingNo,
		FloorNo:    reqInfo.FloorNo,
		RoomNo:     reqInfo.RoomNo,
		MapName:    reqInfo.MapName,
		RoomID:     reqInfo.RoomID,
		AreaID:     reqInfo.AreaID,
	}
	if err = models.EditRobotPoint(rp); err != nil {
		logrus.Error("EditRobotPoint failed : ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	point := reqInfo.BuildingNo + "B-" + reqInfo.BuildingNo + "F"
	if err = models.UpdateRobotStatusTech(reqInfo.SN, point); err != nil {
		logrus.Error("UpdateRobotStatusTech failed : ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary   查询楼层机器人及位置地图信息
// @Tags  坝光酒店部署
// @Accept json
// @Produce  json
// @Param buildingNo query  string true "楼栋"
// @Param floorNo query  string true "楼层"
// @Param roomNo query  string false "房间"
// @Success 200 {string} json "{ "code": 200, "data": {}, "msg": "ok" }"
// @Failure 400 {object} app.Response
// @Router /robot/point [GET]
// @Security ApiKeyAuth
func GetRobotPoint(c *gin.Context) {
	appG := app.Gin{C: c}
	build := c.Query("buildingNo")
	floor := c.Query("floorNo")

	logrus.Info("build,floor ", build, floor)
	point, err := models.GetRobotPoint(build, floor)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, point)
}




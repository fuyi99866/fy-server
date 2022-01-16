package models

import (
	"github.com/jinzhu/gorm"
	"go_server/pkg/logger"
)

//机器人自身状态
type RobotTypeInfo struct {
	RobotType int `json:"type"`
}

//楼层信息
type RobotRoomTech struct {
	gorm.Model
	SN         string `json:"sn"`         //机器人唯一标识
	BuildingNo string `json:"buildingNo"` //所在楼栋编号
	FloorNo    string `json:"floorNo"`    //所在楼栋楼层
	RoomNo     string `json:"roomNo"`     //当前工作房间号码或者准备前往的房间号码（在工作状态下，房间号必填）
	MapName    string `json:"mapName"`    //地图名
	RoomID     int    `json:"roomId"`     //房间ID
	AreaID     int    `json:"areaId"`     //区域ID
}

//机器人自身状态
type RobotStatusTech struct {
	gorm.Model
	SN          string `json:"sn"`
	NickName    string `json:"nickName"`    //机器人编号 ；机器人的昵称
	Company     string `json:"company"`     //机器人厂商编号 ；企业号
	RobotType   int    `json:"robotType"`   //1餐饮 2垃圾回收 3环境消杀 4房间消杀 5巡逻监控 6物流运送
	RobotStatus int    `json:"robotStatus"` //机器人状态:   1 在线 ； 2 离线；3 工作中 ； 4 充电
	PointId     string `json:"pointId"`     //机器人的位置信息，事先进行分配
	Online      string `json:"online"`      //机器人在线离线状态 ： 1 在线 ； 2 离线
}

//任务信息
type RobotTaskTech struct {
	gorm.Model
	SN            string `json:"sn"`            //机器人唯一标识
	TaskStatus    string `json:"taskStatus"`    //1、未开始 ；2、进行中；3、已完成 ； 4、未完成 ；5、失败 ；6、暂停中
	TaskID        string `json:"taskID"`        //机器人平台内部任务ID
	DestinationId string `json:"destinationId"` //任务目标位置ID信息；
}

type RobotRoomTechInfo struct {
	SN         string `json:"sn"`         //机器人唯一标识
	BuildingNo string `json:"buildingNo"` //所在楼栋编号
	FloorNo    string `json:"floorNo"`    //所在楼栋楼层
	RoomNo     string `json:"roomNo"`     //当前工作房间号码或者准备前往的房间号码（在工作状态下，房间号必填）
	MapName    string `json:"mapName"`    //地图名
	RoomID     int    `json:"roomId"`     //房间ID
	AreaID     int    `json:"areaId"`     //区域ID
}


type TaskIDInfo struct {
	TaskId string `json:"taskId"`
}

type OrderInfo struct {
	Priority         string   `json:"priority"`
	VehicleType      string   `json:"vehicleType"`
	DestinationCodes []string `json:"destinationCodes"`
}

type OrderIdListInfo struct {
	OrderIdList []string `json:"orderIdList"`
}

type RobotStatusEdit struct {
	SN              string `json:"sn"`
	RobotNo         string `json:"robotNo"`         //机器人编号 ；机器人的昵称
	RobotBranchCode string `json:"robotBranchCode"` //机器人厂商编号 ；企业号
	RobotType       int    `json:"robotType"`       //1餐饮 2垃圾回收 3环境消杀 4房间消杀 5巡逻监控 6物流运送
	PointId         string `json:"pointId"`         //机器人的部署的楼层位置
}

//获取楼层消杀机器人的自身状态
func GetRobotStatusByPointId(pintId string) (*RobotStatusTech, error) {
	var robot RobotStatusTech
	err := db.Where("point_id = ? ", pintId).First(&robot).Error

	if err != nil {
		return nil, err
	}
	return &robot, nil
}

//获取机器人的统计信息
func GetAllRobotStatistics(robotType int) ([]int, error) {
	var (
		total         int
		onlineCount   int
		offlineCount  int
		workingCount  int
		chargingCount int
	)
	list := make([]int, 0)
	if err := db.Model(&RobotStatusTech{}).Where("robot_type = ?", robotType).Count(&total).Error; err != nil {
		return nil, err
	}
	list = append(list, total)
	if err := db.Model(&RobotStatusTech{}).Where("robot_type = ? AND online = ?", robotType, "true").Count(&onlineCount).Error; err != nil {
		return nil, err
	}
	list = append(list, onlineCount)
	if err := db.Model(&RobotStatusTech{}).Where("robot_type = ? AND online = ?", robotType, "false").Count(&chargingCount).Error; err != nil {
		return nil, err
	}
	list = append(list, offlineCount)
	if err := db.Model(&RobotStatusTech{}).Where("robot_type = ? AND robot_status = ?", robotType, 3).Count(&offlineCount).Error; err != nil {
		return nil, err
	}

	list = append(list, workingCount)
	if err := db.Model(&RobotStatusTech{}).Where("robot_type = ? AND robot_status = ?", robotType, 4).Count(&workingCount).Error; err != nil {
		return nil, err
	}
	list = append(list, chargingCount)

	return list, nil
}

//获取消杀机器人的任务状态
func GetRobotTaskStatusByTaskId(taskId string) (interface{}, error) {
	//联合查询，获取机器人的三个表的数据 join
	return nil, nil
}

//查询机器人列表
func GetAllRobotList(page, pageSize int) ([]*RobotStatusTech, int, error) {
	var robots []*RobotStatusTech
	var total int
	err := db.Model(&RobotStatusTech{}).Limit(pageSize).Offset((page - 1) * pageSize).Find(&robots).Error
	db.Model(RobotStatusTech{}).Count(&total)
	if err != nil {
		return nil, 0, err
	}
	return robots, total, nil
}

//通过SN查询机器人
func GetRobotInfoBySn(sn string) (*RobotStatusTech, error) {
	var robot RobotStatusTech
	err := db.Model(&RobotStatusTech{}).Where("sn = ?", sn).First(&robot).Error
	if err != nil {
		return nil, err
	}
	return &robot, nil
}

//根据任务id，查询机器人
func GetRobotTaskByTaskId(taskId string) (*RobotTaskTech, error) {
	var task RobotTaskTech
	err := db.Model(&RobotTaskTech{}).Where("task_id = ?", taskId).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

//根据SN查询房间信息
func GetRoomInfoBySn(sn string) (*RobotRoomTech, error) {
	var room RobotRoomTech
	err := db.Model(&RobotRoomTech{}).Where("sn = ?", sn).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func IsExistTaskByTaskId(taskId string) (bool, error) {
	var task RobotTaskTech
	err := db.Select("id").Where("task_id = ?", taskId).First(&task).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if task.ID > 0 {
		return true, nil
	}
	return false, nil
}

//楼层是否存在
func IsExistTechRoom(buildNo, floorNo string) (bool, error) {
	var room RobotRoomTech
	err := db.Select("id").Where("building_no = ? AND floor_no = ? ", buildNo, floorNo).First(&room).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if room.ID > 0 {
		return true, nil
	}
	return false, nil
}

func IsExistRobotTech(sn string) (bool, error) {
	var robot RobotStatusTech
	err := db.Select("id").Where("sn = ?", sn).First(&robot).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if robot.ID > 0 {
		return true, nil
	}
	return false, nil
}

//创建现有地图的任务
func CreateTask(sn, taskId, destination string) error {
	task := RobotTaskTech{
		SN:            sn,
		TaskStatus:    "1",
		TaskID:        taskId,
		DestinationId: destination,
	}
	err := db.Create(&task).Error
	return err
}

//更新任务
func UpdateTask(taskId string, data RobotTaskTech) error {
	exist, err := IsExistTaskByTaskId(taskId)
	if err != nil || !exist {
		return err
	}
	err = db.Where("task_id = ?", taskId).First(&RobotTaskTech{}).Update(&data).Error
	return err
}

func GetTaskByTaskId(taskId string) (*RobotTaskTech, error) {
	var task RobotTaskTech
	exist, err := IsExistTaskByTaskId(taskId)
	if err != nil || !exist {
		return nil, err
	}
	err = db.Where("task_id = ?", taskId).First(&task).Error
	return &task, err
}

// 更新楼层机器人位置和地图信息
func EditRobotPoint(data RobotRoomTech) error {
	var err error
	exist, err := IsExistTechRoom(data.BuildingNo, data.FloorNo)
	if err != nil {
		return err
	}

	logger.Info("exist = ", exist)
	if !exist {
		err = db.Model(&RobotRoomTech{}).Create(&data).Error
		return err
	} else {
		err = db.Where("building_no = ? AND floor_no = ? ", data.BuildingNo, data.FloorNo).First(&RobotRoomTech{}).Update(&data).Error
		return err
	}
}

func GetRobotPoint(buildNo, floorNo string) (*RobotRoomTech, error) {
	var point RobotRoomTech
	err := db.Where("building_no = ? AND floor_no = ? ", buildNo, floorNo).First(&point).Error
	logger.Info("GetRobotPoint err ", err)
	if err != nil {
		return nil, err
	}
	return &point, nil
}

func GetRobotPointBySn(sn string) (*RobotRoomTech, error) {
	var point RobotRoomTech
	err := db.Where("sn = ?", sn).First(&point).Error

	if err != nil {
		return nil, err
	}
	return &point, nil
}

//更新机器人部署楼层
func UpdateRobotStatusTech(sn, pointId string) error {
	robot := RobotStatusTech{
		SN:      sn,
		PointId: pointId,
	}
	err := db.Where("sn = ?", sn).First(&RobotStatusTech{}).Update(&robot).Error
	return err
}

//更新机器人在线离线状态
func UpdateRobotOnlineStatusTech(sn string, online bool) error {
	var status string
	if online {
		status = "true"
	} else {
		status = "false"
	}
	robot := RobotStatusTech{
		SN:     sn,
		Online: status,
	}
	err := db.Where("sn = ?", sn).First(&RobotStatusTech{}).Update(&robot).Error
	return err
}

//更新机器人的状态
func UpdateRobotStatusTechBySn(sn string, status int) error {
	robot := RobotStatusTech{
		SN:          sn,
		RobotStatus: status,
	}
	err := db.Where("sn = ?", sn).First(&RobotStatusTech{}).Update(&robot).Error
	return err
}

func EditRobotStatusTech(data RobotStatusTech) error {
	var err error
	exist, err := IsExistRobotTech(data.SN)
	if err != nil {
		return err
	}
	if !exist {
		err = db.Create(&data).Error
		return err
	} else {
		err = db.Where("sn = ?", data.SN).First(&RobotStatusTech{}).Update(&data).Error
		return err
	}
}

func GetAllRobot() ([]*RobotStatusTech, error) {
	var robots []*RobotStatusTech
	err := db.Find(&robots).Error
	if err != nil {
		return nil, err
	}

	return robots, nil
}

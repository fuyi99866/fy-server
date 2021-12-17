package format

const (
	TaskStartRequestTitle     = "request_add_task"    //创建任务
	TaskTerminateRequestTitle = "request_task_cancel" //中断任务
)

//task 创建
type MqttTaskAddRequest struct {
	Title   string `json:"title"`
	Content struct {
		SessionID         string `json:"sessionid"`
		Timestamp         int64  `json:"timestamp"`
		Name              string `json:"name"`                //任务唯一标识
		NickName          string `json:"nickname"`            //任务别名（可忽略
		UvType            string `json:"type"`                //定点消杀：disinfect、自主消杀：explore_disinfect
		Cycle             string `json:"cycle"`               //immediately：立即执行 manual：手动触发 once：单次预约执行 day：每天执行 week：每周执行 month：每月执行
		Flags             int    `json:"flags"`               //门禁牌(1) 、PIR(2) 、摄像头AI(4），支持组合，例如 1|2|4
		KillDuration      int64  `json:"duration"`            //杀毒总时间，单位毫秒
		CountDown         int    `json:"countdown"`           //倒计时，单位秒
		PointList         []int  `json:"pointlist"`           //位置点列表
		Mode              int    `json:"mode"`                //低功率(1)、全功率(2)
		Date              string `json:"date"`                //【暂无】定时消杀日期。2021-03-20
		Time              string `json:"time"`                //【暂无】定时消杀时间。“21:35:00”
		Days              []int  `json:"days"`                //【暂无】每个星期消杀日期
		RoomID            int    `json:"room_id"`             //房间ID
		AreaID            int    `json:"area_id"`             //区域ID
		RoomName          string `json:"room_name"`           //房间名
		MapName           string `json:"map_name"`            //地图名
		AutoFinishTimeout int64  `json:"auto_finish_timeout"` //自动恢复时间
		Config            struct {
			Distance int `json:"distance"`
			Fov      int `json:"fov"`
			Area     int `json:"area"`
			Mode     int `json:"mode"`
			GridSize int `json:"grid_size"`
		}
	} `json:"content"`
}

type MqttTaskFinishRequest struct {
	Title   string `json:"title"`
	Content struct {
		SessionID string `json:"sessionid"`
		Timestamp int64  `json:"timestamp"`
		Name      string `json:"name"`
	} `json:"content"`
}

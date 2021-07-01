package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"go_server/conf"
	"go_server/routers/casbin/enforcer"
)

var db *gorm.DB

type Company struct {
	gorm.Model
	CompanyID string
}

type Robot struct {
	RobotSN      string `json:"robotsn"`
	TaskID       string `json:"taskid";gorm:"unique_index, not null"` //任务名称（也就是ID，唯一标识
	NickName     string `json:"nickname"`                             //任务别名（可忽略
	UvType       string `json:"type"`                                 //定点消杀：disinfect、自主消杀：explore_disinfect
	Cycle        string `json:"cycle"`                                //immediately：立即执行 manual：手动触发 once：单次预约执行 day：每天执行 week：每周执行 month：每月执行
	Flags        int    `json:"flags"`                                //门禁牌(1) 、PIR(2) 、摄像头AI(4），支持组合，例如 1|2|4
	KillDuration int64  `json:"duration"`                             //杀毒总时间，单位毫秒
	CountDown    int    `json:"countdown"`                            //倒计时，单位秒
	PointList    string `json:"pointlist"`                            //位置点列表
	Mode         int    `json:"mode"`                                 //低功率(1)、全功率(2)
	Date         string `json:"data"`                                 //【暂无】定时消杀日期。2021-03-20
	Time         string `json:"time"`                                 //【暂无】定时消杀时间。“21:35:00”
	Days         string `json:"days"`                                 //【暂无】每个星期消杀日期
	Status       string `json:"status"`                               //消杀状态 start:开始 working：消杀中 terminate_by_robot：硬件结束 terminate_by_software软件结束 finish：消杀完成
}

type Authority struct {
	gorm.Model
	AuthorityId   string `json:"authority_id" gorm:"not null;unique"` //权限ID
	AuthorityName string `json:"authority_name"`                      //角色名
	ParentId      string `json:"parent_id"`                           //父角色ID
}

type User struct {
	gorm.Model
	Username  string  `json:"username" gorm:"column:username; index:usr_name;not null"` //唯一索引
	Password  string  `json:"password, omitempty" gorm:"column:password;" `
	NickName  string  `json:"nickname" gorm:"column:nickname;" `
	CompanyID string  `json:"company_id" gorm:"column:company_id;"`
	Robots    []Robot `gorm:"many2many:user_robot;"`
}

type UserRegister struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CompanyID string `json:"company_id"`
	NickName  string `json:"nickname"`
}

//用户登录
type UserLogin struct {
	//10003 账号不存在
	//20001 30002 登录异常
	//400 参数错误
	//40002 用户不存在
	//40003 账号密码错误
	Username string `json:"username"`
	Password string `json:"password"`
}

//用户权限
type UserPolicy struct {
	Username string `json:"username"`
	URL      string `json:"url"`
	Type     string `json:"type"`
}

func Init() {
	//连接数据库
	var err error
	db, err = gorm.Open(conf.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.DatabaseSetting.User,
		conf.DatabaseSetting.Password,
		conf.DatabaseSetting.Host,
		conf.DatabaseSetting.Name))
	if err != nil {
		logrus.Fatal("无法连接数据库... err: %/v", err)
	}

	//指定表的前缀，修改默认的表名
	/*	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return conf.DatabaseSetting.TablePrefix+defaultTableName
	}*/
	//禁用表名复数
	db.SingularTable(true)       //设置全局表名禁用复数
	db.DB().SetMaxOpenConns(100) //设置最大连接数
	db.DB().SetMaxIdleConns(10)  //设置最大闲置连接数

	migration()

	//检查管理员账号是否存在，不在就添加管理员账号
	exist, _ := CheckUser("admin", "123456", )
	if !exist {
		data := map[string]interface{}{
			"username": "admin",
			"password": "123456",
		}
		AddUser(data) //添加的管理员账号不允许修改
	}

	//给管理员赋予改变访问权限的权限

	e := enforcer.EnforcerTool()
	e.AddPolicy("admin", "/policy", "GET")
	e.AddPolicy("admin", "/policy", "POST")
	e.AddPolicy("admin", "/policy", "DELETE")
}

//自动创建数据表
func migration() {
	db.AutoMigrate(&Company{}).AutoMigrate(&Robot{}).AutoMigrate(&User{}).AutoMigrate(&Authority{})
}

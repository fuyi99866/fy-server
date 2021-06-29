package global

import (
	"github.com/jinzhu/gorm"
	"go_server/routers/config"

)

var(
	GVA_DB  *gorm.DB
	//GVA_REDIS *redis.Client
	//GVA_CONFIG config.Server
	//MyLogger *myLogger
	//GVA_VP *viper.Viper
	//GVA_LOG *zap.Logger
	GVA_CONFIG config.Server

)

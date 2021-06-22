package global

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var(
	GVA_DB  *gorm.DB
	//GVA_REDIS *redis.Client
	//GVA_CONFIG config.Server
	GVA_VP *viper.Viper
	GVA_LOG *zap.Logger
)

package service

import (
	"go_server/pkg/logger"
	"go_server/service/adis_server"
	"go_server/service/robot_service"
)

/**
启动服务
 */
type Service interface {
	Start() error
	Stop() error
}

func Start() error {
	logger.Info("robot service start ...")
	robot_service.S.Start()
	adis_server.A.Start()
	return nil
}

func Stop() error {
	robot_service.S.Stop()
	return nil
}

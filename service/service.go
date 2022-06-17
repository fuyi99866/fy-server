package service

import (
	"github.com/sirupsen/logrus"
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
	logrus.Info("robot service start ...")
	//robot_service.S.Start()
	//adis_server.A.Start()
	//adis_server.TestStart()
	return nil
}

func Stop() error {
	robot_service.S.Stop()
	return nil
}

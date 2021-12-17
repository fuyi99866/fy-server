package service

import "go_server/service/robot_service"

/**
启动服务
 */
type Service interface {
	Start() error
	Stop() error
}

func Start() error {
	robot_service.S.Start()
	return nil
}

func Stop() error {
	robot_service.S.Stop()
	return nil
}

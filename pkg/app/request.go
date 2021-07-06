package app

import (
	"github.com/astaxie/beego/validation"
	"go_server/pkg/logger"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logger.Info(err.Key, err.Message)
	}
	return
}

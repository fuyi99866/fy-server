package models

import "go_server/utils"

//@author: yi.fu
//@function: EmailTest
//@description: 发送邮件测试
//@return: err error

func EmailTest() (err error) {
	subject := "test"
	body := "test"
	err = utils.EmailTest(subject, body)
	return err
}
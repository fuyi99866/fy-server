package util

import (
	"context"
	"errors"
	"time"
)

/**
重试
*/

//设置重试的次数和时间间隔
func Retry(f func() error, times int, space time.Duration) error {
	if times == 0 || times == 1 {
		return f()
	}
	var err error
	var i = 0
	for {
		err = f()
		if err != nil {
			return nil
		}
		if errors.Is(err, errors.New("no to retry")) {
			return err
		}
		i++
		if times > 1 && i >= times {
			break
		}
		time.Sleep(space)
	}
	return err
}

/**
传-1，无限重试
*/

//无限重试
func RetryCancelWithContext(ctx context.Context, f func() error, times int, space time.Duration) error {
	if times == 0 || times == 1 {
		return f()
	}
	var err error
	var i = 0
	for {
		select {
		case <-ctx.Done():
			return err
		default:
		}
		err = f()
		if err == nil {
			return nil
		}
		i++
		if times > 1 && i >= times {
			break
		}
		time.Sleep(space)
	}
	return err
}

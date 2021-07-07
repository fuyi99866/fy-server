package main

import (
	"github.com/robfig/cron"
	"go_server/models"
	"go_server/pkg/logger"
	"log"
	"time"
)

/**
定时任务
*/

func CronRun() {
	logger.Info("start corn")
	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		logger.Info("Run models.CleanAllTag...")
		models.CleanAllTag()
	})

	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()
	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}

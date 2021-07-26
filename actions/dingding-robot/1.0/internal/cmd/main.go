package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dingding-robot/1.0/internal/dingdingRobot"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)

	logrus.Info("dingding-robot Testing...")
	if err := dingdingRobot.Run(); err != nil {
		logrus.Errorf("dingding-robot"+
			" failed, err: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

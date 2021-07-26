package dingdingRobot

import (
	"github.com/erda-project/erda-actions/actions/dingding-robot/1.0/internal/conf"
)

func Run() error {
	if err := conf.Load(); err != nil {
		return err
	}
	return handleAPIs()
}

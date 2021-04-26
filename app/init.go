package app

import (
	"github.com/flashguru-git/node-monitor/config"
)

var (
	cfg config.Provider
)

func init() {
	cfg = config.Config()
}

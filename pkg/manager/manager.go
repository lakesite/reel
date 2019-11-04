package manager

import (
	"github.com/lakesite/ls-config/pkg/config"
	"github.com/lakesite/ls-fibre/pkg/service"
)

func RunManagementService() {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ws := service.NewWebService("reel", address)
	ws.RunWebServer()
}

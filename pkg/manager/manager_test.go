package manager

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/lakesite/ls-config"
	"github.com/lakesite/ls-fibre"
)

var (
	cwd_arg = flag.String("cwd", "", "set cwd")
)

func init() {
	flag.Parse()
	if *cwd_arg != "" {
		if err := os.Chdir(*cwd_arg); err != nil {
			fmt.Println("Chdir error:", err)
		}
	}
}

func TestManagementService(t *testing.T) {
	address := config.Getenv("REEL_HOST", "127.0.0.1") + ":" + config.Getenv("REEL_PORT", "7999")
	ws := fibre.NewService("reel", address)

	if ws == nil {
		t.Errorf("Reel management service initialization failed")
	}
	// todo.
}

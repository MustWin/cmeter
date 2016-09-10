package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/MustWin/cmeter/cmd"
	_ "github.com/MustWin/cmeter/cmd/agent"
	_ "github.com/MustWin/cmeter/cmd/api"
	"github.com/MustWin/cmeter/cmd/root"
	"github.com/MustWin/cmeter/context"
	_ "github.com/MustWin/cmeter/monitor/cadvisor"
	_ "github.com/MustWin/cmeter/monitor/embedded"
)

const VERSION = "0.0.1-alpha"

func main() {
	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), VERSION)

	execute := cmd.BuildRootExecutor(ctx, root.Info)
	if err := execute(); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/MustWin/cmeter/cmd"
	_ "github.com/MustWin/cmeter/cmd/api"
	"github.com/MustWin/cmeter/context"
)

const VERSION = "0.0.1-alpha"

func main() {
	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), VERSION)

	execute := cmd.BuildRootExecutor(ctx, rootCmdInfo)
	if err := execute(); err != nil {
		log.Fatalln(err)
	}
}

var (
	rootCmdInfo = &cmd.Info{
		Use:   "cmeter",
		Short: "`cmeter`",
		Long:  "`cmeter`",
	}
)

package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/MustWin/cmeter/cmd"
	_ "github.com/MustWin/cmeter/cmd/agent"
	_ "github.com/MustWin/cmeter/cmd/api"
	"github.com/MustWin/cmeter/cmd/root"
	_ "github.com/MustWin/cmeter/containers/embedded"
	"github.com/MustWin/cmeter/context"
)

const VERSION = "0.0.1-alpha"

func main() {
	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), VERSION)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/MustWin/cmeter/cmd"
	_ "github.com/MustWin/cmeter/cmd/agent"
	"github.com/MustWin/cmeter/cmd/root"
	_ "github.com/MustWin/cmeter/cmd/version"
	_ "github.com/MustWin/cmeter/containers/embedded"
	"github.com/MustWin/cmeter/context"
	_ "github.com/MustWin/cmeter/reporting/cmeterapi"
	_ "github.com/MustWin/cmeter/reporting/http"
	_ "github.com/MustWin/cmeter/reporting/mock"
)

var appVersion string

const DEFAULT_VERSION = "0.0.0-dev"

func main() {
	if appVersion == "" {
		appVersion = DEFAULT_VERSION
	}

	rand.Seed(time.Now().Unix())
	ctx := context.WithVersion(context.Background(), appVersion)

	dispatch := cmd.CreateDispatcher(ctx, root.Info)
	if err := dispatch(); err != nil {
		log.Fatalln(err)
	}
}

package api

import (
	"github.com/MustWin/cmeter/cmd"
	"github.com/MustWin/cmeter/configuration"
	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/mockapi"
)

func init() {
	cmd.Register("api", Info)
}

func run(ctx context.Context, args []string) error {
	config, err := configuration.Resolve(args)
	if err != nil {
		return err
	}

	server, err := mockapi.NewServer(ctx, config)
	if err != nil {
		return err
	}

	return server.ListenAndServe()
}

var (
	Info = &cmd.Info{
		Use:   "api",
		Short: "`api`",
		Long:  "`api`",
		Run:   cmd.ExecutorFunc(run),
	}
)

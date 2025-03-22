package daemon

import (
	"context"
	"net/http"
	"sync"

	"github.com/jantytgat/go-kit/pkg/application"
	"github.com/spf13/cobra"

	"github.com/jantytgat/rpdemo/internal/handlers"
)

const (
	daemonCmdUse   = "daemon"
	daemonCmdShort = "Run daemon"
	daemonCmdLong  = "Run daemon"
)

var Cmd = application.Command{
	Command: &cobra.Command{
		Use:   daemonCmdUse,
		Short: daemonCmdShort,
		Long:  daemonCmdLong,
		RunE:  daemonRunE,
	},
	SubCommands: nil,
	Configure:   nil,
}

func daemonRunE(cmd *cobra.Command, args []string) error {
	mux := http.NewServeMux() // Create sample handler to returns 404
	mux.Handle("/", handlers.NewRootHandler())
	srv := application.NewHttpServer("0.0.0.0", 28080, mux)
	srvCtx, srvCancel := context.WithCancel(cmd.Context())
	defer srvCancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		srv.Run(ctx)
	}(srvCtx, &wg)

	wg.Wait()
	return nil
}

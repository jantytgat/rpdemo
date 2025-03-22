package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jantytgat/go-kit/pkg/application"
	"github.com/jantytgat/go-kit/pkg/semver"
	"github.com/jantytgat/go-kit/pkg/slogd"
)

var (
	version string

	appName   = "rpdemo"
	appTitle  = "Reverse Proxy Demo"
	appBanner = ""
)

var (
	ctx        context.Context
	semVersion semver.Version
)

func main() {
	var err error

	// Configure logging
	slogd.RegisterColoredTextHandler(os.Stdout, true)
	ctx = slogd.WithContext(context.Background())

	if semVersion, err = semver.Parse(version); err != nil {
		slogd.Logger().LogAttrs(ctx, slogd.LevelError, "error running application", slog.Any("error", err))
	}

	application.New(appName, appTitle, appBanner, semVersion)
	// application.RegisterCommands(subCommands, nil)

	if err = application.Run(ctx); err != nil {
		slogd.Logger().LogAttrs(ctx, slogd.LevelError, "error running application", slog.Any("error", err))
		os.Exit(1)
	}
	os.Exit(0)
}

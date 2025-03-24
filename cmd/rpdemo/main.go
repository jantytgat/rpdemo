package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jantytgat/go-kit/pkg/slogd_colored"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/jantytgat/go-kit/pkg/application"
	"github.com/jantytgat/go-kit/pkg/semver"
	"github.com/jantytgat/go-kit/pkg/slogd"

	"github.com/jantytgat/rpdemo/internal/handlers"
)

var (
	version string = "0.0.1"

	appName   = "rpdemo"
	appTitle  = "Reverse Proxy Demo"
	appBanner = ""
)

var (
	ctx        context.Context
	semVersion semver.Version

	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to STING_NUMBER.
	envPrefix = "RPDEMO"
	vConfig   *viper.Viper

	defaultConfigFilename = "rpdemo"
	// Replace hyphenated flag names with camelCase in the config file
	replaceHyphenWithCamelCase = false

	listenPortFlag int
	baseFqdnFlag   string
	colorFlag      string
)

func main() {
	var err error

	// Configure logging
	slogd.Init(slogd.LevelInfo, false)
	slogd_colored.RegisterColoredTextHandler(os.Stdout, true)
	ctx = slogd.WithContext(context.Background())

	if semVersion, err = semver.Parse(version); err != nil {
		slogd.Logger().LogAttrs(ctx, slogd.LevelError, "error running application", slog.Any("error", err))
		os.Exit(1)
	}

	application.New(appName, appTitle, appBanner, semVersion)
	application.RegisterFlag(configureFlags)
	application.OverrideRunE(runE)
	application.RegisterPersistentPreRunE(loadEnv)
	if err = application.Run(ctx); err != nil {
		slogd.Logger().LogAttrs(ctx, slogd.LevelError, "error running application", slog.Any("error", err))
		os.Exit(1)
	}
	os.Exit(0)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		// If using camelCase in the config file, replace hyphens with a camelCased string.
		// Since viper does case-insensitive comparisons, we don't need to bother fixing the case, and only need to remove the hyphens.
		if replaceHyphenWithCamelCase {
			configName = strings.ReplaceAll(f.Name, "-", "")
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func configureFlags(cmd *cobra.Command) {
	// var err error
	cmd.PersistentFlags().IntVarP(&listenPortFlag, "port", "p", 28080, "port to listen on")
	cmd.PersistentFlags().StringVarP(&baseFqdnFlag, "base-fqdn", "", "", "base fqdn to use for the application")

	// Set color page
	cmd.PersistentFlags().StringVarP(&colorFlag, "color-page", "", "white", "set page color (red, green, blue)")
}

func loadEnv(cmd *cobra.Command, args []string) error {
	slogd.FromContext(cmd.Context()).LogAttrs(cmd.Context(), slogd.LevelDebug, "loading environment variables")

	vConfig = viper.New()
	vConfig.SetConfigName(defaultConfigFilename)
	vConfig.SetEnvPrefix(envPrefix)
	vConfig.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	vConfig.AutomaticEnv()

	bindFlags(cmd, vConfig)
	return nil
}

func runE(cmd *cobra.Command, args []string) error {
	mux := http.NewServeMux() // Create sample handler to returns 404

	mux.Handle("/", handlers.NewRootHandler(colorFlag))
	srv := application.NewHttpServer("0.0.0.0", listenPortFlag, mux)
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

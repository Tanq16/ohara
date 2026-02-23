package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/tanq16/ohara/internal/server"
	"github.com/tanq16/ohara/internal/store"
)

var AppVersion = "dev-build"

var debugFlag bool

var serveFlags struct {
	dataDir string
	port    int
}

var rootCmd = &cobra.Command{
	Use:     "ohara",
	Short:   "Track professional achievements and touchpoints",
	Version: AppVersion,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Run: runServe,
}

func init() {
	cobra.OnInitialize(setupLogs)
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	rootCmd.Flags().StringVar(&serveFlags.dataDir, "data-dir", "./data", "Path to data directory")
	rootCmd.Flags().IntVar(&serveFlags.port, "port", 8080, "Server port")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func setupLogs() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.DateTime,
		NoColor:    false,
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runServe(cmd *cobra.Command, args []string) {
	st, err := store.New(store.Config{DataDir: serveFlags.dataDir})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize store")
	}

	srv := server.New(server.Config{Port: serveFlags.port}, st)

	log.Info().
		Str("package", "cmd").
		Int("port", serveFlags.port).
		Str("data", serveFlags.dataDir).
		Msg("Starting Ohara")

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}

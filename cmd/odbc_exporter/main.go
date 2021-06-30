package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/magna-z/prometheus-odbc-exporter/internal/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Version string

var landingPageFmt = `<html>
<head><title>ODBC Exporter</title></head>
<body><h1>ODBC Exporter</h1><p><a href="%s">Metrics</a></p></body>
</html>`

var (
	metricsAddr = flag.String(
		"metrics.addr",
		":9798",
		"Address and port for exporter",
	)
	metricsPath = flag.String(
		"metrics.path",
		"/metrics",
		"URL path for public collected metrics",
	)
	metricsConfig = flag.String(
		"metrics.config",
		"/etc/odbc-exporter/config.yml",
		"Path to config file",
	)
	metricsDSN = flag.String(
		"metrics.dsn",
		"",
		"DSN connection string",
	)
	metricsLowerCase = flag.Bool(
		"metrics.lowercase",
		true,
		"Lower case for all published metrics and labels keys",
	)
	merticsPrefix = flag.String(
		"metrics.prefix",
		"odbc",
		"",
	)

	logLevel = flag.String(
		"log.level",
		"info",
		"Log level (trace|debug|info|warn|error|fatal|panic)",
	)

	version = flag.Bool(
		"version",
		false,
		"Print version and exit",
	)
)

func init() {
	flag.Parse()

	if *version {
		fmt.Printf("odbc-exporter version %s\n", Version)
		os.Exit(0)
	}
}

func httpServer(wg *sync.WaitGroup) *http.Server {
	s := &http.Server{Addr: *metricsAddr}

	var landingPage = []byte(fmt.Sprintf(landingPageFmt, *metricsAddr))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(landingPage)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Landing page write error")
		}
	})
	http.Handle(*metricsPath, promhttp.Handler())
	go func() {
		defer wg.Done()

		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Str("metrics.addr", *metricsAddr).
				Msg("HTTP-server serve error")
		}
	}()

	log.Info().Msgf("Exporter started on %q", *metricsAddr)

	return s
}

func main() {
	var err error

	zLevel, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(zLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DisableSampling(true)

	conf, err := config.GetConfig(*metricsConfig)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("metrics.config", *metricsConfig).
			Msg("Config file read error")
	}

	// Global Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Global WaitGroup
	wg := &sync.WaitGroup{}
	wg.Add(1)
	httpSrv := httpServer(wg)

	// Termination
	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	sig := <-termSign
	log.Info().Str("signal", sig.String()).Msg("Shutdown exporter")

	err = httpSrv.Shutdown(ctx)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("metrics.addr", *metricsAddr).
			Msg("HTTP-server shutdown error")
	}
}

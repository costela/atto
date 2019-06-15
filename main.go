package main // import "github.com/costela/atto"

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nytimes/gziphandler"

	"github.com/sirupsen/logrus"
	"github.com/stevenroose/gonfig"
)

var version = "devel"

var logger *logrus.Logger

var conf = struct {
	LogLevel string `default:"warn" short:"l" desc:"level of logging output (one of debug/info/warn/error)"`
	Port     int    `default:"8080" desc:"port to listen on"`
	Path     string `default:"/www" desc:"path to serve"`
	Prefix   string `default:"" desc:"prefix under which atto will be accessed (this will be stripped before accessing 'path')"`
	ShowList bool   `default:"false" desc:"whether to display folder contents"`
	Compress bool   `default:"true" desc:"whether to transparently compress served files"`
	Timeout  struct {
		ReadHeader *duration `default:"5s" desc:"time to wait for request headers"`
		Shutdown   *duration `default:"30s" desc:"time to wait for ungoing requests to finish before shutting down"`
	}
}{}

func main() {
	logger = logrus.New()

	if err := gonfig.Load(&conf, gonfig.Conf{
		EnvPrefix: "ATTO_",
	}); err != nil {
		logger.Fatalf("could not load config: %s", err)
	}

	if level, err := logrus.ParseLevel(conf.LogLevel); err == nil {
		logger.SetLevel(level)
	} else {
		logger.Fatal(err)
	}
	log.SetOutput(logger.Writer())

	logger.WithField("version", version).Debugf("starting atto")

	logger.Debugf("instantiating server for path %s", conf.Path)
	handler := http.StripPrefix(conf.Prefix, http.FileServer(safeDir(conf.Path)))
	if conf.Compress {
		handler = gziphandler.GzipHandler(handler)
	}
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", conf.Port),
		Handler:           logHandler{handler},
		ReadHeaderTimeout: time.Duration(*conf.Timeout.ReadHeader),
		ReadTimeout:       time.Duration(*conf.Timeout.ReadHeader), // we do not expect any content upload, so headers are enough
		ErrorLog:          log.New(logger.WithField("source", "http.Server").WriterLevel(logrus.WarnLevel), "", 0),
	}

	server.RegisterOnShutdown(func() {
		logger.Debugf("shutting down gracefully (timeout: %s)", time.Duration(*conf.Timeout.Shutdown))
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Debugf("listening on port %d", conf.Port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSignals(&server)
	}()

	wg.Wait()
}

type logHandler struct {
	inner http.Handler
}

func (lh logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{"host": r.Host, "path": r.URL.Path, "method": r.Method}).Debug("got request")
	lh.inner.ServeHTTP(w, r)
}

// work around time.Duration's lack of UnmarshalText
type duration time.Duration

func (d *duration) UnmarshalText(data []byte) error {
	dd, err := time.ParseDuration(string(data))
	*d = duration(dd)
	return err
}

func handleSignals(server *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*conf.Timeout.Shutdown))
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		if err == context.DeadlineExceeded {
			logger.Warnf("timeout exceeded while shutting down")
		} else {
			logger.Errorf("error while shutting down server: %s", err)
		}
	}
}

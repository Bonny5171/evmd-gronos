package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/logger"
	"github.com/besser/cron"
	"github.com/pkg/errors"

	"bitbucket.org/everymind/evmd-gronos/core"
)

func init() {
	// Setting the limits the number of operating system threads that can execute user-level Go code simultaneously.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setting the log output and prefix
	logger.Init("", os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	// Starting flags
	StartFlags()

	if Version {
		fmt.Printf("Version: %s\n", VERSION)
		os.Exit(0)
	}

	logger.Tracef("-> Starting gronos service version %s (%s)", VERSION, runtime.Version())

	os.Setenv("GOTRACE", strconv.FormatBool(Trace))

	if len(os.Getenv("GRONOS_DATABASE_CONFIG_DSN")) == 0 {
		logger.Fatalln("Environment variable 'GRONOS_DATABASE_CONFIG_DSN' not defined!")
	}

	if len(os.Getenv("GRONOS_DATABASE_DSN")) == 0 {
		logger.Fatalln("Environment variable 'GRONOS_DATABASE_DSN' not defined!")
	}

	if len(os.Getenv("GRONOS_STACK_NAME")) == 0 {
		logger.Fatalln("Environment variable 'GRONOS_STACK_NAME' not defined!")
	}

	if len(os.Getenv("GRONOS_SCHEDULE")) == 0 {
		os.Setenv("GRONOS_SCHEDULE", "@every 30s")
	}

	logger.Traceln("Openning conncetion with DBs...")

	// Getting DSN
	cfgDSN := os.Getenv("GRONOS_DATABASE_CONFIG_DSN")
	dtaDSN := os.Getenv("GRONOS_DATABASE_DSN")

	// DB conn variables
	var (
		dbMaxOpenConns int = 5
		dbMaxIdleConns int = 1
		dbMaxLifeTime  int = 10
	)

	if len(os.Getenv("GRONOS_DATABASE_MAXOPENCONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GRONOS_DATABASE_MAXOPENCONNS")); e != nil {
			dbMaxOpenConns = v
		}
	}

	if len(os.Getenv("GRONOS_DATABASE_MAXIDLECONNS")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GRONOS_DATABASE_MAXIDLECONNS")); e != nil {
			dbMaxIdleConns = v
		}
	}

	if len(os.Getenv("GRONOS_DATABASE_MAXLIFETIME")) > 0 {
		if v, e := strconv.Atoi(os.Getenv("GRONOS_DATABASE_MAXLIFETIME")); e != nil {
			dbMaxLifeTime = v
		}
	}

	// Starting config DB connection
	if err := db.Connections.Connect("CONFIG", &db.PostgresDB{
		ConnectionStr: cfgDSN,
		MaxOpenConns:  dbMaxOpenConns,
		MaxIdleConns:  dbMaxIdleConns,
		MaxLifetime:   dbMaxLifeTime,
	}); err != nil {
		logger.Infof("DSN: %s\n", cfgDSN)
		logger.Errorln(err)
	}

	// Starting config DB connection
	if err := db.Connections.Connect("DATA", &db.PostgresDB{
		ConnectionStr: dtaDSN,
		MaxLifetime:   dbMaxLifeTime,
		MaxIdleConns:  dbMaxIdleConns,
		MaxOpenConns:  dbMaxOpenConns,
	}); err != nil {
		logger.Infof("DSN: %s\n", dtaDSN)
		logger.Errorln(err)
	}

	logger.Traceln("Connected!")

	// Create a new cron manager
	c := cron.NewWithLocation(time.UTC)
	defer c.Stop()

	// Add func to cron
	if _, err := c.AddFuncN("client", os.Getenv("GRONOS_SCHEDULE"), func() {
		if err := startJob(c); err != nil {
			logger.Errorln(errors.Wrap(err, "startJob()"))
		}
	}); err != nil {
		logger.Fatalln(errors.Wrap(err, "c.AddFunc()"))
		return
	}

	logger.Tracef("Jobs cron verifications scheduled to: %s", os.Getenv("GRONOS_SCHEDULE"))

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	logger.Traceln("Waiting for job scheduled...")

	// Start cron
	c.Start()

	// Waiting for interrupt by system signal
	killSignal := <-interrupt
	logger.Infoln("Got signal:", killSignal)
}

func startJob(c *cron.Cron) error {
	// Job start here
	if err := core.Run(c); err != nil {
		return errors.Wrap(err, "core.Run()")
	}
	return nil
}

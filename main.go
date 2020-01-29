package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3" // "github.com/besser/cron"

	"bitbucket.org/everymind/evmd-golib/db"
	"bitbucket.org/everymind/evmd-golib/logger"
	"bitbucket.org/everymind/evmd-gronos/cmd"
	"bitbucket.org/everymind/evmd-gronos/core"
)

var version string

func init() {
	// Setting the limits the number of operating system threads that can execute user-level Go code simultaneously.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Setting the log output and prefix
	logger.Init("", os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr)
}

func main() {
	// Starting flags
	cmd.StartFlags()

	if cmd.BuildVersion {
		fmt.Print(version)
		os.Exit(0)
	}

	if cmd.Version {
		fmt.Printf("Version: %s (%s)\n", version, runtime.Version())
		os.Exit(0)
	}

	logger.Tracef("-> Starting gronos service version %s (%s)", version, runtime.Version())

	os.Setenv("GOTRACE", strconv.FormatBool(cmd.Trace))

	if len(os.Getenv("GRONOS_DATABASE_DSN")) == 0 {
		logger.Fatalln("Environment variable 'GRONOS_DATABASE_DSN' not defined!")
	}

	if len(os.Getenv("GRONOS_SCHEDULE")) == 0 {
		os.Setenv("GRONOS_SCHEDULE", "@every 30s")
	}

	logger.Traceln("Openning conncetion with DBs...")

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

	// Getting DSN
	dsn := os.Getenv("GRONOS_DATABASE_DSN")

	// Starting config DB connection
	if err := db.Connections.Connect("CONFIG", &db.PostgresDB{
		ConnectionStr: dsn,
		MaxOpenConns:  dbMaxOpenConns,
		MaxIdleConns:  dbMaxIdleConns,
		MaxLifetime:   dbMaxLifeTime,
	}); err != nil {
		logger.Infof("DSN: %s", dsn)
		logger.Errorln(err)
	}

	logger.Traceln("Connected!")

	// Create a new cron manager
	loc := time.UTC
	if locationStr := os.Getenv("LOCATION"); len(locationStr) > 0 {
		if l, err := time.LoadLocation(locationStr); err == nil {
			loc = l
		}
	}
	c := cron.New(cron.WithLocation(loc))
	defer c.Stop()

	scheduledJobs := make(map[string]core.ScheduledJob, 0)

	// Add func to cron
	if _, err := c.AddFunc(os.Getenv("GRONOS_SCHEDULE"), func() {
		if err := startJob(c, scheduledJobs); err != nil {
			logger.Errorln(fmt.Errorf("startJob(): %w", err))
		}
	}); err != nil {
		logger.Fatalln(fmt.Errorf("c.AddFunc(): %w", err))
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

func startJob(c *cron.Cron, sJobs map[string]core.ScheduledJob) error {
	// Job start here
	if err := core.Run(c, sJobs); err != nil {
		return fmt.Errorf("core.Run(): %w", err)
	}
	return nil
}

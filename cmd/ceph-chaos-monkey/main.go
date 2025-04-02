package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	cephShellDriver "github.com/teran/ceph-chaos-monkey/ceph/drivers/shell"
	"github.com/teran/ceph-chaos-monkey/monkey"
	"github.com/teran/go-collection/random"
)

type config struct {
	LogLevel log.Level `envconfig:"LOG_LEVEL" default:"error"`
}

func main() {
	cfg := config{}
	envconfig.MustProcess("", &cfg)

	log.SetLevel(cfg.LogLevel)

	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s FUSS_INTERVAL GAME_DURATION\n", os.Args[0])
		fmt.Println("Both FUSS_INTERVAL and GAME_DURATION are in seconds")
		os.Exit(1)
	}

	interval, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Incorrect interval value, must be integer: %s", os.Args[1])
		os.Exit(1)
	}

	if interval < 45 {
		fmt.Println("Please set interval >=45s to allow Ceph to react and give you a proper status.")
		os.Exit(1)
	}

	duration, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		fmt.Printf("Incorrect duration value, must be integer: %s", os.Args[1])
		os.Exit(1)
	}

	if duration > 60*60 {
		fmt.Println("Please avoid duration >=1h to have some time to analyze what you've done and what actually happened")
		os.Exit(1)
	}

	ctx := context.TODO()

	runner := cephShellDriver.NewRunner("/usr/bin/ceph", "/usr/bin/rados")
	cluster := cephShellDriver.New(runner)
	printer := monkey.NewPrinter()

	m := monkey.New(cluster, random.GetRand(), printer, time.Duration(interval)*time.Second, time.Duration(duration)*time.Second)
	if err := m.Run(ctx); err != nil {
		panic(err)
	}
}

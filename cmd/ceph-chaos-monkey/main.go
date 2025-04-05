package main

import (
	"context"
	"fmt"
	"os"

	kingpin "github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
	"github.com/teran/go-collection/random"

	cephShellDriver "github.com/teran/ceph-chaos-monkey/ceph/drivers/shell"
	"github.com/teran/ceph-chaos-monkey/monkey"
)

const (
	appName = "ceph-chaos-monkey"

	runCmd     = "run"
	versionCmd = "version"
)

var (
	appVersion     = "n/a (dev build)"
	buildTimestamp = "undefined"

	app = kingpin.New(appName, "Ceph Chaos Monkey")

	isTrace = app.
		Flag("trace", "set verbosity level to trace").
		Bool()

	cephBinaryPath = app.
			Flag("ceph-binary", "path to the ceph binary").
			Default("/usr/bin/ceph").
			String()

	radosBinaryPath = app.
			Flag("rados-binary", "path to the rados binary").
			Default("/usr/bin/rados").
			String()

	isRun        = app.Command(runCmd, "run the game")
	fussInterval = isRun.
			Flag("fuss-interval", "set fuss interval i.e. how often to trigger chaos behavior. Example: 2m for 2 minutes").
			Required().
			Duration()

	gameDuration = isRun.
			Flag("game-duration", "set game duration i.e. overall time for chaos monkey to destroy Ceph cluster. Example 10m for 10 minutes").
			Required().
			Duration()

	_ = app.Command(versionCmd, "print version and exit")
)

func main() {
	ctx := context.TODO()
	appCmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	if isTrace != nil && *isTrace {
		log.SetLevel(log.TraceLevel)
	}

	switch appCmd {
	case runCmd:
		runner := cephShellDriver.NewRunner(*cephBinaryPath, *radosBinaryPath)
		cluster := cephShellDriver.New(runner)
		printer := monkey.NewPrinter()

		m := monkey.New(cluster, random.GetRand(), printer, *fussInterval, *gameDuration)
		if err := m.Run(ctx); err != nil {
			panic(err)
		}
		return
	case versionCmd:
		fmt.Printf("%s v%s (built @ %s)\n", appName, appVersion, buildTimestamp)
		os.Exit(1)
	}
}

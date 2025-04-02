package monkey

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teran/go-collection/random"

	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

type Monkey interface {
	Run(ctx context.Context) error
}

type monkey struct {
	cluster  drivers.Cluster
	duration time.Duration
	interval time.Duration
	printer  Printer
	rnd      random.Random
}

func New(cluster drivers.Cluster, rnd random.Random, printer Printer, interval time.Duration, duration time.Duration) Monkey {
	return &monkey{
		cluster:  cluster,
		duration: duration,
		interval: interval,
		printer:  printer,
		rnd:      rnd,
	}
}

func (m *monkey) Run(ctx context.Context) error {
	m.printer.Println(`This software is designed to train Ceph engineers to recover Ceph clusters in
various ways by interacting with Ceph components and data to trigger errors
in the cluster. Therefore it could damage the data stored within the cluster
and that's why there are some limitations where you can run ceph-chaos-monkey:

* >0 && <=10 OSD daemons
* <=500 GB of raw space

These restrictions are hardcoded and cannot be changed in runtime but anyway
if you have such a small clusters with important data please check twice where
you're running ceph-chaos-monkey.`)
	m.printer.Println()

	questions := []string{
		"Your Ceph cluster could be permanently damaged, are you sure you want to proceed?",
		"The data in your Ceph cluster could be permanently lost, are you still sure to proceed?",
		"You can just damage the cluster, lost data and not learn anything, still proceed?",
		"This mean you could never recover the data from the cluster, proceed?",
		"Check once again what I just asked you, sure?",
		"Sure???",
	}

	m.printer.Println()

	for _, q := range questions {
		if !askForConfirmation(q) {
			m.printer.Println()
			m.printer.Println("Ain't brave enough for this? No worries, get back later")

			return nil
		}
	}

	m.printer.Println()

	if ok := m.preflightCheck(ctx); !ok {
		return nil
	}

	m.printer.Printf(
		"Huh... that's what you wanted, let's go! Waiting %d seconds for the first action ...\n",
		int(m.interval.Seconds()),
	)

	ctx, cancel := context.WithTimeout(ctx, m.duration)
	defer cancel()

	ticker := time.NewTicker(m.interval)
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != context.DeadlineExceeded {
				return err
			}

			m.printer.Println("Game is over! Go check your cluster if it's still alive :-)")

			return nil
		case <-ticker.C:
			m.printer.Println("Tick! Running something dangerous in the cluster ...")
			if err := m.doSomeFuss(ctx); err != nil {
				if err != context.DeadlineExceeded {
					log.Debugf("error doSomeFuss(): %s", err)
					m.printer.Println("Hm... I'm starting getting errors from cluster, I'm leading! :-)")
					continue
				}

				m.printer.Println("Game is over! Go check your cluster if it's still alive :-)")

				return nil
			}
		}
	}
}

func (m *monkey) doSomeFuss(ctx context.Context) error {
	cases := []func(context.Context, drivers.Cluster, random.Random) error{
		setRandomFlag, unsetRandomFlag,
		destroyRandomOSD,
		randomlyResizeRandomPool, randomlyChangePGNumForRandomPool,
		reweightByUtilization,
		createPoolAndPutAmountOfObjects,
		setRandomNearFullRatio, setRandomBackfillfullRatio, setRandomFullRatio,
	}

	fn := cases[m.rnd.Intn(len(cases))]

	return fn(ctx, m.cluster, m.rnd)
}

func (m *monkey) preflightCheck(ctx context.Context) bool {
	osds, err := m.cluster.GetOSDs(ctx)
	if err != nil {
		m.printer.Println("Can't do a preflight check, sorry ...")
		return false
	}

	if len(osds) == 0 || len(osds) > 10 {
		m.printer.Printf("OSDs count must be >0 && <=10, you have: %d\n", len(osds))
		return false
	}

	var total uint64
	for _, osd := range osds {
		total += (osd.KbUsed + osd.KbAvailable)
	}

	if total > 500*1024*1024*1024 {
		m.printer.Printf("Total cluster space must be <=500GB, you have: %d bytes\n", total)
		return false
	}

	return true
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		}
	}
}

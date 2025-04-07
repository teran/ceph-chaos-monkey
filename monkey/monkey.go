package monkey

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/teran/go-collection/random"
	"golang.org/x/sync/errgroup"

	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

type Monkey interface {
	Run(ctx context.Context) error
}

type JournalEntry struct {
	Timestamp time.Time
	Entry     string
}

type monkey struct {
	cluster      drivers.Cluster
	duration     time.Duration
	interval     time.Duration
	printer      Printer
	stats        Stats
	rnd          random.Random
	bgIOPoolName string

	journal []JournalEntry
}

type fuss struct {
	name string
	fn   func(context.Context, drivers.Cluster, random.Random) error
}

func New(cluster drivers.Cluster, rnd random.Random, printer Printer, stats Stats, interval time.Duration, duration time.Duration) Monkey {
	return &monkey{
		cluster:      cluster,
		duration:     duration,
		interval:     interval,
		printer:      printer,
		rnd:          rnd,
		stats:        stats,
		bgIOPoolName: fmt.Sprintf("chaos-monkey-%d", rnd.Uint32()*rnd.Uint32()),
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

	go func(ctx context.Context) { _ = m.doBackgroundIO(ctx) }(ctx)

	ticker := time.NewTicker(m.interval)

outer:
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != context.DeadlineExceeded {
				return err
			}

			break outer
		case <-ticker.C:
			m.printer.Println("Tick! Running something dangerous in the cluster ...")
			if err := m.doSomeFuss(ctx); err != nil {
				if err != context.DeadlineExceeded {
					log.Debugf("error doSomeFuss(): %s", err)
					continue
				}

				break outer
			}
		}
	}

	m.printer.Println()
	m.printer.Println("Game is over! Go check your cluster if it's still alive :-)")
	m.printer.Println()

	s := m.stats.Dump()
	m.printer.Printf("Avg Reads latency = %.3fs\n", s.AvgReadsLatency.Seconds())
	m.printer.Printf("Avg Writes latency = %.3fs\n", s.AvgWritesLatency.Seconds())
	m.printer.Printf("Read operations succeeded = %.2f%%\n", s.ReadsSuccessPercent*100)
	m.printer.Printf("Write operations succeeded = %.2f%%\n", s.ReadsSuccessPercent*100)

	m.printer.Println()
	m.printer.Println("Here's the journal of your adventure during the game:")
	for _, j := range m.journal {
		fmt.Printf("- %s: %s\n", j.Timestamp.Format(time.RFC3339), j.Entry)
	}

	return nil
}

func (m *monkey) doSomeFuss(ctx context.Context) error {
	cases := []fuss{
		{
			name: "set random flag",
			fn:   setRandomFlag,
		},
		{
			name: "unset random flag",
			fn:   unsetRandomFlag,
		},
		{
			name: "destroy random OSD",
			fn:   destroyRandomOSD,
		},
		{
			name: "randomly resize random pool",
			fn:   randomlyResizeRandomPool,
		},
		{
			name: "randomly change pg_num for random pool",
			fn:   randomlyChangePGNumForRandomPool,
		},
		{
			name: "run reweight-by-utilization",
			fn:   reweightByUtilization,
		},
		{
			name: "set random value for nearfull-ratio",
			fn:   setRandomNearFullRatio,
		},
		{
			name: "set random value for backfillfull-ratio",
			fn:   setRandomBackfillfullRatio,
		},
		{
			name: "set random value for full-ratio",
			fn:   setRandomFullRatio,
		},
		{
			name: "remove random monitor",
			fn:   removeRandomMonitor,
		},
		{
			name: "drain random host",
			fn:   drainRandomHost,
		},
	}

	c := cases[m.rnd.Intn(len(cases))]

	m.journal = append(m.journal, JournalEntry{
		Timestamp: time.Now(),
		Entry:     c.name,
	})

	err := c.fn(ctx, m.cluster, m.rnd)
	if err != context.DeadlineExceeded {
		m.journal = append(m.journal, JournalEntry{
			Timestamp: time.Now(),
			Entry:     fmt.Sprintf("cluster operations are failing (during %s)", c.name),
		})
	}

	return c.fn(ctx, m.cluster, m.rnd)
}

func (m *monkey) doBackgroundIO(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return m.doBackgroundIOReads(ctx)
	})

	g.Go(func() error {
		return m.doBackgroundIOWrites(ctx)
	})

	return g.Wait()
}

func (m *monkey) doBackgroundIOReads(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != context.DeadlineExceeded {
				return err
			}
			return nil
		default:
			pools, err := m.cluster.GetPools(ctx)
			if err != nil {
				log.Debug("error getting pool list")
				continue
			}

			isPoolExists := false
			for _, p := range pools {
				if p.PoolName == m.bgIOPoolName {
					isPoolExists = true
				}
			}

			if !isPoolExists {
				log.Tracef("pool %s is not exists yet", m.bgIOPoolName)
				continue
			}

			start := time.Now()
			objs, err := m.cluster.ListRADOSObjects(ctx, m.bgIOPoolName)
			if err != nil {
				m.stats.ObserveRead(time.Since(start), err)
				continue
			}

			if len(objs) == 0 {
				continue
			}

			obj := objs[m.rnd.Intn(len(objs))]

			data, err := m.cluster.ReadRADOSObject(ctx, m.bgIOPoolName, obj)
			if err != nil {
				m.stats.ObserveRead(time.Since(start), err)
				continue
			}

			hasher := sha256.New()
			if _, err := hasher.Write(data); err != nil {
				return err
			}

			if obj != hex.EncodeToString(hasher.Sum(nil)) {
				m.stats.ObserveRead(time.Since(start), err)
				continue
			}
			m.stats.ObserveRead(time.Since(start), nil)
		}
	}
}

func (m *monkey) doBackgroundIOWrites(ctx context.Context) error {
	pools, err := m.cluster.GetPools(ctx)
	if err != nil {
		return err
	}

	isPoolExists := false
	for _, v := range pools {
		if v.PoolName == m.bgIOPoolName {
			isPoolExists = true
			break
		}
	}

	if !isPoolExists {
		if err := m.cluster.CreateDefaultPool(ctx, m.bgIOPoolName); err != nil {
			return err
		}
	}

	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != context.DeadlineExceeded {
				return err
			}
			return nil
		default:
			buf := make([]byte, m.rnd.Intn(1024*1024*1024))
			if _, err := m.rnd.Read(buf); err != nil {
				return err
			}

			hasher := sha256.New()
			if _, err := hasher.Write(buf); err != nil {
				return err
			}

			start := time.Now()
			err = m.cluster.CreateRADOSObject(ctx, m.bgIOPoolName, hex.EncodeToString(hasher.Sum(nil)), buf)
			m.stats.ObserveWrite(time.Since(start), err)
		}
	}
}

func (m *monkey) preflightCheck(ctx context.Context) bool {
	health, err := m.cluster.GetHealth(ctx)
	if err != nil {
		m.printer.Println("Can't do a preflight check, sorry ...")
		return false
	}

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

	m.journal = append(m.journal, JournalEntry{
		Timestamp: time.Now(),
		Entry: fmt.Sprintf(
			"your Ceph cluster is up and running in %s state with %d OSDs and %d bytes total raw space",
			health.Status, len(osds), total,
		),
	})

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

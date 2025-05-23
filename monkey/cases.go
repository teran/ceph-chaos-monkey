package monkey

import (
	"context"
	"errors"
	"strconv"

	"github.com/teran/go-collection/random"

	"github.com/teran/ceph-chaos-monkey/ceph"
	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

var cephFlags = []ceph.Flag{
	ceph.FlagNoBackfill,
	ceph.FlagNoDeepScrub,
	ceph.FlagNoIn,
	ceph.FlagNoOut,
	ceph.FlagNoRebalance,
	ceph.FlagNoRecover,
	ceph.FlagNoScrub,
	ceph.FlagNoUp,
}

func setRandomFlag(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	return c.SetFlag(ctx, cephFlags[rnd.Intn(len(cephFlags))])
}

func unsetRandomFlag(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	return c.UnsetFlag(ctx, cephFlags[rnd.Intn(len(cephFlags))])
}

func destroyRandomOSD(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	ids, err := c.GetOSDIDs(ctx)
	if err != nil {
		return err
	}

	if len(ids) == 0 {
		return errors.New("no OSDs are present in the cluster")
	}

	id := ids[rnd.Intn(len(ids))]

	return c.DestroyOSD(ctx, id)
}

func randomlyResizeRandomPool(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	pools, err := c.GetPools(ctx)
	if err != nil {
		return err
	}

	if len(pools) == 0 {
		return errors.New("no Pools are present in the cluster")
	}

	pool := pools[rnd.Intn(len(pools))]

	return c.ResizePool(ctx, pool.PoolName, uint64(rnd.Intn(10)))
}

func randomlyChangePGNumForRandomPool(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	pools, err := c.GetPools(ctx)
	if err != nil {
		return err
	}

	if len(pools) == 0 {
		return errors.New("no Pools are present in the cluster")
	}

	pool := pools[rnd.Intn(len(pools))]

	pgNumMax := 256
	if pool.Options.PgNumMax > 0 {
		pgNumMax = pool.Options.PgNumMax
	}

	return c.ChangePoolPGNum(ctx, pool.PoolName, uint64(rnd.Intn(pgNumMax)+1))
}

func reweightByUtilization(ctx context.Context, c drivers.Cluster, _ random.Random) error {
	return c.ReweightByUtilization(ctx)
}

func setRandomNearFullRatio(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	return c.SetNearFullRatio(ctx, rnd.Float64())
}

func setRandomBackfillfullRatio(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	return c.SetNearFullRatio(ctx, rnd.Float64())
}

func setRandomFullRatio(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	return c.SetNearFullRatio(ctx, rnd.Float64())
}

func removeRandomMonitor(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	mons, err := c.GetMons(ctx)
	if err != nil {
		return err
	}

	mon := mons[rnd.Intn(len(mons))]

	return c.RemoveMonitor(ctx, mon.Name)
}

func drainRandomHost(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	hosts, err := c.ListHosts(ctx)
	if err != nil {
		return err
	}

	host := hosts[rnd.Intn(len(hosts))]

	return c.DrainHost(ctx, host.Hostname)
}

func setRandomFlagForRandomGroup(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	targets := []string{}

	osdIDs, err := c.GetOSDIDs(ctx)
	if err != nil {
		return err
	}

	for _, o := range osdIDs {
		targets = append(targets, "osd."+strconv.FormatUint(o, 10))
	}

	hosts, err := c.ListHosts(ctx)
	if err != nil {
		return err
	}

	for _, h := range hosts {
		targets = append(targets, h.Hostname)
	}

	return c.SetGroupFlag(ctx, cephFlags[rnd.Intn(len(cephFlags))], targets[:int(len(targets)/3)]...)
}

func unsetRandomFlagFromRandomGroup(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	targets := []string{}

	osdIDs, err := c.GetOSDIDs(ctx)
	if err != nil {
		return err
	}

	for _, o := range osdIDs {
		targets = append(targets, "osd."+strconv.FormatUint(o, 10))
	}

	hosts, err := c.ListHosts(ctx)
	if err != nil {
		return err
	}

	for _, h := range hosts {
		targets = append(targets, h.Hostname)
	}

	return c.UnsetGroupFlag(ctx, cephFlags[rnd.Intn(len(cephFlags))], targets[:int(len(targets)/3)]...)
}

func deepScrubRandomPG(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	pgs, err := c.ListPGs(ctx)
	if err != nil {
		return err
	}

	pg := pgs[rnd.Intn(len(pgs))]

	return c.DeepScrubPG(ctx, pg.PGID)
}

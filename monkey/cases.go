package monkey

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"time"

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
	ceph.FlagPause,
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

	return c.ChangePoolPGNum(ctx, pool.PoolName, uint64(rnd.Intn(256)))
}

func reweightByUtilization(ctx context.Context, c drivers.Cluster, _ random.Random) error {
	return c.ReweightByUtilization(ctx)
}

func createPoolAndPutAmountOfObjects(ctx context.Context, c drivers.Cluster, rnd random.Random) error {
	poolName := "test-pool-" + strconv.FormatInt(rnd.Int63n(math.MaxInt64), 10)
	if err := c.CreateDefaultPool(ctx, poolName); err != nil {
		return err
	}

	maxSize := 150 * 1024 * 1024
	amount := rnd.Intn(50)

	for i := 0; i < amount; i++ {
		buf := make([]byte, rnd.Intn(maxSize))
		if _, err := rnd.Read(buf); err != nil {
			return err
		}

		hasher := sha256.New()
		if _, err := hasher.Write(buf); err != nil {
			return err
		}

		if err := c.CreateRADOSObject(ctx, poolName, hex.EncodeToString(hasher.Sum(nil)), buf); err != nil {
			time.Sleep(2 * time.Second)
		}
	}
	return nil
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

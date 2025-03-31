package monkey

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/teran/ceph-chaos-monkey/ceph"
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

func setRandomFlag(ctx context.Context, c ceph.Cluster) error {
	return c.SetFlag(ctx, cephFlags[getRandomChoice(len(cephFlags))])
}

func unsetRandomFlag(ctx context.Context, c ceph.Cluster) error {
	return c.UnsetFlag(ctx, cephFlags[getRandomChoice(len(cephFlags))])
}

func destroyRandomOSD(ctx context.Context, c ceph.Cluster) error {
	osds, err := c.GetOSDs(ctx)
	if err != nil {
		return err
	}

	if len(osds) == 0 {
		return errors.New("no OSDs are present in the cluster")
	}

	osd := osds[getRandomChoice(len(osds))]

	return c.DestroyOSD(ctx, osd.ID)
}

func randomlyResizeRandomPool(ctx context.Context, c ceph.Cluster) error {
	pools, err := c.GetPools(ctx)
	if err != nil {
		return err
	}

	if len(pools) == 0 {
		return errors.New("no Pools are present in the cluster")
	}

	pool := pools[getRandomChoice(len(pools))]

	return c.ResizePool(ctx, pool.PoolName, uint64(getRandomChoice(10)))
}

func randomlyChangePGNumForRandomPool(ctx context.Context, c ceph.Cluster) error {
	pools, err := c.GetPools(ctx)
	if err != nil {
		return err
	}

	if len(pools) == 0 {
		return errors.New("no Pools are present in the cluster")
	}

	pool := pools[getRandomChoice(len(pools))]

	return c.ChangePoolPGNum(ctx, pool.PoolName, uint64(2^getRandomChoice(16)))
}

func reweightByUtilization(ctx context.Context, c ceph.Cluster) error {
	return c.ReweightByUtilization(ctx)
}

func createPoolAndPutAmountOfObjects(ctx context.Context, c ceph.Cluster) error {
	poolName := "test-pool-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := c.CreateDefaultPool(ctx, poolName); err != nil {
		return err
	}

	maxSize := 150 * 1024 * 1024
	amount := getRandomChoice(50)

	for i := 0; i < amount; i++ {
		buf := make([]byte, getRandomChoice(maxSize))
		if _, err := rng.Read(buf); err != nil {
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

func setRandomNearFullRatio(ctx context.Context, c ceph.Cluster) error {
	return c.SetNearFullRatio(ctx, float64(getRandomChoice(100))/100.0)
}

func setRandomBackfillfullRatio(ctx context.Context, c ceph.Cluster) error {
	return c.SetNearFullRatio(ctx, float64(getRandomChoice(100))/100.0)
}

func setRandomFullRatio(ctx context.Context, c ceph.Cluster) error {
	return c.SetNearFullRatio(ctx, float64(getRandomChoice(100))/100.0)
}

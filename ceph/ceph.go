package ceph

import (
	"context"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Flag string

const (
	FlagNoOut       Flag = "noout"
	FlagNoIn        Flag = "noin"
	FlagNoRecover   Flag = "norecover"
	FlagNoScrub     Flag = "noscrub"
	FlagNoDeepScrub Flag = "nodeep-scrub"
	FlagPause       Flag = "pause"
	FlagNoBackfill  Flag = "nobackfill"
	FlagNoUp        Flag = "noup"
	FlagNoRebalance Flag = "norebalance"
)

type Cluster interface {
	GetOSDs(ctx context.Context) ([]OSD, error)
	GetOSDIDs(ctx context.Context) ([]uint64, error)
	GetMons(ctx context.Context) ([]Mon, error)

	DestroyOSD(ctx context.Context, id uint64) error
	StopOSDDaemon(ctx context.Context, id uint64) error

	SetFlag(ctx context.Context, flag Flag) error
	UnsetFlag(ctx context.Context, flag Flag) error

	GetPools(ctx context.Context) ([]Pool, error)
	CreateDefaultPool(ctx context.Context, name string) error
	ResizePool(ctx context.Context, name string, size uint64) error
	ChangePoolPGNum(ctx context.Context, name string, pgs uint64) error
	ReweightByUtilization(ctx context.Context) error

	CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error

	SetNearFullRatio(ctx context.Context, value float64) error
	SetBackfillfullRatio(ctx context.Context, value float64) error
	SetFullRatio(ctx context.Context, value float64) error
}

type cluster struct {
	runner Runner
}

func New(runner Runner) Cluster {
	return &cluster{
		runner: runner,
	}
}

func (c *cluster) GetOSDs(ctx context.Context) ([]OSD, error) {
	type osds struct {
		OSDs []OSD `json:"OSDs"`
	}

	stdout, stderr, err := c.runner.RunCephBinary(ctx, nil, "osd", "status", "--format=json")
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return nil, err
	}

	data := osds{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data.OSDs, nil
}

func (c *cluster) GetOSDIDs(ctx context.Context) ([]uint64, error) {
	stdout, stderr, err := c.runner.RunCephBinary(ctx, nil, "osd", "ls", "--format=json")
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return nil, err
	}

	data := []uint64{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *cluster) GetMons(ctx context.Context) ([]Mon, error) {
	type mons struct {
		// ...
		Mons []Mon `json:"mons"`
		// ...
	}

	stdout, stderr, err := c.runner.RunCephBinary(ctx, nil, "mon", "dump", "--format=json")
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return nil, err
	}

	data := mons{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data.Mons, nil
}

func (c *cluster) GetPools(ctx context.Context) ([]Pool, error) {
	stdout, stderr, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "ls", "detail", "--format=json")
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return nil, err
	}

	data := []Pool{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *cluster) ResizePool(ctx context.Context, name string, size uint64) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "set", name, "size", strconv.FormatUint(size, 10))
	return err
}

func (c *cluster) ChangePoolPGNum(ctx context.Context, name string, pgs uint64) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "set", name, "pg_num", strconv.FormatUint(pgs, 10))
	return err
}

func (c *cluster) ReweightByUtilization(ctx context.Context) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "reweight-by-utilization")
	return err
}

func (c *cluster) DestroyOSD(ctx context.Context, id uint64) error {
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "out", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "down", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "orch", "daemon", "rm", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "destroy", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "purge", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "rm", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "auth", "del", "osd."+strconv.FormatUint(id, 10))
	_, _, _ = c.runner.RunCephBinary(ctx, nil, "osd", "crush", "rm", "osd."+strconv.FormatUint(id, 10))

	return nil
}

func (c *cluster) StopOSDDaemon(ctx context.Context, id uint64) error {
	_, stderr, err := c.runner.RunCephBinary(ctx, nil, "orch", "daemon", "stop", "osd."+strconv.FormatUint(id, 10))
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return err
	}

	return nil
}

func (c *cluster) SetFlag(ctx context.Context, flag Flag) error {
	_, stderr, err := c.runner.RunCephBinary(ctx, nil, "osd", "set", string(flag))
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return err
	}

	return nil
}

func (c *cluster) UnsetFlag(ctx context.Context, flag Flag) error {
	_, stderr, err := c.runner.RunCephBinary(ctx, nil, "osd", "unset", string(flag))
	if err != nil {
		log.Debugf("command stderr: %s", string(stderr))
		return err
	}

	return nil
}

func (c *cluster) CreateDefaultPool(ctx context.Context, name string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "create", name)
	return err
}

func (c *cluster) CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error {
	_, _, err := c.runner.RunRadosBinary(ctx, data, "put", "--pool="+pool, objectName, "-")
	return err
}

func (c *cluster) SetNearFullRatio(ctx context.Context, value float64) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "set-nearfull-ratio", strconv.FormatFloat(value, 'f', -1, 64))
	return err
}

func (c *cluster) SetBackfillfullRatio(ctx context.Context, value float64) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "set-backfillfull-ratio", strconv.FormatFloat(value, 'f', -1, 64))
	return err
}

func (c *cluster) SetFullRatio(ctx context.Context, value float64) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "set-full-ratio", strconv.FormatFloat(value, 'f', -1, 64))
	return err
}

package shell

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/teran/ceph-chaos-monkey/ceph"
	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

type cluster struct {
	runner Runner
}

func New(runner Runner) drivers.Cluster {
	return &cluster{
		runner: runner,
	}
}

func (c *cluster) GetHealth(ctx context.Context) (ceph.Health, error) {
	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "health", "--format=json")
	if err != nil {
		return ceph.Health{}, err
	}

	data := ceph.Health{}
	return data, json.Unmarshal(stdout, &data)
}

func (c *cluster) GetOSDs(ctx context.Context) ([]ceph.OSD, error) {
	type osds struct {
		OSDs []ceph.OSD `json:"OSDs"`
	}

	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "status", "--format=json")
	if err != nil {
		return nil, err
	}

	data := osds{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data.OSDs, nil
}

func (c *cluster) GetOSDIDs(ctx context.Context) ([]uint64, error) {
	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "ls", "--format=json")
	if err != nil {
		return nil, err
	}

	data := []uint64{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *cluster) GetMons(ctx context.Context) ([]ceph.Mon, error) {
	type mons struct {
		// ...
		Mons []ceph.Mon `json:"mons"`
		// ...
	}

	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "mon", "dump", "--format=json")
	if err != nil {
		return nil, err
	}

	data := mons{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data.Mons, nil
}

func (c *cluster) GetPools(ctx context.Context) ([]ceph.Pool, error) {
	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "ls", "detail", "--format=json")
	if err != nil {
		return nil, err
	}

	data := []ceph.Pool{}
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
	_, _, err := c.runner.RunCephBinary(ctx, nil, "orch", "daemon", "stop", "osd."+strconv.FormatUint(id, 10))
	return err
}

func (c *cluster) SetFlag(ctx context.Context, flag ceph.Flag) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "set", string(flag))
	return err
}

func (c *cluster) UnsetFlag(ctx context.Context, flag ceph.Flag) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "unset", string(flag))
	return err
}

func (c *cluster) SetGroupFlag(ctx context.Context, flag ceph.Flag, group ...string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, append([]string{"osd", "set-group", string(flag)}, group...)...)
	return err
}

func (c *cluster) UnsetGroupFlag(ctx context.Context, flag ceph.Flag, group ...string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, append([]string{"osd", "unset-group", string(flag)}, group...)...)
	return err
}

func (c *cluster) CreateDefaultPool(ctx context.Context, name string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "osd", "pool", "create", name)
	return err
}

func (c *cluster) CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error {
	_, _, err := c.runner.RunRadosBinary(ctx, data, "put", "--pool="+pool, objectName, "-")
	return err
}

func (c *cluster) ReadRADOSObject(ctx context.Context, pool, objectName string) ([]byte, error) {
	stdout, _, err := c.runner.RunRadosBinary(ctx, nil, "get", "--pool="+pool, objectName, "-")
	if err != nil {
		return nil, err
	}

	return stdout, nil
}

func (c *cluster) ListRADOSObjects(ctx context.Context, pool string) ([]string, error) {
	type object struct {
		Namespace string `json:"namespace"`
		Name      string `json:"name"`
	}

	stdout, _, err := c.runner.RunRadosBinary(ctx, nil, "ls", "--pool="+pool, "--format=json")
	if err != nil {
		return nil, err
	}

	data := []object{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	out := []string{}
	for _, v := range data {
		out = append(out, v.Name)
	}

	return out, nil
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

func (c *cluster) RemoveMonitor(ctx context.Context, name string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "mon", "remove", name)
	return err
}

func (c *cluster) ListHosts(ctx context.Context) ([]ceph.Host, error) {
	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "orch", "host", "ls", "--format=json")
	if err != nil {
		return nil, err
	}

	data := []ceph.Host{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *cluster) DrainHost(ctx context.Context, hostname string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "orch", "host", "drain", hostname)
	return err
}

func (c *cluster) ListPGs(ctx context.Context) ([]ceph.PGStat, error) {
	type pgstat struct {
		PgStats []ceph.PGStat `json:"pg_stats"`
	}

	stdout, _, err := c.runner.RunCephBinary(ctx, nil, "pg", "ls", "--format=json")
	if err != nil {
		return nil, err
	}

	data := pgstat{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		return nil, err
	}

	return data.PgStats, nil
}

func (c *cluster) DeepScrubPG(ctx context.Context, target string) error {
	_, _, err := c.runner.RunCephBinary(ctx, nil, "pg", "deep-scrub", target)
	if err != nil {
		return err
	}
	return err
}

package drivers

import (
	"context"

	"github.com/teran/ceph-chaos-monkey/ceph"
)

type Cluster interface {
	GetHealth(ctx context.Context) (ceph.Health, error)

	RemoveMonitor(ctx context.Context, name string) error

	GetOSDs(ctx context.Context) ([]ceph.OSD, error)
	GetOSDIDs(ctx context.Context) ([]uint64, error)
	GetMons(ctx context.Context) ([]ceph.Mon, error)

	DestroyOSD(ctx context.Context, id uint64) error
	StopOSDDaemon(ctx context.Context, id uint64) error

	SetFlag(ctx context.Context, flag ceph.Flag) error
	UnsetFlag(ctx context.Context, flag ceph.Flag) error

	GetPools(ctx context.Context) ([]ceph.Pool, error)
	CreateDefaultPool(ctx context.Context, name string) error
	ResizePool(ctx context.Context, name string, size uint64) error
	ChangePoolPGNum(ctx context.Context, name string, pgs uint64) error
	ReweightByUtilization(ctx context.Context) error

	CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error
	ReadRADOSObject(ctx context.Context, pool, objectName string) ([]byte, error)
	ListRADOSObjects(ctx context.Context, pool string) ([]string, error)

	SetNearFullRatio(ctx context.Context, value float64) error
	SetBackfillfullRatio(ctx context.Context, value float64) error
	SetFullRatio(ctx context.Context, value float64) error

	ListHosts(ctx context.Context) ([]ceph.Host, error)
	DrainHost(ctx context.Context, hostname string) error
}

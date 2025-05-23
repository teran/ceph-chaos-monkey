package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/teran/ceph-chaos-monkey/ceph"
	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

var _ drivers.Cluster = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

func New() *Mock {
	return &Mock{}
}

func (m *Mock) GetHealth(context.Context) (ceph.Health, error) {
	args := m.Called()
	return args.Get(0).(ceph.Health), args.Error(1)
}

func (m *Mock) GetOSDs(context.Context) ([]ceph.OSD, error) {
	args := m.Called()
	return args.Get(0).([]ceph.OSD), args.Error(1)
}

func (m *Mock) GetOSDIDs(ctx context.Context) ([]uint64, error) {
	args := m.Called()
	return args.Get(0).([]uint64), args.Error(1)
}

func (m *Mock) GetMons(context.Context) ([]ceph.Mon, error) {
	args := m.Called()
	return args.Get(0).([]ceph.Mon), args.Error(1)
}

func (m *Mock) DestroyOSD(_ context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Mock) StopOSDDaemon(_ context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *Mock) SetFlag(_ context.Context, flag ceph.Flag) error {
	args := m.Called(flag)
	return args.Error(0)
}

func (m *Mock) UnsetFlag(_ context.Context, flag ceph.Flag) error {
	args := m.Called(flag)
	return args.Error(0)
}

func (m *Mock) SetGroupFlag(_ context.Context, flag ceph.Flag, group ...string) error {
	args := m.Called(flag, group)
	return args.Error(0)
}

func (m *Mock) UnsetGroupFlag(ctx context.Context, flag ceph.Flag, group ...string) error {
	args := m.Called(flag, group)
	return args.Error(0)
}

func (m *Mock) GetPools(context.Context) ([]ceph.Pool, error) {
	args := m.Called()
	return args.Get(0).([]ceph.Pool), args.Error(1)
}

func (m *Mock) ResizePool(_ context.Context, name string, size uint64) error {
	args := m.Called(name, size)
	return args.Error(0)
}

func (m *Mock) ChangePoolPGNum(_ context.Context, name string, pgs uint64) error {
	args := m.Called(name, pgs)
	return args.Error(0)
}

func (m *Mock) ReweightByUtilization(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *Mock) CreateDefaultPool(ctx context.Context, name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *Mock) CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error {
	args := m.Called(pool, objectName, data)
	return args.Error(0)
}

func (m *Mock) ReadRADOSObject(_ context.Context, pool, objectName string) ([]byte, error) {
	args := m.Called(pool, objectName)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *Mock) ListRADOSObjects(_ context.Context, pool string) ([]string, error) {
	args := m.Called(pool)
	return args.Get(0).([]string), args.Error(1)
}

func (m *Mock) SetNearFullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *Mock) SetBackfillfullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *Mock) SetFullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *Mock) RemoveMonitor(_ context.Context, name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *Mock) ListHosts(ctx context.Context) ([]ceph.Host, error) {
	args := m.Called()
	return args.Get(0).([]ceph.Host), args.Error(1)
}

func (m *Mock) DrainHost(ctx context.Context, hostname string) error {
	args := m.Called(hostname)
	return args.Error(0)
}

func (m *Mock) ListPGs(context.Context) ([]ceph.PGStat, error) {
	args := m.Called()
	return args.Get(0).([]ceph.PGStat), args.Error(1)
}

func (m *Mock) DeepScrubPG(_ context.Context, target string) error {
	args := m.Called(target)
	return args.Error(0)
}

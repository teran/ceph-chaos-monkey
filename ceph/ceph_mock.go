package ceph

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Cluster = (*ClusterMock)(nil)

type ClusterMock struct {
	mock.Mock
}

func NewClusterMock() *ClusterMock {
	return &ClusterMock{}
}

func (m *ClusterMock) GetOSDs(context.Context) ([]OSD, error) {
	args := m.Called()
	return args.Get(0).([]OSD), args.Error(1)
}

func (m *ClusterMock) GetMons(context.Context) ([]Mon, error) {
	args := m.Called()
	return args.Get(0).([]Mon), args.Error(1)
}

func (m *ClusterMock) DestroyOSD(_ context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ClusterMock) StopOSDDaemon(_ context.Context, id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ClusterMock) SetFlag(_ context.Context, flag Flag) error {
	args := m.Called(flag)
	return args.Error(0)
}

func (m *ClusterMock) UnsetFlag(_ context.Context, flag Flag) error {
	args := m.Called(flag)
	return args.Error(0)
}

func (m *ClusterMock) GetPools(context.Context) ([]Pool, error) {
	args := m.Called()
	return args.Get(0).([]Pool), args.Error(1)
}

func (m *ClusterMock) ResizePool(_ context.Context, name string, size uint64) error {
	args := m.Called(name, size)
	return args.Error(0)
}

func (m *ClusterMock) ChangePoolPGNum(_ context.Context, name string, pgs uint64) error {
	args := m.Called(name, pgs)
	return args.Error(0)
}

func (m *ClusterMock) ReweightByUtilization(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *ClusterMock) CreateDefaultPool(ctx context.Context, name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *ClusterMock) CreateRADOSObject(ctx context.Context, pool, objectName string, data []byte) error {
	args := m.Called(pool, objectName)
	return args.Error(0)
}

func (m *ClusterMock) SetNearFullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *ClusterMock) SetBackfillfullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

func (m *ClusterMock) SetFullRatio(ctx context.Context, value float64) error {
	args := m.Called(value)
	return args.Error(0)
}

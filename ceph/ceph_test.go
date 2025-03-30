package ceph

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *cephTestSuite) TestGetOSDs() {
	stdout, err := os.ReadFile("testdata/osd-status.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "status", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	osds, err := s.cluster.GetOSDs(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]OSD{
		{
			HostName:     "",
			ID:           0,
			KbAvailable:  0,
			KbUsed:       0,
			ReadByteRate: 0,
			ReadOpsRate:  0,
			State: []string{
				"autoout",
				"exists",
			},
			WriteByteRate: 0,
			WriteOpsRate:  0,
		},
		{
			HostName:     "",
			ID:           1,
			KbAvailable:  0,
			KbUsed:       0,
			ReadByteRate: 0,
			ReadOpsRate:  0,
			State: []string{
				"exists",
			},
			WriteByteRate: 0,
			WriteOpsRate:  0,
		},
		{
			HostName:     "ceph03",
			ID:           2,
			KbAvailable:  29950885888,
			KbUsed:       2257174528,
			ReadByteRate: 0,
			ReadOpsRate:  0,
			State: []string{
				"exists",
				"up",
			},
			WriteByteRate: 0,
			WriteOpsRate:  0,
		},
	}, osds)
}

func (s *cephTestSuite) TestGetMons() {
	stdout, err := os.ReadFile("testdata/mon-dump.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"mon", "dump", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	mons, err := s.cluster.GetMons(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]Mon{
		{
			Rank:       0,
			Name:       "ceph01",
			Addr:       "100.64.65.19:6789/0",
			PublicAddr: "100.64.65.19:6789/0",
			Priority:   0,
			Weight:     0,
		},
		{
			Rank:       1,
			Name:       "ceph02",
			Addr:       "100.64.65.20:6789/0",
			PublicAddr: "100.64.65.20:6789/0",
			Priority:   0,
			Weight:     0,
		},
		{
			Rank:       2,
			Name:       "ceph03",
			Addr:       "100.64.65.21:6789/0",
			PublicAddr: "100.64.65.21:6789/0",
			Priority:   0,
			Weight:     0,
		},
	}, mons)
}

func (s *cephTestSuite) TestSetFlag() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "set", "norecover"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.SetFlag(s.ctx, FlagNoRecover)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestUnsetFlag() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "unset", "norecover"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.UnsetFlag(s.ctx, FlagNoRecover)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestDestroyOSD() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "out", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "down", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"orch", "daemon", "rm", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "destroy", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "purge", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "rm", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"auth", "del", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "crush", "rm", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.DestroyOSD(s.ctx, 10)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestGetPools() {
	stdout, err := os.ReadFile("testdata/osd-pool-ls-detail.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "pool", "ls", "detail", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	mons, err := s.cluster.GetPools(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]Pool{
		{
			PoolID:             1,
			PoolName:           ".mgr",
			CreateTime:         "2025-03-30T14:47:07.151988+0000",
			Size:               3,
			MinSize:            2,
			CrushRule:          0,
			PgAutoscaleMode:    "on",
			PgNum:              1,
			LastChange:         "19",
			QuotaMaxBytes:      0,
			QuotaMaxObjects:    0,
			ErasureCodeProfile: "",
		},
	}, mons)
}

func (s *cephTestSuite) TestResizePool() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "pool", "set", "test-pool", "size", "10"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.ResizePool(s.ctx, "test-pool", 10)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestChangePoolPGNum() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "pool", "set", "test-pool", "pg_num", "10"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.ChangePoolPGNum(s.ctx, "test-pool", 10)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestStopOSDDaemon() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"orch", "daemon", "stop", "osd.10"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.StopOSDDaemon(s.ctx, 10)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestReweightByUtilization() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "reweight-by-utilization", "100"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.ReweightByUtilization(s.ctx)
	s.Require().NoError(err)
}

// ======================= definitions =======================
type cephTestSuite struct {
	suite.Suite

	ctx        context.Context
	cluster    Cluster
	runnerMock *runnerMock
}

func (s *cephTestSuite) SetupTest() {
	s.ctx = context.TODO()
	s.runnerMock = newRunnerMock()
	s.cluster = New(s.runnerMock)
}

func (s *cephTestSuite) TearDownTest() {
	s.runnerMock.AssertExpectations(s.T())
	s.cluster = nil
}

func TestCephTestSuite(t *testing.T) {
	suite.Run(t, &cephTestSuite{})
}

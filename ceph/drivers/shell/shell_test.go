package shell

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/teran/ceph-chaos-monkey/ceph"
	"github.com/teran/ceph-chaos-monkey/ceph/drivers"
)

func (s *cephTestSuite) TestGetOSDs() {
	stdout, err := os.ReadFile("testdata/osd-status.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "status", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	osds, err := s.cluster.GetOSDs(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]ceph.OSD{
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

func (s *cephTestSuite) TestGetOSDIDs() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "ls", "--format=json"}).Return([]byte(`[0,3,4,5,67]`), []byte{}, nil).Once()

	ids, err := s.cluster.GetOSDIDs(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]uint64{0, 3, 4, 5, 67}, ids)
}

func (s *cephTestSuite) TestGetMons() {
	stdout, err := os.ReadFile("testdata/mon-dump.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"mon", "dump", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	mons, err := s.cluster.GetMons(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]ceph.Mon{
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

	err := s.cluster.SetFlag(s.ctx, ceph.FlagNoRecover)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestUnsetFlag() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "unset", "norecover"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.UnsetFlag(s.ctx, ceph.FlagNoRecover)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetGroupFlag() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "set-group", "norecover", "group1", "group2"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.SetGroupFlag(s.ctx, ceph.FlagNoRecover, "group1", "group2")
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestUnsetGroupFlag() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "unset-group", "norecover", "group1", "group2"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.UnsetGroupFlag(s.ctx, ceph.FlagNoRecover, "group1", "group2")
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
	s.Require().Equal([]ceph.Pool{
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
			Options: ceph.PoolOptions{
				PgNumMax: 32,
				PgNumMin: 1,
			},
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
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "reweight-by-utilization"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.ReweightByUtilization(s.ctx)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestCreateDefaultPool() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "pool", "create", "test-pool"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.CreateDefaultPool(s.ctx, "test-pool")
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestCreatRADOSObject() {
	s.runnerMock.On("RunRadosBinary", []byte(`test data`), []string{"put", "--pool=test-pool", "test-object", "-"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.CreateRADOSObject(s.ctx, "test-pool", "test-object", []byte(`test data`))
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetNearFullRatio() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "set-nearfull-ratio", "0.15"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.SetNearFullRatio(s.ctx, 0.15)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetBackfillfullRatio() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "set-backfillfull-ratio", "0.2"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.SetBackfillfullRatio(s.ctx, 0.2)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetFullRatio() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"osd", "set-full-ratio", "0.3"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.SetFullRatio(s.ctx, 0.3)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestRemoveMonitor() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"mon", "remove", "test"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.RemoveMonitor(s.ctx, "test")
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestReadRADOSObject() {
	s.runnerMock.On("RunRadosBinary", []byte(nil), []string{"get", "--pool=test-pool", "object-name", "-"}).Return([]byte("test data"), []byte{}, nil).Once()

	data, err := s.cluster.ReadRADOSObject(s.ctx, "test-pool", "object-name")
	s.Require().NoError(err)
	s.Require().Equal("test data", string(data))
}

func (s *cephTestSuite) TestListRADOSObjects() {
	s.runnerMock.On("RunRadosBinary", []byte(nil), []string{"ls", "--pool=test-pool", "--format=json"}).Return([]byte(`[{"name":"obj1"},{"name":"obj2"}]`), []byte{}, nil).Once()

	pools, err := s.cluster.ListRADOSObjects(s.ctx, "test-pool")
	s.Require().NoError(err)
	s.Require().Equal([]string{"obj1", "obj2"}, pools)
}

func (s *cephTestSuite) TestListHosts() {
	stdout, err := os.ReadFile("testdata/orch-host-ls.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"orch", "host", "ls", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	hosts, err := s.cluster.ListHosts(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]ceph.Host{
		{
			Hostname: "ceph01",
			Addr:     "100.64.65.19",
			Labels:   []string{"_admin", "mon", "mgr"},
		},
		{
			Hostname: "ceph02",
			Addr:     "100.64.65.20",
			Labels:   []string{"_admin", "mon", "mgr"},
		},
		{
			Hostname: "ceph03",
			Addr:     "100.64.65.21",
			Labels:   []string{"_admin", "mon", "mgr"},
		},
		{
			Hostname: "ceph04",
			Addr:     "100.64.65.23",
			Labels:   []string{"_admin", "_no_schedule", "_no_conf_keyring"},
		},
		{
			Hostname: "ceph05",
			Addr:     "100.64.65.24",
			Labels:   []string{"_admin", "_no_schedule", "_no_conf_keyring"},
		},
	}, hosts)
}

func (s *cephTestSuite) TestDrainHost() {
	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"orch", "host", "drain", "test-host"}).Return([]byte{}, []byte{}, nil).Once()

	err := s.cluster.DrainHost(s.ctx, "test-host")
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestListPGs() {
	stdout, err := os.ReadFile("testdata/pg-ls.json")
	s.Require().NoError(err)

	s.runnerMock.On("RunCephBinary", []byte(nil), []string{"pg", "ls", "--format=json"}).Return(stdout, []byte{}, nil).Once()
	pgs, err := s.cluster.ListPGs(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal([]ceph.PGStat{
		{
			PGID:  "1.0",
			State: "active+clean",
			Up:    []uint64{3, 2, 1},
		},
	}, pgs)
}

// ======================= definitions =======================
type cephTestSuite struct {
	suite.Suite

	ctx        context.Context
	cluster    drivers.Cluster
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

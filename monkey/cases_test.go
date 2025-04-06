package monkey

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/teran/go-collection/random"

	"github.com/teran/ceph-chaos-monkey/ceph"
	clusterMock "github.com/teran/ceph-chaos-monkey/ceph/drivers/mock"
)

func (s *cephTestSuite) TestSetRandomFlag() {
	s.rnd.On("Intn", len(cephFlags)).Return(3).Once()
	s.cluster.On("SetFlag", ceph.FlagNoOut).Return(nil).Once()

	err := setRandomFlag(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestUnsetRandomFlag() {
	s.rnd.On("Intn", len(cephFlags)).Return(4).Once()
	s.cluster.On("UnsetFlag", ceph.FlagNoRebalance).Return(nil).Once()

	err := unsetRandomFlag(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestDestroyRandomOSD() {
	s.cluster.On("GetOSDIDs").Return([]uint64{3, 5, 7, 9}, nil).Once()
	s.rnd.On("Intn", 4).Return(3).Once()
	s.cluster.On("DestroyOSD", uint64(9)).Return(nil).Once()

	err := destroyRandomOSD(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestRandomlyResizeRandomPool() {
	s.cluster.On("GetPools").Return([]ceph.Pool{
		{PoolID: 3, PoolName: "pool1"},
		{PoolID: 4, PoolName: "pool2"},
	}, nil).Once()
	s.rnd.On("Intn", 2).Return(1).Once()
	s.rnd.On("Intn", 10).Return(3).Once()
	s.cluster.On("ResizePool", "pool2", uint64(3)).Return(nil).Once()

	err := randomlyResizeRandomPool(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestRandomlyChangePGNumForRandomPool() {
	s.cluster.On("GetPools").Return([]ceph.Pool{
		{PoolID: 1, PoolName: "pool1", Options: ceph.PoolOptions{PgNumMax: 3}},
		{PoolID: 2, PoolName: "pool2", Options: ceph.PoolOptions{PgNumMax: 5}},
	}, nil).Once()
	s.rnd.On("Intn", 2).Return(1).Once()
	s.rnd.On("Intn", 5).Return(4).Once()
	s.cluster.On("ChangePoolPGNum", "pool2", uint64(4+1)).Return(nil).Once()

	err := randomlyChangePGNumForRandomPool(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestReweightByUtilization() {
	s.cluster.On("ReweightByUtilization").Return(nil).Once()

	err := reweightByUtilization(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetRandomNearFullRatio() {
	s.rnd.On("Float64").Return(0.75).Once()
	s.cluster.On("SetNearFullRatio", 0.75).Return(nil).Once()

	err := setRandomNearFullRatio(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetRandomBackfillfullRatio() {
	s.rnd.On("Float64").Return(0.85).Once()
	s.cluster.On("SetNearFullRatio", 0.85).Return(nil).Once()

	err := setRandomBackfillfullRatio(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestSetRandomFullRatio() {
	s.rnd.On("Float64").Return(0.95).Once()
	s.cluster.On("SetNearFullRatio", 0.95).Return(nil).Once()

	err := setRandomFullRatio(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestRemoveMonitor() {
	s.cluster.On("GetMons").Return([]ceph.Mon{
		{Name: "test1"},
		{Name: "test2"},
		{Name: "test3"},
	}, nil).Once()
	s.rnd.On("Intn", 3).Return(1).Once()
	s.cluster.On("RemoveMonitor", "test2").Return(nil).Once()

	err := removeRandomMonitor(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

// ======================= definitions =======================
type cephTestSuite struct {
	suite.Suite

	ctx     context.Context
	cluster *clusterMock.Mock
	rnd     *random.Mock
}

func (s *cephTestSuite) SetupTest() {
	s.ctx = context.TODO()

	s.rnd = random.NewMock()
	s.cluster = clusterMock.New()
}

func (s *cephTestSuite) TearDownTest() {
	s.rnd.AssertExpectations(s.T())
	s.cluster.AssertExpectations(s.T())
}

func TestCephTestSuite(t *testing.T) {
	suite.Run(t, &cephTestSuite{})
}

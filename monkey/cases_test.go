package monkey

import (
	"context"
	"math"
	"strconv"
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
		{PoolID: 1, PoolName: "pool1"},
		{PoolID: 2, PoolName: "pool2"},
	}, nil).Once()
	s.rnd.On("Intn", 2).Return(1).Once()
	s.rnd.On("Intn", 256).Return(4).Once()
	s.cluster.On("ChangePoolPGNum", "pool2", uint64(4)).Return(nil).Once()

	err := randomlyChangePGNumForRandomPool(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestReweightByUtilization() {
	s.cluster.On("ReweightByUtilization").Return(nil).Once()

	err := reweightByUtilization(s.ctx, s.cluster, s.rnd)
	s.Require().NoError(err)
}

func (s *cephTestSuite) TestCreatePoolAndPutAmountOfObjects() {
	poolNameEntropy := int64(1234567890)
	bufSize := 1024
	amountOfObjects := 2

	poolName := "test-pool-" + strconv.FormatInt(1234567890, 10)
	s.rnd.On("Int63n", int64(math.MaxInt64)).Return(poolNameEntropy).Once()
	s.rnd.On("Intn", 50).Return(amountOfObjects).Once()
	s.rnd.On("Intn", 150*1024*1024).Return(bufSize).Times(amountOfObjects)

	buf := make([]byte, bufSize)
	copy(buf, []byte("test-data"))

	s.cluster.On("CreateDefaultPool", poolName).Return(nil).Once()
	s.rnd.On("Read").Return(buf, bufSize, nil).Times(amountOfObjects)
	s.cluster.
		On(
			"CreateRADOSObject",
			poolName,
			"ec7c4235ee1166be806c2b7d69a05939726355bbbf8520743f2192b2c930cdaf",
			buf,
		).
		Return(nil).Times(amountOfObjects)

	err := createPoolAndPutAmountOfObjects(s.ctx, s.cluster, s.rnd)
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

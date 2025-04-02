package shell

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Runner = (*runnerMock)(nil)

type runnerMock struct {
	mock.Mock
}

func newRunnerMock() *runnerMock {
	return &runnerMock{}
}

func (m *runnerMock) RunCephBinary(_ context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error) {
	p := m.Called(stdin, args)
	return p.Get(0).([]byte), p.Get(1).([]byte), p.Error(2)
}

func (m *runnerMock) RunRadosBinary(_ context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error) {
	p := m.Called(stdin, args)
	return p.Get(0).([]byte), p.Get(1).([]byte), p.Error(2)
}

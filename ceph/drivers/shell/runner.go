package shell

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var _ Runner = (*runner)(nil)

type Runner interface {
	RunCephBinary(ctx context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error)
	RunRadosBinary(ctx context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error)
}

type runner struct {
	cephBinaryPath  string
	radosBinaryPath string
}

func NewRunner(cephBinaryPath, radosBinaryPath string) Runner {
	return &runner{
		cephBinaryPath:  cephBinaryPath,
		radosBinaryPath: radosBinaryPath,
	}
}

func (r *runner) RunRadosBinary(ctx context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error) {
	return run(ctx, stdin, r.radosBinaryPath, args...)
}

func (r *runner) RunCephBinary(ctx context.Context, stdin []byte, args ...string) (stdoutContents []byte, stderrContents []byte, err error) {
	return run(ctx, stdin, r.cephBinaryPath, args...)
}

func run(ctx context.Context, stdin []byte, cmd string, args ...string) (stdoutContents []byte, stderrContents []byte, err error) {
	log.Tracef("preparing command: %s %#v", cmd, args)
	c := exec.CommandContext(ctx, cmd, args...)

	if stdin != nil {
		c.Stdin = bytes.NewReader(stdin)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	c.Stdout = stdout
	c.Stderr = stderr

	if err := c.Start(); err != nil {
		return nil, nil, err
	}

	if err := c.Wait(); err != nil {
		return nil, nil, err
	}

	outStdout := stdout.Bytes()
	outStderr := stderr.Bytes()

	isBinaryGetOp := false
	for _, arg := range args {
		if arg == "get" {
			isBinaryGetOp = true
		}
	}
	if !strings.HasSuffix(cmd, "rados") || !isBinaryGetOp {
		log.Debugf("data received [stdout]: %s\n", string(outStdout))
	}

	log.Debugf("data received [stderr]: %s\n", string(outStderr))
	log.Debugf("exit code: %d", c.ProcessState.ExitCode())

	return outStdout, outStderr, nil
}

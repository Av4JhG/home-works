package main

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("should work w/o arguments", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		if err != nil {
			require.Fail(t, "unexpected error while reading env")
		}

		exitCode := RunCmd([]string{"ls"}, env)

		require.Equal(t, 0, exitCode)
	})

	t.Run("should work w/o env", func(t *testing.T) {
		exitCode := RunCmd([]string{"ls", "-l"}, EmptyEnv())

		require.Equal(t, 0, exitCode)
	})
}

func TestRunCmdError(t *testing.T) {
	t.Run("run without command provided", func(t *testing.T) {
		exitCode := RunCmd([]string{}, EmptyEnv())

		require.Equal(t, int(syscall.EINVAL), exitCode)
	})
}

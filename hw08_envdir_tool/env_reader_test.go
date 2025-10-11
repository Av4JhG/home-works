package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("error reading dir", func(t *testing.T) {
		dir := "/root"

		env, err := ReadDir(dir)

		require.Nil(t, env)
		require.ErrorContains(t, err, ErrEnvRead.Error())
	})
}

func TestSanitizeVal(t *testing.T) {
	tests := []struct {
		name, output string
		input        []byte
	}{
		{
			name:   "trailing tab should be removed",
			input:  []byte("\tabracadabra\t\t"),
			output: "\tabracadabra",
		},
		{
			name:   "trailing space should be removed",
			input:  []byte("\tabracadabra    "),
			output: "\tabracadabra",
		},
		{
			name:   "NUL symbol should be replaced with new line symbol",
			input:  []byte("abrac\x00adabra"),
			output: "abrac\nadabra",
		},
		{
			name:   "complex",
			input:  []byte("  abrac\x00adabra\t \t"),
			output: "  abrac\nadabra",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.output, sanitizeVal(test.input))
		})
	}
}

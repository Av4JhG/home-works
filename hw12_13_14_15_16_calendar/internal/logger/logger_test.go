package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		level         string
		expectedLog   string
		expectedLevel logrus.Level
	}{
		{"debug", "Это Debug сообщение", logrus.DebugLevel},
		{"info", "Это Info сообщение", logrus.InfoLevel},
		{"warn", "Это Warn сообщение", logrus.WarnLevel},
		{"error", "Это Error сообщение", logrus.ErrorLevel},
		{"invalid", "Это Info сообщение", logrus.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger := New(tt.level)
			var buf bytes.Buffer
			logger.SetOutput(&buf)
			logger.Log(tt.expectedLevel, tt.expectedLog)
			require.Contains(t, buf.String(), tt.expectedLog)
			require.Equal(t, tt.expectedLevel, logger.Logger.GetLevel())
		})
	}
}

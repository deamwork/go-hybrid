package errorw

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/XSAM/go-hybrid/trace"
)

func TestError_MarshalLogObject(t *testing.T) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := zap.New(ob)

	logger.Info("test", zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e})

	log := logs.All()[0]
	contextMap := log.ContextMap()
	err := contextMap["error"].(map[string]interface{})
	assert.Equal(t, "test", log.Message)
	assert.Equal(t, "wrap: testing error", err["msg"])
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
	}, err["fields"])
	assert.Contains(t, err["stack"], "go-hybrid/errorw/stack_test.go:12")
}

func TestError_MarshalLogObject2(t *testing.T) {
	ob, logs := observer.New(zapcore.InfoLevel)
	logger := zap.New(ob)

	// Set trace id
	e := New(trace.SetTraceIDToContext(context.Background(), "trace_id"), errors.New("testing error")).WithWrap("wrap")
	// Set cause error is nil
	e.Err = nil

	logger.Info("test", zap.Field{Key: "error", Type: zapcore.ObjectMarshalerType, Interface: e})

	log := logs.All()[0]
	contextMap := log.ContextMap()
	err := contextMap["error"].(map[string]interface{})
	assert.Equal(t, "test", log.Message)
	assert.Equal(t, "wrap: nil", err["msg"])
}

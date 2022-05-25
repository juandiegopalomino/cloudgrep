package util

import (
	"context"
	"testing"

	"go.uber.org/zap"

	"github.com/matishsiao/goInfo"
	"github.com/run-x/cloudgrep/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAmplitudeEvent(t *testing.T) {
	systemInfo, err := goInfo.GetInfo()
	if err != nil {
		t.Error(err)
	}
	type args struct {
		eventType       string
		eventProperties map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "TestGenerateAmplitudeEvent",
			args: args{
				eventType:       BaseEvent,
				eventProperties: nil,
			},
			want: map[string]interface{}{
				"user_id":          userId,
				"device_id":        systemInfo.Hostname,
				"event_type":       BaseEvent,
				"event_properties": map[string]interface{}{"application": application},
				"app_version":      version.Version,
				"platform":         systemInfo.Platform,
				"insert_id":        "test",
				"session_id":       sessionId,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := generateAmplitudeEvent(tt.args.eventType, tt.args.eventProperties)
			assert.Equal(t, got["user_id"], userId)
			assert.Equal(t, got["device_id"], systemInfo.Hostname)
			assert.Equal(t, got["event_type"], tt.args.eventType)
			assert.Equal(t, got["app_version"], version.Version)
			assert.Equal(t, got["platform"], systemInfo.Platform)
			gotEventProperties := got["event_properties"].(map[string]interface{})
			assert.Equal(t, gotEventProperties["application"], application)
		})
	}
}

func TestSendAmplitudeEvent(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	ctx := context.Background()

	t.Run("TestSendAmplitudeEventDevVersion", func(t *testing.T) {
		returnValue, err := sendAmplitudeEvent(ctx, logger, BaseEvent, nil)
		assert.Equal(t, returnValue, 1)
		assert.ErrorContains(t, err, "dev application, not sending events to amplitude")
	})

	version.Version = "test"
	t.Run("TestSendAmplitudeEventInvalidEvent", func(t *testing.T) {
		returnValue, err := sendAmplitudeEvent(ctx, logger, "INVALID_EVENT", nil)
		assert.Equal(t, returnValue, 1)
		assert.ErrorContains(t, err, "invalid event type: INVALID_EVENT, not sending events to amplitude\n")
	})

}

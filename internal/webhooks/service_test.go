package webhooks

import (
	"testing"

	coreTesting "github.com/brigadecore/brigade/sdk/v2/testing/core"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	s, ok := NewService(
		// Totally unusable client that is enough to fulfill the dependencies for
		// this test...
		&coreTesting.MockEventsClient{
			LogsClient: &coreTesting.MockLogsClient{},
		},
		ServiceConfig{},
	).(*service)
	require.True(t, ok)
	require.NotNil(t, s.eventsClient)
	require.NotNil(t, s.config)
}

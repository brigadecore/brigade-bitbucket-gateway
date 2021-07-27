package webhooks

import (
	"testing"

	coreTesting "github.com/brigadecore/brigade/sdk/v2/testing/core"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	s := NewService(
		// Totally unusable client that is enough to fulfill the dependencies for
		// this test...
		&coreTesting.MockEventsClient{
			LogsClient: &coreTesting.MockLogsClient{},
		},
		ServiceConfig{},
	).(*service)
	require.NotNil(t, s.eventsClient)
	require.NotNil(t, s.config)
}

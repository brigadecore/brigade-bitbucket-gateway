package webhooks

import (
	"testing"

	sdkTesting "github.com/brigadecore/brigade/sdk/v3/testing"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	s, ok := NewService(
		// Totally unusable client that is enough to fulfill the dependencies for
		// this test...
		&sdkTesting.MockEventsClient{
			LogsClient: &sdkTesting.MockLogsClient{},
		},
	).(*service)
	require.True(t, ok)
	require.NotNil(t, s.eventsClient)
}

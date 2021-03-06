package auth

import (
	"os/exec"
	"testing"

	"github.com/pachyderm/pachyderm/src/client/pkg/require"
)

func TestValidateActivationCode(t *testing.T) {
	code := "eyJ0b2tlbiI6IntcImV4cGlyeVwiOlwiMjAyNy0wNy0xMlQwMzowOTowNi4xODBaXCIsXCJzY29wZXNcIjp7XCJiYXNpY1wiOnRydWV9LFwibmFtZVwiOlwicGFjaHlkZXJtRW5naW5lZXJpbmdcIn0iLCJzaWduYXR1cmUiOiJWdjZQbEkrL3RJamlWYUNHMGw0TUd6WldDS2YrUFMyWS9WUzFkZ3plcVdjS3RETlJvdkRnSnd3TXFXbWdCOUs5a2lPemVQRlh4eTh3U2dMbTJ4dnBlTHN2bGlsTlc5MEhKbGxxcjhKWEVTbVV4R2tKQldMTHZHak5mYUlHZ0IvZTFEMzQzMi95eUVnSW1LZDlpZ3J3RXZsRCtGdW0wa1hqS3Rrb2pPRmhkMDR6RHFEMSt5ZWpsTmRtUzBTaDJKWHRTMnFqWk0zTE5lWlpTRldLcEVJTmlXa2dhOTdTNUw2ZVlCdXFZcFJLMTkwd1pXNTVCOVFJSHJNNWtDWGQrWUN5aTh0QU9kcFY2a3FMSDNoVGgxVDIwVjYveFNZNUVheHZObm8yRmFYbDUyQzRFSWIvZ05RWW8xVExDd1hJN0FYL2lpL0VTckVBQmYzdDlYZmlwWGxleE9OMmhJaWY5dDROZFBaQ1pmYlErbW8vSlQ3Um5VTGpTb2J3alNWVk1qMUozLzZKbmhQRFpFSWNDdlVvUnMyL2M2WUZxOVo1TFRJNkUxV2Q0bE1RczRJYXVsTHVQOEFVa3R3ejBiQmY2dUhPd3VvTlk4UjJ3ZTA1MmUxWVVGbmNyUE4wd2ZJVHo5Vm51M1dNcktpaDhhRzNmMzRLb2x0R3hpWXJHL2JZQjgweUFaTytCbzFmTTJwaDB0emRXejFLR0lNQUlEbjBFWHU2V0duSUFFUWN1NHVFc1pSVXRzNFhuYk5PTC9vYU1NK3RLV3UzdnFMdEhMWWlPaWZHNHpEcUxwYnNNN2NhZGNXWjJ3QzNoZVh6Y1loaUwzMHJlOGJ4MFc3Vm1FOSt4elJHZisyNEdvRjFaS1BvaDNhY3hCS0dsZzRxN2JQd0c3QWJESmxkak1HbkVEdz0ifQ=="

	require.NoError(t, validateActivationCode(code))
}

func TestGetSetBasic(t *testing.T) {
	// TODO(msteffen): replace this with a test that uses the client and checks
	// output. For now, this just tests that the auth service doesn't crash
	cmd := exec.Command("bash", "-c", `
	set -ex
	pachctl delete-repo test || true
	pachctl create-repo test
	pachctl auth get test
	pachctl auth set testuser reader test
	pachctl auth get testuser test
	`)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)
	t.Log(string(output))
}

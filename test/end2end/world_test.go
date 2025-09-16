package end2end

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var baseURL = "http://localhost:8080"
var client = &http.Client{}

func DoRequest[T any](t *testing.T, method, path string, body interface{}, headers map[string]string) (T, *http.Response) {
	t.Helper()

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	defer func() {
		// We close after decoding
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	var out T
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil && err != io.EOF {
		require.NoError(t, err)
	}

	return out, resp
}

func TestWorldsLifecycle(t *testing.T) {
	userID := uuid.New().String()
	_, resp := DoRequest[interface{}](t, http.MethodPost, "/user/"+userID, nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	worlds, resp := DoRequest[[]map[string]interface{}](t, http.MethodGet, "/worlds", nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, worlds, 0)

	newWorld := map[string]string{
		"name":        "Test World",
		"description": "from e2e test",
	}
	authHeaders := map[string]string{"Authorization": "Bearer " + userID}
	created, resp := DoRequest[map[string]interface{}](t, http.MethodPost, "/worlds", newWorld, authHeaders)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.Equal(t, "Test World", created["name"])
	require.NotEmpty(t, created["id"])

	worlds, resp = DoRequest[[]map[string]interface{}](t, http.MethodGet, "/worlds", nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, worlds, 1)
	require.Equal(t, "Test World", worlds[0]["name"])
}

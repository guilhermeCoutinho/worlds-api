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

	var out T
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return out, resp
	}

	defer func() {
		// We close after decoding
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil && err != io.EOF {
		require.NoError(t, err)
	}

	return out, resp
}

func TestWorldsLifecycle(t *testing.T) {
	userID := uuid.New().String()
	_, resp := DoRequest[interface{}](t, http.MethodPost, "/user/"+userID, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	worlds, resp := DoRequest[[]map[string]interface{}](t, http.MethodGet, "/worlds", nil, nil)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	oldLen := len(worlds)

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
	require.Len(t, worlds, oldLen+1)
	require.Equal(t, "Test World", worlds[0]["name"])
}

func TestCannotEditWorldsFromOtherUsers(t *testing.T) {
	userA := uuid.New().String()
	userB := uuid.New().String()

	_, resp := DoRequest[interface{}](t, http.MethodPost, "/user/"+userA, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	_, resp = DoRequest[interface{}](t, http.MethodPost, "/user/"+userB, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	newWorld := map[string]string{
		"name":        "Test World",
		"description": "from e2e test",
	}
	authHeadersUserA := map[string]string{"Authorization": "Bearer " + userA}
	worldA, resp := DoRequest[map[string]interface{}](t, http.MethodPost, "/worlds", newWorld, authHeadersUserA)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	worldAID := worldA["id"].(string)

	authHeadersUserB := map[string]string{"Authorization": "Bearer " + userB}
	_, resp = DoRequest[interface{}](t, http.MethodPut, "/worlds/"+worldAID, newWorld, authHeadersUserB)
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestUpdateWorld(t *testing.T) {
	user := uuid.New().String()
	_, resp := DoRequest[interface{}](t, http.MethodPost, "/user/"+user, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	originalWorld := map[string]string{
		"name":        "Old Name",
		"description": "old description",
	}
	authHeaders := map[string]string{"Authorization": "Bearer " + user}
	oldWorld, resp := DoRequest[map[string]interface{}](t, http.MethodPost, "/worlds", originalWorld, authHeaders)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.Equal(t, originalWorld["name"], oldWorld["name"])
	require.Equal(t, originalWorld["description"], oldWorld["description"])

	updatedWorld := map[string]string{
		"name":        "New Name",
		"description": "new description",
	}
	updatedWorldResponse, resp := DoRequest[map[string]interface{}](t, http.MethodPut, "/worlds/"+oldWorld["id"].(string), updatedWorld, authHeaders)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, updatedWorldResponse["name"], updatedWorld["name"])
	require.Equal(t, updatedWorldResponse["description"], updatedWorld["description"])

	worlds, resp := DoRequest[[]map[string]interface{}](t, http.MethodGet, "/worlds?ownerId="+user, nil, authHeaders)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, worlds, 1)
	require.Equal(t, updatedWorldResponse["name"], worlds[0]["name"])
	require.Equal(t, updatedWorldResponse["description"], worlds[0]["description"])
}

func TestGetWorldsByOwnerID(t *testing.T) {
	userA := uuid.New().String()
	userB := uuid.New().String()

	_, resp := DoRequest[interface{}](t, http.MethodPost, "/user/"+userA, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	_, resp = DoRequest[interface{}](t, http.MethodPost, "/user/"+userB, nil, nil)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	newWorld := map[string]string{
		"name":        "Test World",
		"description": "from e2e test",
	}
	authHeadersUserA := map[string]string{"Authorization": "Bearer " + userA}
	newWorldA, resp := DoRequest[map[string]interface{}](t, http.MethodPost, "/worlds", newWorld, authHeadersUserA)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	worlds, resp := DoRequest[[]map[string]interface{}](t, http.MethodGet, "/worlds?ownerId="+userA, nil, authHeadersUserA)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, worlds, 1)
	require.Equal(t, newWorldA["name"], worlds[0]["name"])
}

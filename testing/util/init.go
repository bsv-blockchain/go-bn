package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bsv-blockchain/go-bn/models"
)

type closeFunc func()

// TestServer creates a test server for testing.
func TestServer(t *testing.T, expReq *models.Request, testFile string) (*httptest.Server, closeFunc) { //nolint: revive // test code
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req models.Request
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		assert.Equal(t, *expReq, req)

		//nolint:gosec // test code
		response, err := os.ReadFile(path.Join("./testing/data", testFile+".json"))
		assert.NoError(t, err)

		mm := map[string]interface{}{}
		assert.NoError(t, json.Unmarshal(response, &mm))
		bb, err := json.Marshal(mm)
		assert.NoError(t, err)
		_, _ = w.Write(bb)
	}))

	return svr, svr.Close
}

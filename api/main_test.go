package api

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func getReaderFor[T any](t *testing.T, request T) io.Reader {
	buf, err := json.Marshal(request)
	require.NoError(t, err)
	return bytes.NewReader(buf)
}
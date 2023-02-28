package api

import (
	"os"
	"testing"
	"time"

	db "github.com/MasonPhan2110/SimpleBank/db/sqlc"
	"github.com/MasonPhan2110/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(store, config)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

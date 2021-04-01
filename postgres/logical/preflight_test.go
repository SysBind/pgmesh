package logical

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/sysbind/pgmesh/postgres"
)

func TestPreflight_WalLevel(t *testing.T) {
	is := is.New(t)

	src := postgres.ConnConfig{
		Host:     "localhost",
		Database: "postgres",
		User:     "postgres",
		Pass:     "q1w2e3r4"}

	ctx := context.Background()

	err := PreFlight(ctx, src, postgres.ConnConfig{})
	is.NoErr(err)
}

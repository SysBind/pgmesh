package util

import (
	"os/exec"

	"github.com/sysbind/pgmesh/postgres"
)

func Dump(conn postgres.ConnConfig) (cmd *exec.Cmd) {
	cmd = exec.Command("pg_dump", "-U", conn.User, conn.Database)
	return
}

func DumpGlobals(conn postgres.ConnConfig) (cmd *exec.Cmd) {
	cmd = exec.Command("pg_dumpall", "--globals-only", "-U", conn.User)
	return
}

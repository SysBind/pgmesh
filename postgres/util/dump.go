package util

import (
	"io"

	"github.com/sysbind/pgmesh/postgres"
)

func dump(conn postgres.ConnConfig) (stdout io.ReadCloser, err error) {
	cmd := exec.Command("pg_dump", "moodle")
	stdout, err = exec.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	return
}

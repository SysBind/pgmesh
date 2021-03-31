package util

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/sysbind/pgmesh/postgres"
)

func Dump(conn postgres.ConnConfig) (<-chan string, <-chan string, error) {
	outchan := make(chan string)
	errchan := make(chan string)

	cmd := exec.Command("pg_dump", "-U", conn.User, conn.Database)

	err := run(cmd, outchan, errchan)
	if err != nil {
		return nil, nil, err
	}

	return outchan, errchan, nil
}

func DumpGlobals(conn postgres.ConnConfig) (<-chan string, <-chan string, error) {
	outchan := make(chan string)
	errchan := make(chan string)

	cmd := exec.Command("pg_dumpall", "--globals-only", "-U", conn.User)

	err := run(cmd, outchan, errchan)
	if err != nil {
		return nil, nil, err
	}

	return outchan, errchan, nil
}

func run(cmd *exec.Cmd, outchan chan<- string, errchan chan<- string) (err error) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outchan <- scanner.Text()
		}
		close(outchan)
	}()

	go func() {
		stderr_str, _ := io.ReadAll(stderr)
		errchan <- string(stderr_str)
		close(errchan)
		cmd.Wait()
	}()

	return
}

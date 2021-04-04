package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"

	"github.com/sysbind/pgmesh/postgres"
)

type DumpSection int

const (
	PreData DumpSection = iota
	PostData
)

func (sec DumpSection) String() string {
	switch sec {
	case PreData:
		return "pre-data"
	case PostData:
		return "post-data"
	}
	return fmt.Sprintf("no such DumpSection %d", sec)
}

// DumpSchema will execute pg_dump --schema-only --section=[pre-data,post-data]
// returning chnnels for stdout & stderr
func DumpSchema(conn postgres.ConnConfig, section DumpSection) (<-chan string, <-chan string, error) {
	pgpass(conn)
	outchan := make(chan string)
	errchan := make(chan string)

	cmd := exec.Command("pg_dump",
		fmt.Sprintf("--section=%s", section),
		"-U", conn.User,
		"--host", conn.Host,
		conn.Database)

	err := run(cmd, outchan, errchan)
	if err != nil {
		return nil, nil, err
	}

	return outchan, errchan, nil
}

// DumpGlobals will execute pg_dumpall --globals-only
// returning chnnels for stdout & stderr
func DumpGlobals(conn postgres.ConnConfig) (<-chan string, <-chan string, error) {
	pgpass(conn)
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
		defer close(outchan)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			outchan <- scanner.Text()
		}
	}()

	go func() {
		defer close(errchan)
		stderr_str, _ := io.ReadAll(stderr)
		errchan <- string(stderr_str)
		cmd.Wait()
	}()

	return
}

// write $HOME/.pgpass for passwordless operation
func pgpass(conn postgres.ConnConfig) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	str := fmt.Sprintf("%s:%s:%s:%s:%s",
		conn.Host,
		"5432",
		conn.Database,
		conn.User,
		conn.Pass)

	err = ioutil.WriteFile(usr.HomeDir+"/.pgpass", []byte(str), 0600)
	if err != nil {
		log.Fatal(err)
	}
}

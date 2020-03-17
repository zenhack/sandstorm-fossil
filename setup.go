package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func initMySchema(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS sandstorm_user (
		fossil_uid INTEGER
			REFERENCES user(uid)
			ON DELETE CASCADE,
		sandstorm_uid VARCHAR PRIMARY KEY
	)`)
	chkfatal(err)
}

func setupSystem() {
	chkfatal(ioutil.WriteFile(
		"/var/passwdfile",
		[]byte(fmt.Sprintf(
			"%s:x:%d:%d::%s:/bin/bash",
			os.Getenv("USER"),
			unix.Getuid(),
			unix.Getgid(),
			os.Getenv("HOME"),
		)),
		0644,
	))
}

func fossilCmd(args ...string) {
	args = append([]string{"--user", "grain"}, args...)
	log.Print("Running fossil with args: ", args)
	cmd := exec.Command("fossil", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	chkfatal(cmd.Run())
}

func maybeCreateFossil() {
	_, err := os.Stat(repoPath)
	if err == nil {
		return
	}
	log.Print("Couldn't stat repo: ", err)
	fossilCmd("new", repoPath)
}

func migrateFossil() {
	fossilCmd("rebuild", repoPath)
}

func runFossilServer() {
	fossilCmd("server", "--port", fossilPort, repoPath)
	panic("Fossil server exited!")
}

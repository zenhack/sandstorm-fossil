package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	repoPath   = os.Getenv("REPO_PATH")
	fossilPort = os.Getenv("FOSSIL_PORT")
	fossilAddr = net.JoinHostPort("127.0.0.1", fossilPort)
	proxyPort  = os.Getenv("PROXY_PORT")
)

func main() {
	setupSystem()
	maybeCreateFossil()
	migrateFossil()
	go runFossilServer()

	db, err := sql.Open("sqlite3", repoPath)
	chkfatal(err)
	defer db.Close()
	initMySchema(db)
	panic(http.ListenAndServe(":"+proxyPort, makeProxyHandler(db)))
}

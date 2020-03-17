package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

// Synchronize the information passed to use from sandstorm-http-bridge with
// that in the fossil database.
//
// Updates the basic auth info of the request to use the fossil username & password.
func syncUser(req *http.Request, db *sql.DB) {
	tx, err := db.BeginTx(req.Context(), nil)
	chkfatal(err)

	// Will have no effect if we've already committed by the time this runs:
	defer tx.Rollback()

	sandstormUid := req.Header.Get("X-Sandstorm-User-Id")
	if sandstormUid == "" {
		return
	}

	var username, password string

	// Try to find an existing user:
	err = tx.QueryRow(`
		SELECT login, pw
		FROM user, sandstorm_user
		WHERE uid = fossil_uid
		AND sandstorm_uid = ?
	`, sandstormUid).Scan(&username, &password)
	if err != nil {
		// TODO: check the specific error.
		username, password = createUser(req, tx, sandstormUid)
	}
	req.SetBasicAuth(username, password)
	chkfatal(tx.Commit())
}

func createUser(req *http.Request, tx *sql.Tx, sandstormUid string) (username, password string) {
	caps := getUserCaps(req)
	handle := req.Header.Get("X-Sandstorm-Preferred-Handle")
	username = handle
	password = generatePassword()

	userExists := func(name string) bool {
		var count int
		chkfatal(tx.QueryRow(`
			SELECT count(*)
			FROM user
			WHERE login = ?
		`, name).Scan(&count))
		return count > 0
	}

	for i := 2; userExists(username); i++ {
		username = fmt.Sprintf("%s-%d", handle, i)
	}

	res, err := tx.Exec(`
		INSERT INTO
		user(login, pw, cap, info, mtime)
		VALUES (?, ?, ?, ?, ?)
	`, username, password, caps, "", time.Now())
	chkfatal(err)
	fossilUid, err := res.LastInsertId()
	chkfatal(err)
	_, err = tx.Exec(`
		INSERT INTO
		sandstorm_user(sandstorm_uid, fossil_uid)
		VALUES(?, ?)
	`, sandstormUid, fossilUid)
	chkfatal(err)
	return
}

func getUserCaps(req *http.Request) string {
	// TODO: return a value actually based on the headers.
	return "u"
}

func generatePassword() string {
	var buf [16]byte
	_, err := rand.Read(buf[:])
	chkfatal(err)
	return fmt.Sprintf("%", buf)
}

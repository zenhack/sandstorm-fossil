package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

func makeProxyHandler(db *sql.DB) http.Handler {
	errDontRedirect := errors.New("Don't Redirect")
	client := &http.Client{
		// Let the client handle redirects itself:
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return errDontRedirect
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic in request handler: %s\n%s", r, debug.Stack())
				w.WriteHeader(500)
				w.Write([]byte("Internal Server Error\n"))
			}
		}()
		syncUser(req, db)

		req.Host = fossilAddr
		req.URL.Host = req.Host
		req.URL.Scheme = "http"
		// Server only field; we get an error if we try to pass this to a client:
		req.RequestURI = ""

		resp, err := client.Do(req)
		if err != nil && errors.Is(err, errDontRedirect) {
			err = nil
		}
		chkfatal(err)
		//defer resp.Body.Close()

		hdr := w.Header()
		for k, v := range resp.Header {
			hdr[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		_, err = io.Copy(io.MultiWriter(w, os.Stdout), resp.Body)

		if err != nil {
			// Not really anything we can do about this but log it.
			log.Println("Error copying request body from fossil to client:", err)
		}
	})
}

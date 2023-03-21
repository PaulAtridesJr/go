package main

import (
	"context"
	"errors"

	//"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
    "dummy"
	"github.com/google/uuid"
)

var store sync.Map

func profileAvailable(w http.ResponseWriter, r *http.Request) {
	log.Print("Profile available request\n")
	io.WriteString(w, "{\"result\":true}")
}

func startBioSession(w http.ResponseWriter, r *http.Request) {
	log.Print("Start bio session request")
	// take callbackURL from body
	// {"person":{"pid":"rshb_admin","pidType":"tab_num"},
	// "callbackURL":"https://pr-rshb-app-101.rshb-tc.local/auth/login?signin=c0f58476859a02ad63d1308479afde1e",

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read body: %s\n", err)
	}

	sid_t := "1234_" + uuid.New().String() + "_bio_sid"
	log.Print(sid_t)
	sid := "rshb_admin_bio_sid"
	sidTTL := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)

	parts := strings.Split(string(body), ",")
	for _, part := range parts {
		if strings.HasPrefix(part, "\"callbackURL") {
			//log.Printf("%s\n", part)
			s := strings.ReplaceAll(part, "\"", "")
			callbackURL := strings.ReplaceAll(s, "callbackURL:", "")
			log.Printf("%s\n", callbackURL)
			store.Store(sid, callbackURL)
			break
		}
	}

	io.WriteString(w, "{\"sid\":\""+sid+"\",\"sidTTL\":"+sidTTL+",\"redirectURL\":\"\"}")
}

func bioSessionResult(w http.ResponseWriter, r *http.Request) {
	log.Print("Get bio session result request")
	io.WriteString(w, "{\"data\":{\"person\":{\"pid\":\"1\",\"pidType\":\"tab_num\"},\"matchResult\":true},\"status\":0}")
}

func uiRedirect(w http.ResponseWriter, r *http.Request) {
	//corrID := r.Header.Get("X-Correlation-ID")

	// get from URL Path
	sid := "rshb_admin_bio_sid"

	cUrl, _ := store.LoadAndDelete(sid)

	log.Printf("UI redirect to '%s'", cUrl.(string))

	http.Redirect(w, r, cUrl.(string), http.StatusSeeOther)
}

func main() {
	err := http.ListenAndServe(":3333", DummyServe)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/bio/profile/1/verification/available", profileAvailable)
	mux.HandleFunc("/api/v1/bio/2fa/verification", startBioSession)
	mux.HandleFunc("/api/v1/bio/2fa/verification/rshb_admin_bio_sid/result", bioSessionResult)
	mux.HandleFunc("/ui/2fa/verification/rshb_admin_bio_sid", uiRedirect)

	ctx, cancelCtx := context.WithCancel(context.Background())
	// 10.2.7.24:9000
	// :9000
	// 192.168.4.91:9000
	serverOne := &http.Server{
		Addr:    ":9000",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	log.Print("starting\n")
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("server one closed\n")
		} else if err != nil {
			log.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}

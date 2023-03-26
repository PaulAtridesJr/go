package advanced

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var debug bool

type bioSessionInfo struct {
	PID         string
	CallbackURL string
	MatchResult bool
}

var store map[string]bioSessionInfo

func AdvancedServe(debugmode bool) func(http.ResponseWriter, *http.Request) {
	debug = debugmode
	store = make(map[string]bioSessionInfo)
	return serve
}

func serve(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	switch {
	case n == 7:
		if p[6] == "available" {
			// GET /api/v1/bio/profile/1/verification/available"
			profileAvailable(w, r, p[4])
		} else {
			// GET /api/v1/bio/2fa/verification/rshb_admin_bio_sid/result
			bioSessionResult(w, r, p[5])
		}
	case n == 5:
		// POST /api/v1/bio/2fa/verification
		startBioSession(w, r)
	case n == 4:
		// /ui/2fa/verification/rshb_admin_bio_sid
		uiRedirect(w, r, p[3])
		return
	}
}

func profileAvailable(w http.ResponseWriter, r *http.Request, pid string) {
	print(fmt.Sprintf("Profile available request for persion with ID: '%s'", pid))
	io.WriteString(w, "{\"result\":true}")
}

func startBioSession(w http.ResponseWriter, r *http.Request) {
	print("Start bio session request")
	// take callbackURL from body
	// {"person":{"pid":"rshb_admin","pidType":"tab_num"},
	// "callbackURL":"https://pr-rshb-app-101.rshb-tc.local/auth/login?signin=ef1222bdf7d5850cc1b3267a7055790d",
	// "userDevice":{
	// "ipv4":"1.1.1.1",
	// "userAgent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
	// "language":"RU",
	// "screen":{"colorDepth":0.0,"availHeight":0.0,"availWidth":0.0}}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	requestData := struct {
		Person      Person `json:"person"`
		CallbackURL string `json:"callbackURL"`
		UserDevice  struct {
			IPV4      string `json:"ipv4"`
			UserAgent string `json:"userAgent"`
			Language  string `json:"language"`
			Screen    struct {
				ColorDepth  float32 `json:"colorDepth"`
				AvailHeight float32 `json:"availHeight"`
				AvailWidth  float32 `json:"availWidth"`
			} `json:"screen"`
		} `json:"userDevice"`
	}{}
	err := decoder.Decode(&requestData)
	if err != nil {
		log.Println(err)
		return
	}

	sid := requestData.Person.Pid + "_" + uuid.New().String() + "_bio_sid"

	store[sid] = bioSessionInfo{PID: requestData.Person.Pid, CallbackURL: requestData.CallbackURL}

	m := struct {
		Sid         string `json:"sid"`
		SidTTL      string `json:"sidTTL"`
		RedirectURL string `json:"redirectURL"`
	}{
		Sid:         sid,
		SidTTL:      strconv.FormatInt(time.Now().Add(time.Hour*10).UTC().Unix()*1000, 10),
		RedirectURL: "",
	}

	res, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, fmt.Sprintf("%s\n", res))
	//io.WriteString(w, "{\"sid\":\""+sid+"\",\"sidTTL\":"+sidTTL+",\"redirectURL\":\"\"}")
}

type Person struct {
	Pid     string `json:"pid"`
	PidType string `json:"pidType"`
}

type Data struct {
	Person      Person `json:"person"`
	MatchResult bool   `json:"matchResult"`
}

type BioSessionResult struct {
	Data   Data `json:"data"`
	Status int  `json:"status"`
}

func bioSessionResult(w http.ResponseWriter, r *http.Request, sid string) {
	bio, ok := store[sid]
	if ok {
		print(fmt.Sprintf("Get bio session result for session ID: '%s'", sid))

		m := BioSessionResult{
			Data: Data{
				Person: Person{
					Pid:     bio.PID,
					PidType: "tab_num",
				},
				MatchResult: bio.MatchResult,
			},
			Status: 0,
		}

		res, err := json.Marshal(m)
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(w, fmt.Sprintf("%s\n", res))
		delete(store, sid)
		//io.WriteString(w, "{\"data\":{\"person\":{\"pid\":\"1\",\"pidType\":\"tab_num\"},\"matchResult\":true},\"status\":0}")
	} else {
		print(fmt.Sprintf("Can't get bio session result for session ID: '%s'", sid))
	}
}

func uiRedirect(w http.ResponseWriter, r *http.Request, sid string) {
	bio, ok := store[sid]
	if ok {
		print(fmt.Sprintf("UI redirect to '%s'", bio.CallbackURL))
		bio.MatchResult = true
		store[sid] = bio
		http.Redirect(w, r, bio.CallbackURL, http.StatusSeeOther)
	} else {
		print(fmt.Sprintf("Can't UI redirect. SID '%s' not found", sid))
	}

}

func print(text string) {
	if debug {
		log.Println(text)
	}
}

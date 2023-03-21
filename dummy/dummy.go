package dummy

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func DummyServe(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")[1:]
	n := len(p)

	switch {
	case n == 7:
		if p[6] == "available" {
			// GET /api/v1/bio/profile/1/verification/available"
			profileAvailable(w, r)
		} else {
			// GET /api/v1/bio/2fa/verification/rshb_admin_bio_sid/result
			bioSessionResult(w, r)
		}
	case n == 5:
		// POST /api/v1/bio/2fa/verification
		startBioSession(w, r)

		return
	}
}

func profileAvailable(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{\"result\":true}")
}

func startBioSession(w http.ResponseWriter, r *http.Request) {
	sid := uuid.New().String() + "_bio_sid"
	sidTTL := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	io.WriteString(w, "{\"sid\":\""+sid+"\",\"sidTTL\":"+sidTTL+",\"redirectURL\":\"\"}")
}

func bioSessionResult(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{\"data\":{\"person\":{\"pid\":\"1\",\"pidType\":\"tab_num\"},\"matchResult\":true},\"status\":0}")
}

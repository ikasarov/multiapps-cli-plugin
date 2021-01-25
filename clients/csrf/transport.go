package csrf

import (
	"net/http"
	"time"

	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/log"
	"github.com/jinzhu/copier"
)

type Csrf struct {
	Header              string
	Token               string
	IsInitialized       bool
	NonProtectedMethods map[string]bool
}

type Cookies struct {
	Cookies []*http.Cookie
}

type Transport struct {
	Transport http.RoundTripper
	Csrf      *Csrf
	Cookies   *Cookies
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := http.Request{}
	copier.Copy(&req2, req)

	if t.Cookies != nil {
		UpdateCookiesIfNeeded(t.Cookies.Cookies, &req2)
	}

	tokenNotYetInitialized := t.Csrf.IsInitialized == false
	csrfTokenManager := NewDefaultCsrfTokenUpdater(&t, &req2, NewDefaultCsrfTokenFetcher(&t))

	err := csrfTokenManager.updateCsrfToken()
	if err != nil {
		return nil, err
	}
	if t.Csrf.IsInitialized {
		log.Tracef("Sending a protected request with CSRF '" + t.Csrf.Token + "'\n")
	}

	if tokenNotYetInitialized && t.Csrf.IsInitialized {
		log.Tracef("Sleeping before first protected request!\n")
		time.Sleep(time.Second * 90)
	}

	log.Tracef("Sending a request with CSRF '" + req2.Header.Get("X-Csrf-Header") + " / " + req2.Header.Get("X-Csrf-Token") + "'\n")
	log.Tracef("The sticky-session headers are: " + prettyPrintCookies(req2.Cookies()) + "\n")
	res, err := t.Transport.RoundTrip(&req2)
	if err != nil {
		return nil, err
	}
	isRetryNeeded, err := csrfTokenManager.isRetryNeeded(res)
	if err != nil {
		return nil, err
	}

	if isRetryNeeded {
		log.Tracef("Response code '" + string(res.StatusCode) + "' from bad token\n")
		log.Tracef("Will retry with newer CSRF!\n")
		return nil, &ForbiddenError{}
	}

	return res, err
}

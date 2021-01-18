package csrf

import (
	"net/http"

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
	originalTransport http.RoundTripper
	csrf              *Csrf
	cookies           *Cookies
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := http.Request{}
	copier.Copy(&req2, req)

	if t.cookies != nil {
		UpdateCookiesIfNeeded(t.cookies.Cookies, &req2)
	}

	csrfTokenManager := NewDefaultCsrfTokenManager(&t, &req2)

	err := csrfTokenManager.updateToken()
	if err != nil {
		return nil, err
	}

	res, err := t.originalTransport.RoundTrip(&req2)
	if err != nil {
		return nil, err
	}
	tokenWasRefreshed, err := csrfTokenManager.refreshTokenIfNeeded(res)
	if err != nil {
		return nil, err
	}

	if tokenWasRefreshed {
		return nil, &ForbiddenError{}
	}

	return res, err
}

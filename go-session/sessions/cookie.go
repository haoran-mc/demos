package sessions

import "net/http"

// newCookieFromOptions returns an http.Cookie with the options set.
// name:   sessionName    string
// value:  sessionValue   [interface{}]interface{}
func newCookieFromOptions(name, value string, options *Options) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

/*
cookie := &http.Cookie{
	Name:     "token",
	Value:    "root",
	Path:     "/",
	MaxAge:   300,
	Secure:   true,
	HttpOnly: true,
	Expires:  time.Now().Add(time.Duration(300) * time.Second),
}
*/

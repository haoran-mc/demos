package sessions

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Store is an interface for custom session stores.
//
// See CookieStore and FilesystemStore for examples.
type Store interface {
	// Get should return a cached session.
	Get(r *http.Request, name string) (*Session, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(r *http.Request, name string) (*Session, error)

	// Save should persist session to the underlying store implementation.
	Save(r *http.Request, w http.ResponseWriter, s *Session) error
}

// CookieStore ----------------------------------------------------------------

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	codec   Codec
	Options *Options
}

// NewCookieStore returns a new CookieStore
//
// Only one key is supported, it's used for authentication.
//
// It is recommended to use an authentication key with 32 or 64 bytes.
func NewCookieStore(key []byte) *CookieStore {
	cs := &CookieStore{
		codec: CodecFromKey(key),
		Options: &Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
	}
	return cs
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (s *CookieStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).GetSession(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (s *CookieStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		// sessionName was found in the request
		err = s.codec.Decode(name, c.Value, &session.Values) // DECODE
		if err == nil {
			session.IsNew = false
		}
	}
	return session, err
}

// Save adds a single session to the response.
func (s *CookieStore) Save(r *http.Request, w http.ResponseWriter,
	session *Session) error {
	encoded, err := s.codec.Encode(session.Name(), session.Values) // ENCODE
	if err != nil {
		return err
	}
	http.SetCookie(w, NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// MaxAge sets the maximum age for the store. Individual sessions
// can be deleted by setting Options.MaxAge = -1 for that session.
func (s *CookieStore) MaxAge(age int) {
	s.Options.MaxAge = age
}

// FilesystemStore ------------------------------------------------------------

var fileMutex sync.RWMutex

// FilesystemStore stores sessions in the filesystem.
//
// It also serves as a reference for custom stores.
//
// This store is still experimental and not well tested. Feedback is welcome.
type FilesystemStore struct {
	codec   Codec
	Options *Options
	path    string
}

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, key []byte) *FilesystemStore {
	if path == "" {
		path = os.TempDir()
	}
	fs := &FilesystemStore{
		codec: CodecFromKey(key),
		Options: &Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		path: path,
	}

	fs.MaxAge(fs.Options.MaxAge)
	return fs
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *FilesystemStore) Get(r *http.Request, name string) (*Session, error) {
	return GetRegistry(r).GetSession(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *FilesystemStore) New(r *http.Request, name string) (*Session, error) {
	session := NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true
	var err error
	if c, errCookie := r.Cookie(name); errCookie == nil {
		err = s.codec.Decode(name, c.Value, &session.ID) // DECODE
		if err == nil {
			err = s.load(session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store path. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *FilesystemStore) Save(r *http.Request, w http.ResponseWriter,
	session *Session) error {
	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.erase(session); err != nil {
			return err
		}
		http.SetCookie(w, NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strconv.FormatInt(time.Now().UnixNano(), 10) // GENERATE ID
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := s.codec.Encode(session.Name(), session.ID) // ENCODE
	if err != nil {
		return err
	}
	http.SetCookie(w, NewCookie(session.Name(), encoded, session.Options))
	return nil
}

// MaxAge sets the maximum age for the store. Individual sessions
// can be deleted by setting Options.MaxAge = -1 for that session.
func (s *FilesystemStore) MaxAge(age int) {
	s.Options.MaxAge = age
}

// save writes encoded session.Values to a file.
func (s *FilesystemStore) save(session *Session) error {
	encoded, err := s.codec.Encode(session.Name(), session.Values) // ENCODE
	if err != nil {
		return err
	}
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return os.WriteFile(filename, []byte(encoded), 0600)
}

// load reads a file and decodes its content into session.Values.
func (s *FilesystemStore) load(session *Session) error {
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.RLock()
	defer fileMutex.RUnlock()
	fdata, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = s.codec.Decode(session.Name(), string(fdata), // DECODE
		&session.Values); err != nil {
		return err
	}
	return nil
}

// delete session file
func (s *FilesystemStore) erase(session *Session) error {
	filename := filepath.Join(s.path, "session_"+session.ID)

	fileMutex.RLock()
	defer fileMutex.RUnlock()

	err := os.Remove(filename)
	return err
}

package Proxy

import (
	"net"

	"github.com/Acebond/ReverseSocks5/statute"
)

// AuthContext A Request encapsulates authentication state provided
// during negotiation
type AuthContext struct {
	// Provided auth method
	Method uint8
	// Payload provided during negotiation.
	// Keys depend on the used auth method.
	// For UserPass auth contains username/password
	// Payload map[string]string
}

// Authenticator provide auth
type Authenticator interface {
	Authenticate(conn net.Conn) error
	GetCode() uint8
}

// NoAuthAuthenticator is used to handle the "No Authentication" mode
type NoAuthAuthenticator struct{}

// GetCode implement interface Authenticator
func (a NoAuthAuthenticator) GetCode() uint8 { return statute.MethodNoAuth }

// Authenticate implement interface Authenticator
func (a NoAuthAuthenticator) Authenticate(conn net.Conn) error {
	_, err := conn.Write([]byte{statute.VersionSocks5, statute.MethodNoAuth})
	return err
}

// UserPassAuthenticator is used to handle username/password based
// authentication
type UserPassAuthenticator struct {
	Username string
	Password string
}

// GetCode implement interface Authenticator
func (a UserPassAuthenticator) GetCode() uint8 { return statute.MethodUserPassAuth }

// Authenticate implement interface Authenticator
func (a UserPassAuthenticator) Authenticate(conn net.Conn) error {

	// reply the client to use user/pass auth
	if _, err := conn.Write([]byte{statute.VersionSocks5, statute.MethodUserPassAuth}); err != nil {
		return err
	}

	// get user and user's password
	nup, err := statute.ParseUserPassRequest(conn)
	if err != nil {
		return err
	}

	// Verify the password
	if a.Username == string(nup.User) && a.Password == string(nup.Pass) {
		if _, err = conn.Write([]byte{statute.UserPassAuthVersion, statute.AuthSuccess}); err != nil {
			return err
		}
		return nil
	}

	if _, err := conn.Write([]byte{statute.UserPassAuthVersion, statute.AuthFailure}); err != nil {
		return err
	}
	return statute.ErrUserAuthFailed
}

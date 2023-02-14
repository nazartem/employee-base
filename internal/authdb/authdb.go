// Package authdb emulates an "authentication database".
//
// It just exposes a global map of username -> bcrypt-ed password; in real life
// this would be a wrapper around a real DB table.
package authdb

import "golang.org/x/crypto/bcrypt"

var usersPasswords = map[string][]byte{
	"joe":  []byte("$2a$12$aMfFQpGSiPiYkekov7LOsu63pZFaWzmlfm1T8lvG6JFj2Bh4SZPWS"),
	"mary": []byte("$2a$12$u.Q6ehmzh.Qd4UnCM52Gq.2Ip/jQ5/XdtODV//gvLWxMonGZFWQGy"),
}

// VerifyUserPass verifies that username/password is a valid pair matching
// our userPasswords "database".
func VerifyUserPass(username, password string) bool {
	wantPass, hasUser := usersPasswords[username]
	if !hasUser {
		return false
	}
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}

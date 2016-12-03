package main

import (
    "golang.org/x/crypto/scrypt"
    "errors"
    "crypto/rand"
    _ "fmt"
    _ "encoding/hex"
    "net/mail"
    "bytes"
)

type LoginHelper struct {
    *SessionManager
    *DBHelper
    *MailServer
}

func (lh *LoginHelper) Init(sm *SessionManager, dbh *DBHelper, ms *MailServer) () {
    lh.SessionManager = sm
    lh.DBHelper = dbh
    lh.MailServer = ms
}

func (lh *LoginHelper) Login(form *LoginForm) ([]byte, error) {
    var token, salt, hash []byte
    var activated bool
    var err error
    if salt, hash, activated, err = lh.GetUserCreds(&form.User); err == nil {
        dk := make([]byte, 32)
        if dk, err = scrypt.Key([]byte(form.Password), salt, 16384, 8, 1, 32); err == nil {
            if bytes.Equal(dk, hash) {
                token = make([]byte, 32)
                if _, err = rand.Read(token); err == nil {
                    if !activated {
                        lh.RenewSession(form.User, token, ACTIVATION_TIMEOUT)
                        if err = lh.SendActivation([]string{"albeec13@gmail.com" /*form.User*/}, form.User, token); err == nil {
                            err = errors.New("Account was not activated. Please check email for activation link.")
                        }
                    } else {
                        lh.RenewSession(form.User, token, API_TIMEOUT)
                    }
                }
            } else {
                err = errors.New("Invalid email address or password")
            }
        }
    }
    return token, err
}

func (lh *LoginHelper) LoginCreate(form *LoginForm) (err error) {
    if _, err = mail.ParseAddress(form.User); err != nil {
        err = errors.New("Invalid email address")
    } else {
        salt := make([]byte, 32)
        dk := make([]byte, 32)

        if _, err = rand.Read(salt); err == nil {
            if dk, err = scrypt.Key([]byte(form.Password), salt, 16384, 8, 1, 32); err == nil {
                if _, err = lh.CreateUser(&form.User, &salt, &dk); err == nil {
                    /* fmt.Printf("Salt: %s\n", hex.EncodeToString(salt))
                     * fmt.Printf("Hash: %s\n", hex.EncodeToString(dk))
                     */
                    token := make([]byte, 32)
                    if _, err = rand.Read(token); err == nil {
                        lh.RenewSession(form.User, token, ACTIVATION_TIMEOUT)
                        err = lh.SendActivation([]string{"albeec13@gmail.com"/*form.user*/}, form.User, token)
                    }
                }
            }
        }
    }
    return err
}

func (lh *LoginHelper) LoginActivate(user string, token []byte) (err error) {
    if lh.IsTokenValid(user, token) {
        _, err = lh.ActivateUser(&user)
    } else {
        err = errors.New("Invalid activation code")
    }

    return err
}

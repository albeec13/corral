package main

import (
    "golang.org/x/crypto/scrypt"
    "errors"
    "crypto/rand"
    "fmt"
    "encoding/hex"
    "net/mail"
    "bytes"
)

type LoginHelper struct {
    *SessionManager
    *DBHelper
}

func (lh *LoginHelper) Login(form *LoginForm) ([]byte, error) {
    var token, salt, hash []byte
    var err error
    if salt, hash, err = lh.GetUserCreds(&form.User); err == nil {
        dk := make([]byte, 32)
        if dk, err = scrypt.Key([]byte(form.Password), salt, 16384, 8, 1, 32); err == nil {
            if bytes.Equal(dk, hash) {
                token = make([]byte, 32)
                _, err = rand.Read(token)
                lh.RenewSession(form.User, token)
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
            fmt.Printf("Salt: %s\n", hex.EncodeToString(salt))
            if dk, err = scrypt.Key([]byte(form.Password), salt, 16384, 8, 1, 32); err == nil {
                fmt.Printf("Hash: %s\n", hex.EncodeToString(dk))
            }
            _, err = lh.CreateUser(&form.User, &salt, &dk)
        }
    }
    return err
}

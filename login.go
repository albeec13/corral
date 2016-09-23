package main

import (
    "golang.org/x/crypto/scrypt"
    "errors"
    "crypto/rand"
    "fmt"
    "encoding/hex"
)

/*type LoginForm struct {
    User        string `form:"email" json:"email" binding:"required"`
    Password    string `form:"password" json:"password" binding:"required"`
}*/

type LoginHelper struct {
    dbh *DBHelper
}

func (lh *LoginHelper) Login(form *LoginForm) ([]byte, error) {
    token := make([]byte, 32)
    if form.User == "herp" && form.Password == "derp" {
        _, err := rand.Read(token)
        return token, err
    } else {
        return token, errors.New("invalid email or password")
    }
}

func (lh *LoginHelper) LoginCreate(form *LoginForm) (error) {
    salt := make([]byte, 32)
    dk := make([]byte, 32)
    _, err := rand.Read(salt)
    
    if err == nil {
        fmt.Printf("Salt: %s\n", hex.EncodeToString(salt))
        if dk, err = scrypt.Key([]byte(form.Password), salt, 16384, 8, 1, 32); err == nil {
            fmt.Printf("Hash: %s\n", hex.EncodeToString(dk))
        }
    }

    return err
}

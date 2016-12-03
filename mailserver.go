package main

import (
    "net/smtp"
    "fmt"
    "encoding/hex"
)

type MailServerInfo struct {
    User    string
    Pass    string
    Server  string
    Port    string
}


type MailServer struct {
    msi     *MailServerInfo
    auth    smtp.Auth
}

func (ms *MailServer) Init (msi *MailServerInfo) () {
    ms.auth = smtp.PlainAuth("", msi.User, msi.Pass, msi.Server)
    ms.msi = msi
}

func (ms *MailServer) SendActivation (to []string, user string, code []byte) (error) {
    msg := []byte("To: ")

    for i, t := range to {
        if i == 0 {
            msg = append(msg, []byte(t)...)
        } else {
            msg = append(msg, ", "...)
            msg = append(msg, []byte(t)...)
        }
    }

    msg = append(msg, []byte("\r\n" +
        "Subject: Corral Activation Link\r\n" +
        "\r\n" +
        "Thanks for registering for Corral!\r\n" +
        "\r\n" +
        "Please click the following link to activate your account: \r\n" +
        "\r\n" +
        "https://thewalr.us/corral/API/activate/" + hex.EncodeToString([]byte(user)) + "/" + hex.EncodeToString(code) + "\r\n" +
        "\r\n" +
        "If you have received this message in error, please ignore it.\r\n" +
        "\r\n" +
        "- Corral Team")...)

    err := smtp.SendMail(ms.msi.Server + ":" + ms.msi.Port, ms.auth, ms.msi.User, to, msg)

    fmt.Println(msg)

    return err
}

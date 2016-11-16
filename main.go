package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "log"
    "encoding/hex"
)

type LoginForm struct {
    User        string `form:"email" json:"email" binding:"required"`
    Password    string `form:"password" json:"password" binding:"required"`
}

func main() {
    var dbh DBHelper
    var config ConfigFile
    var lh LoginHelper
    var sm SessionManager
    var ms MailServer
    path := "corral.conf"

    if err := config.ReadConfigFile(path); err == nil {
        if err = dbh.Open(&config); err == nil {
            var tables []string
            if tables, err = dbh.GetTables(); err == nil {
                for _, table := range tables {
                    fmt.Printf("Table: %s\n", table)
                }
            } else {
                fmt.Printf("DATABASE ERROR: %s\n", err)
            }
        } else {
            log.Fatalf("FATAL ERROR: %s\n", err)
        }
    } else {
        log.Fatalf("FATAL ERROR: %s\n", err)
    }

    fmt.Println(lh.DBHelper)

    // start mail server interface
    ms.Init(&config.MailInfo)

    // start session manager
    sm.Init()

    // start login helper
    lh.Init(&sm, &dbh)
    fmt.Println(lh)

    // Configure routes
    router := gin.Default()

    router.Static("/corral/","./html")

    router.POST("/corral/login", func(c *gin.Context) {
        var form LoginForm
        if c.Bind(&form) == nil {
            if sess_token, err := lh.Login(&form); err == nil {
                c.JSON(http.StatusOK, gin.H{"session_token" : hex.EncodeToString(sess_token)})
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"error" : err.Error()})
            }
            fmt.Println(form)
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid API access"})
        }
    })

    router.POST("/corral/createLogin", func(c *gin.Context) {
        var form LoginForm
        if c.Bind(&form) == nil {
            if err := lh.LoginCreate(&form); err == nil {
                ms.SendActivation([]string{form.User}, []byte("testcode"))
                c.JSON(http.StatusOK, gin.H{"status" : "Please check your email for a confirmation link."})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
            }
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid API access"})
        }
    })

    // Run server
    router.Run(":4569")
}

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
        if err = dbh.Open(&config.DBInfo); err == nil {
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
    lh.Init(&sm, &dbh, &ms)
    fmt.Println(lh)


    // Configure routes
    router := gin.Default()
	router.Use(CORSMiddleware())
    routerStatic := gin.Default()

    // Static routes to html
    routerStatic.Static("/","./www")

    // Remaining routes are API routes
    router.POST("/corral/API/login", func(c *gin.Context) {
        fmt.Println(c.Request.Host)
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

    router.POST("/corral/API/createLogin", func(c *gin.Context) {
        var form LoginForm
        if c.Bind(&form) == nil {
            if err := lh.LoginCreate(&form); err == nil {
                c.JSON(http.StatusOK, gin.H{"status" : "Please check your email for a confirmation link."})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error" : err.Error()})
            }
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error" : "invalid API access"})
        }
    })

    router.GET("/corral/API/activate/:token", func(c *gin.Context) {
        token := c.Param("token")
        fmt.Println("Token confirmed:" + token)
    })

    // Run servers
    go router.Run(":4569")
    routerStatic.Run(":4570")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
        origin := ""
        if len(c.Request.Header["Origin"]) > 0 {
            origin = c.Request.Header["Origin"][0]
            fmt.Println("Origin: " + origin)
        }
        allowed := false;
        whiteList := []string{"https://thewalr.us","https://www.thewalr.us","https://corral.thewalr.us"}

        for _, dom := range whiteList {
            if dom == origin {
                allowed = true;
            }
        }

        if(allowed) {
		    c.Writer.Header().Set("Access-Control-Allow-Origin",origin)
            c.Writer.Header().Set("Access-Control-Allow-Credentials","true") 
            c.Writer.Header().Set("Vary","Origin")
            c.Writer.Header().Set("Access-Control-Expose-Headers","Location")

		    if c.Request.Method == "OPTIONS" {
			    c.Writer.Header().Set("Access-Control-Allow-Headers","Authorization,Content-Type,Accept,Origin,User-Agent,DNT,Cache-Control,X-Mx-ReqToken,Keep-Alive,X-Requested-With,If-Modified-Since")
                c.Writer.Header().Set("Access-Control-Allow-Methods","GET, POST, OPTIONS")
                c.Writer.Header().Set("Access-Control-Max-Age","1728000")
                c.Writer.Header().Set("Content-Length","0")
                c.Writer.Header().Set("Content-Type","text/plain charset=UTF-8")
			    c.AbortWithStatus(204)
			    return
		    }
		    c.Next()
        } else {
            c.AbortWithStatus(404)
        }
    }
}

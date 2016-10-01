package main

import (
    "time"
    "sync"
    "bytes"
    "fmt"
)

type Session struct {
    token []byte
    expTimeUnix int64
}

type SessionManager struct{
    mut sync.RWMutex
    userSessions map[string]Session
}

func (sm *SessionManager) Init() () {
    sm.mut.Lock()
    defer sm.mut.Unlock()
    sm.userSessions = make(map[string]Session)
    sm.startCleanup()
}

func (sm *SessionManager) startCleanup() () {
    go func() {
        for true {
            now := time.Now().Unix()
            sm.mut.RLock()
            for user, sess := range sm.userSessions {
                if now > sess.expTimeUnix {
                    sm.mut.RUnlock()
                    sm.mut.Lock()
                        fmt.Printf("Cleaned up %s.\n", user)
                        delete(sm.userSessions, user)
                    sm.mut.Unlock()
                    sm.mut.RLock()
                }
            }
            sm.mut.RUnlock()

            // clean up sessions every 15 seconds
            time.Sleep(15 * time.Second)
        }
    }()
}

func (sm *SessionManager) RenewSession(user string, token []byte) () {
    sm.mut.Lock()
    defer sm.mut.Unlock()

    // Set session expiration time to now + 10 minutes (Unix time in seconds + 10 * 60s)
    sess := Session{token, time.Now().Unix() + 10 * 60}
    sm.userSessions[user] = sess
    fmt.Println(sm.userSessions)
}

func (sm *SessionManager) IsTokenValid(user string, token []byte) (ret bool) {
    sm.mut.RLock()
    defer sm.mut.RUnlock()

    if sess, exists := sm.userSessions[user]; exists == true {
        if bytes.Equal(sess.token, token) {
            ret = true
        }
    }
    ret = false
    return ret
}

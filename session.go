package main

import (
    "time"
    "sync"
    "bytes"
)

type TimeMinutes int64

const (
    API_TIMEOUT TimeMinutes = 10
    ACTIVATION_TIMEOUT TimeMinutes = (24 * 60)
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

func (sm *SessionManager) RenewSession(user string, token []byte, timeout TimeMinutes) () {
    sm.mut.Lock()
    defer sm.mut.Unlock()

    tmrLength := time.Now().Unix() + (int64(timeout) * 60)
    sess := Session{token, tmrLength}
    sm.userSessions[user] = sess
}

func (sm *SessionManager) IsTokenValid(user string, token []byte) (ret bool) {
    sm.mut.RLock()
    defer sm.mut.RUnlock()

    if sess, exists := sm.userSessions[user]; exists == true {
        if bytes.Equal(sess.token, token) {
            ret = true
        }
    }
    return ret
}

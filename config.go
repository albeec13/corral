package main

import (
    "io/ioutil"
    "encoding/json"
)

type ConfigFile struct {
    DBInfo   DatabaseInfo
    MailInfo MailServerInfo
}

func (cf *ConfigFile) ReadConfigFile(path string) (error) {
    file, err  := ioutil.ReadFile(path)
    if file != nil {
        err = json.Unmarshal(file, cf)
    }
    return err
}

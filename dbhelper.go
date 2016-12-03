package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "errors"
)

type DatabaseInfo struct {
    User    string
    Pass    string
    Name    string
    Tables  []string
}

type DBHelper struct{
    db *sql.DB
}

func (dbh *DBHelper) Open(dbi *DatabaseInfo) (err error) {
    if dbh.db, err = sql.Open("mysql", dbi.User + ":" + dbi.Pass + "@/" + dbi.Name); err == nil { 
        var tables []string
        if tables, err = dbh.GetTables(); tables == nil && err == nil {
            _, err = dbh.CreateTables(dbi.Tables)
        }
    }
    return err
}

func (dbh *DBHelper) GetTables() ([]string, error) {
    rows, err := dbh.db.Query("SHOW TABLES")
    var tables []string

    if err == nil {
        for rows.Next() {
            var name string
            if err := rows.Scan(&name); err == nil {
                tables = append(tables, name)
            }
        }
    }

    return tables, err
}

func (dbh *DBHelper) GetUserCreds(user *string) (salt []byte, hash []byte, activated bool, err error) {
    salt = make([]byte, 32)
    hash = make([]byte, 32)
    activated = false 

    if err = dbh.db.QueryRow("SELECT pwdsalt, pwdhash, activated FROM users WHERE email LIKE(?)", user).Scan(&salt, &hash, &activated); err == sql.ErrNoRows {
        err = errors.New("Email address does not exist")
    }
    return salt, hash, activated, err
}

func (dbh *DBHelper) CreateUser(user *string, salt *[]byte, hash *[]byte) (result sql.Result, err error) {
    var email string

    if err = dbh.db.QueryRow("SELECT email FROM users WHERE email LIKE(?)", user).Scan(&email); err == sql.ErrNoRows {
        // no hits, so we can create this user
        result, err = dbh.db.Exec(
            "INSERT INTO users (email, pwdsalt, pwdhash, activated) VALUES (?,?,?,?)",
            user,
            salt,
            hash,
            false,
        )
    } else if err != nil {
        // some other error occured, don't need to do anything
    } else {
        // non-error, user exists, so we can't create this user. also destroy result before returning
        err = errors.New("Email address already exists")
        result = nil
    }

    return result, err
}

func (dbh *DBHelper) ActivateUser(user *string) (result sql.Result, err error) {
    result, err = dbh.db.Exec("UPDATE users SET activated=TRUE WHERE email=?", user)
    return result, err
}

func (dbh *DBHelper) DeactivateUser(user *string) (result sql.Result, err error) {
    result, err = dbh.db.Exec("UPDATE users SET activated=FALSE WHERE email=?", user)
    return result, err
}

/*func (dbh *DBHelper) LogDevice(de *DHCPEvent) (result sql.Result, err error) {
    var device_id int64

    // Look for existing device_id for this MAC address, and if we don't find it, or return error on failure
    if err = dbh.db.QueryRow("SELECT device_id FROM devices WHERE mac LIKE(?)", de.MAC).Scan(&device_id); err == sql.ErrNoRows {
        // Since no rows were returned, assume this is a new MAC and try to insert it, and return error on failure
        if result, err = dbh.db.Exec("INSERT INTO devices (mac) VALUES (?)", de.MAC); err != nil {
            return result, err
        } else {
            // On successful insert, get the ID to use in the dhcpevent table, and return error if it fails
            if device_id, err = result.LastInsertId(); err != nil {
                return result, err
            }
        }
    } else if err != nil {
        return result, err
    }
    
    // Now insert DHCPEvent details into dhcpevents table, linked by device_id to devices table
    result, err = dbh.db.Exec(
        "INSERT INTO dhcpevents (event, ip, hostname, device_id) VALUES (?, ?, ?, ?)",
        de.Event,
        de.IP,
        de.Hostname,
        device_id,
    )
    return result, err
}*/

func (dbh *DBHelper) CreateTables(tables []string) (result sql.Result, err error) {
    for _, table := range tables {
        if result, err = dbh.db.Exec(table); err != nil {
            return result, err
        }
    }
    return result, err
}

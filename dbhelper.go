package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type DBHelper struct{
    db *sql.DB
}

func (dbh *DBHelper) Open(config *ConfigFile) (err error) {
    if dbh.db, err = sql.Open("mysql", config.DBUser + ":" + config.DBPass + "@/" + config.DBName); err == nil { 
        var tables []string
        if tables, err = dbh.GetTables(); tables == nil && err == nil {
            _, err = dbh.CreateTables(config.DBTables)
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

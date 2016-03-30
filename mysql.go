package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"strings"
    _ "github.com/go-sql-driver/mysql"
)

func hello(w http.ResponseWriter, r *http.Request) {
	var mobileNo string = r.URL.Query().Get("mobile")
	fmt.Fprint(w, mobileNo)
	fmt.Fprint(w, '\n')
	db, err := sql.Open("mysql", "<username>:<password>@/<dbname>")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    rows, err := db.Query("SELECT id, message_text FROM genericsms where mobile = " + mobileNo + " order by id desc limit 1;")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Make a slice for the values
    values := make([]sql.RawBytes, len(columns))

    // rows.Scan wants '[]interface{}' as an argument, so we must copy the
    // references into such a slice
    // See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    // Fetch rows
    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }

        // Now do something with the data.
        // Here we just print each column as a string.
        var value string
        for i, col := range values {
            // Here we can check if the value is nil (NULL value)
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            sArr:= []string{columns[i], value }
            row := strings.Join(sArr,": ")
            fmt.Print(row)
            fmt.Fprint(w, row)
            fmt.Fprint(w,"\n")
        }
        fmt.Fprint(w, "-----------------------------------\n")
    }
    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

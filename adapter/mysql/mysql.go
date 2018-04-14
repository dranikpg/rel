package mysql

import (
	"database/sql"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sqlutil"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
)

// Adapter definition for mysql database.
type Adapter struct {
	*sqlutil.Adapter
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error
	adapter := &Adapter{
		&sqlutil.Adapter{
			Placeholder:   "?",
			IsOrdinal:     false,
			IncrementFunc: incrementFunc,
			ErrorFunc:     errorFunc,
		},
	}

	adapter.DB, err = sql.Open("mysql", dsn)
	return adapter, err
}

func incrementFunc(adapter sqlutil.Adapter) int {
	var variable string
	var increment int
	var err error
	if adapter.TX != nil {
		err = adapter.TX.QueryRow("SHOW VARIABLES LIKE 'auto_increment_increment';").Scan(&variable, &increment)
	} else {
		err = adapter.DB.QueryRow("SHOW VARIABLES LIKE 'auto_increment_increment';").Scan(&variable, &increment)
	}
	paranoid.Panic(err)

	return increment
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	} else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return errors.DuplicateError(e.Message, "")
	}

	return err
}

package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "utilware/sqlite"
)

func main() {
	if err := main1(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main1() error {
	// dir, err := ioutil.TempDir("", "test-")
	// if err != nil {
	// 	return err
	// }

	// defer os.RemoveAll(dir)

	db, err := sql.Open("sqlite", "./qaq.db")
	if err != nil {
		return err
	}

	if _, err = db.Exec(`
drop table if exists t;
create table t(i);
insert into t values(42), (314);
`); err != nil {
		return err
	}

	rows, err := db.Query("select i from t where i > ?", 4)
	if err != nil {
		return err
	}

	for rows.Next() {
		var i int
		if err = rows.Scan(&i); err != nil {
			return err
		}

		fmt.Println(i)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	if err = db.Close(); err != nil {
		return err
	}

	// fi, err := os.Stat(fn)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("%s size: %v\n", fn, fi.Size())
	return nil
}

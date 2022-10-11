package sqlx

import (
	"context"
	"database/sql"
	"sync"
	"utilware/sqlx/reflectx"
)

type DB struct {
	*Logger

	raw        *sql.DB
	driverName string
	unsafe     bool
	Mapper     *reflectx.Mapper
}

func NewDb(db *sql.DB, driverName string) *DB {
	return &DB{raw: db, Logger: &Logger{
		output: DefaultLoggerRawOutput,
		locker: &sync.RWMutex{},
	}, driverName: driverName, Mapper: mapper()}
}

func (db *DB) DriverName() string {
	return db.driverName
}

func Open(ctx context.Context, driverName, dataSourceName string) (*DB, error) {
	raw, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db := NewDb(raw, driverName)

	err = db.raw.PingContext(ctx)
	return db, err
}

func (db *DB) DB() *sql.DB {
	return db.raw
}

func (db *DB) Unsafe() *DB {
	return &DB{Logger: db.Logger, raw: db.raw, driverName: db.driverName, unsafe: true, Mapper: db.Mapper}
}

func (db *DB) Scan(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	r := db.Query(ctx, query, args...)
	return r.scanAny(dest, false)
}

func (db *DB) ScanRows(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, e := db.QueryRows(ctx, query, args...)
	if e != nil {
		return e
	}

	defer rows.Close()

	return scanAll(rows, dest, false)
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (row *Row) {
	defer func(l loggerDone) {
		l.done(row.err)
	}(db.Logger.append(query, args...))

	rows, err := db.raw.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, unsafe: db.unsafe, Mapper: db.Mapper}
}

func (db *DB) QueryRows(ctx context.Context, query string, args ...interface{}) (rows *Rows, e error) {
	defer func(l loggerDone) {
		l.done(e)
	}(db.Logger.append(query, args...))

	r, e := db.raw.QueryContext(ctx, query, args...)
	if e != nil {
		return nil, e
	}

	return &Rows{Rows: r, unsafe: db.unsafe, Mapper: db.Mapper}, e
}

func (db *DB) Begin(ctx context.Context, opts ...*sql.TxOptions) (*Tx, error) {
	var opt *sql.TxOptions
	if len(opts) > 1 {
		opt = opts[0]
	}

	tx, e := db.raw.BeginTx(ctx, opt)
	if e != nil {
		return nil, e
	}

	return &Tx{raw: tx, Logger: db.Logger, driverName: db.driverName, unsafe: db.unsafe, Mapper: db.Mapper}, nil
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (resp *Result, e error) {
	defer func(l loggerDone) {
		l.done(e)
	}(db.Logger.append(query, args...))

	res, e := db.raw.ExecContext(ctx, query, args...)
	if e != nil {
		return nil, e
	}

	return &Result{Result: res}, nil
}

func (db *DB) Prepare(ctx context.Context, query string) (*Stmt, error) {
	stmt, e := db.raw.PrepareContext(ctx, query)
	if e != nil {
		return nil, e
	}

	return &Stmt{raw: stmt, Logger: db.Logger, query: query, unsafe: db.unsafe, Mapper: db.Mapper}, nil
}

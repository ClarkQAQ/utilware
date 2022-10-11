package sqlx

import (
	"context"
	"database/sql"
	"time"
	"utilware/sqlx/reflectx"
)

type Tx struct {
	*Logger
	count      uint32
	duration   time.Duration
	raw        *sql.Tx
	driverName string
	unsafe     bool
	Mapper     *reflectx.Mapper
}

func (tx *Tx) Unsafe() *Tx {
	return &Tx{Logger: tx.Logger, raw: tx.raw, driverName: tx.driverName, unsafe: true, Mapper: tx.Mapper}
}

func (tx *Tx) Scan(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	r := tx.Query(ctx, query, args...)
	return r.scanAny(dest, false)
}

func (tx *Tx) ScanRows(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	rows, e := tx.QueryRows(ctx, query, args...)
	if e != nil {
		return e
	}

	defer rows.Close()

	return scanAll(rows, dest, false)
}

func (tx *Tx) Query(ctx context.Context, query string, args ...interface{}) (row *Row) {
	tx.count++
	defer func(t time.Time, l loggerDone) {
		tx.duration += time.Since(t)
		l.done(row.err)
	}(time.Now(), tx.Logger.append(query, args...))

	rows, err := tx.raw.QueryContext(ctx, query, args...)
	return &Row{rows: rows, err: err, unsafe: tx.unsafe, Mapper: tx.Mapper}
}

func (tx *Tx) QueryRows(ctx context.Context, query string, args ...interface{}) (rows *Rows, e error) {
	tx.count++
	defer func(t time.Time, l loggerDone) {
		tx.duration += time.Since(t)
		l.done(e)
	}(time.Now(), tx.Logger.append(query, args...))

	r, e := tx.raw.QueryContext(ctx, query, args...)
	if e != nil {
		return nil, e
	}

	return &Rows{Rows: r, unsafe: tx.unsafe, Mapper: tx.Mapper}, e
}

func (tx *Tx) Exec(ctx context.Context, query string, args ...interface{}) (resp *Result, e error) {
	tx.count++
	defer func(t time.Time, l loggerDone) {
		tx.duration += time.Since(t)
		l.done(e)
	}(time.Now(), tx.Logger.append(query, args...))

	res, e := tx.raw.ExecContext(ctx, query, args...)
	if e != nil {
		return nil, e
	}

	return &Result{Result: res}, nil
}

func (tx *Tx) Prepare(ctx context.Context, query string) (*Stmt, error) {
	stmt, e := tx.raw.PrepareContext(ctx, query)
	if e != nil {
		return nil, e
	}

	return &Stmt{raw: stmt, Logger: tx.Logger, query: query, unsafe: tx.unsafe, Mapper: tx.Mapper}, nil
}

func (tx *Tx) Count() uint32 {
	return tx.count
}

func (tx *Tx) Duration() time.Duration {
	return tx.duration
}

func (tx *Tx) Roolback() error {
	return tx.raw.Rollback()
}

func (tx *Tx) Commit() error {
	return tx.raw.Commit()
}

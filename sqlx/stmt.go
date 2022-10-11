package sqlx

import (
	"context"
	"database/sql"
	"utilware/sqlx/reflectx"
)

type Stmt struct {
	*Logger

	raw    *sql.Stmt
	query  string
	unsafe bool
	Mapper *reflectx.Mapper
}

func (stmt *Stmt) Unsafe() *Stmt {
	return &Stmt{Logger: stmt.Logger, raw: stmt.raw, query: stmt.query, unsafe: true, Mapper: stmt.Mapper}
}

func (stmt *Stmt) Scan(ctx context.Context, dest interface{}, args ...interface{}) error {
	r := stmt.Query(ctx, args...)
	return r.scanAny(dest, false)
}

func (stmt *Stmt) ScanRows(ctx context.Context, dest interface{}, args ...interface{}) error {
	rows, e := stmt.QueryRows(ctx, args...)
	if e != nil {
		return e
	}

	defer rows.Close()

	return scanAll(rows, dest, false)
}

func (stmt *Stmt) Query(ctx context.Context, args ...interface{}) (row *Row) {
	defer func(l loggerDone) {
		l.done(row.err)
	}(stmt.Logger.append(stmt.query, args...))

	rows, err := stmt.raw.QueryContext(ctx, args...)
	return &Row{rows: rows, err: err, unsafe: stmt.unsafe, Mapper: stmt.Mapper}
}

func (stmt *Stmt) QueryRows(ctx context.Context, args ...interface{}) (rows *Rows, e error) {
	defer func(l loggerDone) {
		l.done(e)
	}(stmt.Logger.append(stmt.query, args...))

	r, e := stmt.raw.QueryContext(ctx, args...)
	if e != nil {
		return nil, e
	}

	return &Rows{Rows: r, unsafe: stmt.unsafe, Mapper: stmt.Mapper}, e
}

func (stmt *Stmt) Exec(ctx context.Context, args ...interface{}) (resp *Result, e error) {
	defer func(l loggerDone) {
		l.done(e)
	}(stmt.Logger.append(stmt.query, args...))

	res, e := stmt.raw.ExecContext(ctx, args...)
	if e != nil {
		return nil, e
	}

	return &Result{Result: res}, nil
}

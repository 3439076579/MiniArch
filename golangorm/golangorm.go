package golangorm

import (
	"MiniArch/golangorm/dialect"
	"MiniArch/golangorm/logger"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func Open(driver, source string) (e *Engine) {

	db, err := sql.Open(driver, source)
	if err != nil {
		logger.Error("database connect fail:", err)
		panic("database connect fail")
	}

	if err := db.Ping(); err != nil {
		logger.Error("database transform information fail:", err)
		panic("database transform information fail")
	}

	dial := dialect.GetDialect("mysql")
	e = &Engine{db: db, dialect: dial}

	logger.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		logger.Error("Fail to Close database Connection")
	}

	logger.Info("Close database connection success")

}

func (e *Engine) NewSession() *Session {
	return &Session{
		db:      e.db,
		dialect: e.dialect,
	}
}
func (e *Engine) DB() *sql.DB {
	return e.db
}

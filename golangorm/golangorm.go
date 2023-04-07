package golangorm

import (
	"MiniArch/golangorm/dialect"
	"MiniArch/golangorm/logger"
	"MiniArch/golangorm/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}
type TxFunc func(s *session.Session) (interface{}, error)

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

func (e *Engine) NewSession() *session.Session {
	return &session.Session{
		Db:      e.db,
		Dialect: e.dialect,
	}
}
func (e *Engine) DB() *sql.DB {
	return e.db
}
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}

	}()

	return f(s)

}

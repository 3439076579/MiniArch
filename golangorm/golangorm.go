package golangorm

import (
	"MiniArch/golangorm/dialect"
	"MiniArch/golangorm/logger"
	"MiniArch/golangorm/session"
	"database/sql"
	"fmt"
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

// b - a
func difference(a []string, b []string) (diff []string) {
	tmp := make(map[string]bool)
	for _, str := range a {
		tmp[str] = true
	}

	for _, str := range b {
		if _, ok := tmp[str]; ok {
			diff = append(diff, str)
		}

	}
	return

}

func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (interface{}, error) {
		s = s.Model(value)
		table := s.Table()

		if !s.HasTable() {
			logger.Error("table does not exist")
			return nil, s.CreateTable()
		}

		rows, err := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.TableName)).QueryRows()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		columns, err := rows.Columns()
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		delCol := difference(table.FieldName, columns)
		addCol := difference(columns, table.FieldName)

		if len(delCol) == 0 {
			return nil, nil
		} else {
			for _, str := range delCol {
				SQL := fmt.Sprintf("ALTER TABLE %s DROP %s", table.TableName, str)
				_, err := s.Raw(SQL).Exec()
				if err != nil {
					logger.Error(err)
					return nil, err
				}
			}
		}

		if len(addCol) == 0 {
			return nil, nil
		}
		for _, str := range addCol {
			SQL := fmt.Sprintf("ALTER TABLE %s ADD %s %s", table.TableName, str, table.GetFieldByName(str).Type)
			_, err := s.Raw(SQL).Exec()
			if err != nil {
				logger.Error(err)
				return nil, err
			}
		}
		return nil, nil
	})

	return err
}

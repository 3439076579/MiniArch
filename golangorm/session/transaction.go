package session

import (
	"MiniArch/golangorm/logger"
)

func (s *Session) Begin() error {
	logger.Info("transaction begin")
	tx, err := s.Db.Begin()
	if err != nil {
		logger.Error(err)
		return err
	}
	s.tx = tx
	return nil
}

func (s *Session) Commit() error {
	logger.Info("transaction commit")
	if s.tx == nil {
		panic("transaction has been not started")
	}
	if err := s.tx.Commit(); err != nil {
		logger.Error(err)
		return err
	}

	return nil

}

func (s *Session) Rollback() error {
	logger.Info("transaction rollback")
	if s.tx == nil {
		panic("transaction has been not started")
	}
	if err := s.tx.Rollback(); err != nil {
		logger.Error(err)
		return err
	}

	return nil

}

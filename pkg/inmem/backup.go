package inmem

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/aridae/go-metrics-store/pkg/logger"
	"io"
	"time"
)

func (s *Storage[Key, Value]) InitBackup(
	ctx context.Context,
	backupFile file,
	backupInterval time.Duration,
	registerTypes map[string]any,
) error {
	s.backupFile = backupFile
	s.backupInterval = backupInterval

	// register types for gob serialization
	for name, val := range registerTypes {
		gob.RegisterName(name, val)
	}

	go s.runBackupLoop(ctx)

	return nil
}

func (s *Storage[Key, Value]) runBackupLoop(ctx context.Context) {
	ticker := time.NewTicker(s.backupInterval)
	defer func() {
		s.shutBackup()
	}()

	for {
		select {
		case <-ticker.C:
			if err := s.dumpBackup(); err != nil {
				logger.Errorf("[timeseriesstorage.Storage.runBackupLoop][CRITICAL] failed to dump data to backup file: %v", err)
			}
		case <-ctx.Done():
			logger.Infof("stopping backup service downstreams...")
			return
		}
	}
}

func (s *Storage[Key, Value]) shutBackup() {
	s.backupFileMu.Lock()
	defer s.backupFileMu.Unlock()

	logger.Infof("closing backup file...")

	err := s.backupFile.Close()
	if err != nil {
		logger.Errorf("failed to close backup file: %v", err)
	}

	logger.Infof("backup is shut")
}

func (s *Storage[Key, Value]) dumpBackup() error {
	s.backupFileMu.Lock()
	s.storeMu.RLock()
	defer s.backupFileMu.Unlock()
	defer s.storeMu.RUnlock()

	s.backupFile.Truncate(0) //nolint:errcheck
	s.backupFile.Seek(0, 0)  //nolint:errcheck

	return gob.NewEncoder(s.backupFile).Encode(s.store)
}

func (s *Storage[Key, Value]) LoadFromBackup() error {
	if s.backupFile == nil {
		return fmt.Errorf("no backup file found, make sure to call InitBackup() method to init backing options before loading from backup file")
	}

	s.backupFileMu.RLock()
	s.storeMu.Lock()
	defer s.backupFileMu.RUnlock()
	defer s.storeMu.Unlock()

	newStore := s.store
	err := gob.NewDecoder(s.backupFile).Decode(&newStore)
	if err == io.EOF {
		logger.Infof("backup file is empty, nothing to load")
		return nil
	}

	_ = newStore

	return err
}

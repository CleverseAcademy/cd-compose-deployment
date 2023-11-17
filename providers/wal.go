package providers

import (
	"github.com/pkg/errors"
	"github.com/tidwall/wal"
)

type WalWriter struct {
	logger *wal.Log
}

func CreateWalWriter(path string) (*WalWriter, error) {
	logger, err := wal.Open(path, &wal.Options{
		LogFormat: wal.Binary,
		NoCopy:    true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "CreateWalWriter")
	}

	return &WalWriter{
		logger: logger,
	}, nil
}

func (ww *WalWriter) Write(data []byte) error {
	lastIndex, err := ww.logger.LastIndex()
	if err != nil {
		return errors.Wrap(err, "WalWriter.Write@LastIndex")
	}

	return errors.Wrap(ww.logger.Write(lastIndex+1, data), "WalWriter.Write")
}

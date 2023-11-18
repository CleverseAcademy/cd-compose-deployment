package providers

import (
	"github.com/pkg/errors"
	"github.com/tidwall/wal"
)

type WalWriter struct {
	logger  *wal.Log
	entropy *Entropy
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

func (ww *WalWriter) Write(data []byte) (int, error) {
	lastIndex, err := ww.logger.LastIndex()
	if err != nil {
		return 0, errors.Wrap(err, "LastIndex")
	}

	if ww.entropy != nil {
		err = ww.entropy.Update(data)

		if err != nil {
			return 0, errors.Wrap(err, "entropy.Update")
		}
	}

	return len(data), errors.Wrap(ww.logger.Write(lastIndex+1, data), "logger.Write")
}

func (ww *WalWriter) RegisterEntropyObserver(e *Entropy) error {
	ww.entropy = e

	firstWalIndex, err := ww.logger.FirstIndex()
	if err != nil {
		panic(err)
	}

	lastWalIndex, err := ww.logger.LastIndex()
	if err != nil {
		panic(err)
	}

	if lastWalIndex == firstWalIndex {
		return nil
	}

	for i := firstWalIndex; i <= lastWalIndex; i++ {
		data, err := ww.logger.Read(i)
		if err != nil {
			return errors.Wrap(err, "Read")
		}

		err = e.Update(data)
		if err != nil {
			return errors.Wrap(err, "Entropy.Update")
		}
	}

	return nil
}

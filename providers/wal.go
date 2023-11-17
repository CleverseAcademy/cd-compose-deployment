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

func (ww *WalWriter) Write(data []byte) error {
	lastIndex, err := ww.logger.LastIndex()
	if err != nil {
}

func (ww *WalWriter) RegisterEntropy(e *Entropy) error {
	firstWalIndex, err := ww.logger.FirstIndex()
	if err != nil {
		return errors.Wrap(err, "WalWriter.GetEntropy@FirstIndex")
	}

	lastWalIndex, err := ww.logger.LastIndex()
	if err != nil {
		return errors.Wrap(err, "WalWriter.GetEntropy@LastIndex")
	}

	if lastWalIndex == firstWalIndex {
		return nil
	}

	for i := firstWalIndex; i <= lastWalIndex; i++ {
		data, err := ww.logger.Read(i)
		if err != nil {
			return errors.Wrap(err, "WalWriter.GetEntropy@Read")
		}

		err = e.Update(data)
		if err != nil {
			return errors.Wrap(err, "WalWriter.GetEntropy@Entropy.Update")
		}
	}

	ww.entropy = e
	return nil
}

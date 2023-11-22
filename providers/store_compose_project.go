package providers

import (
	"fmt"
	"os"
	"time"

	"github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
)

type IArgsStoreComposeProject struct {
	BackupFile string
	OldProject *types.Project

	TargetFile string
	NewProject *types.Project
}

func StoreComposeProject(args IArgsStoreComposeProject) error {
	oData, err := args.OldProject.MarshalYAML()
	if err != nil {
		return errors.Wrap(err, "OldProject.MarshalYAML")
	}

	bckFile, err := os.OpenFile(args.BackupFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return errors.Wrap(err, "OpenBackupFile")
	}

	defer bckFile.Close()

	_, err = bckFile.Write([]byte(fmt.Sprintf("\n---\n# StoreComposeProject at %s", time.Now())))
	if err != nil {
		return errors.Wrap(err, "bckFile.Write")
	}
	_, err = bckFile.Write(oData)
	if err != nil {
		return errors.Wrap(err, "bckFile.Write")
	}

	nData, err := args.NewProject.MarshalYAML()
	if err != nil {
		return errors.Wrap(err, "NewProject.MarshalYAML")
	}
	cfgFile, err := os.OpenFile(args.TargetFile, os.O_WRONLY, 0o644)
	if err != nil {
		return errors.Wrap(err, "OpenTargetFile")
	}

	defer cfgFile.Close()

	_, err = cfgFile.WriteAt(nData, 0)
	if err != nil {
		return errors.Wrap(err, "cfgFile.Write")
	}

	return nil
}

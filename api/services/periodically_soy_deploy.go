package services

import (
	"fmt"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/providers"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
)

func (s Service) PeriodicallySoyDeploy(clnt *client.Client, composeAPI api.Service, initialPrj *types.Project, backupPath string) {
	currentprj := initialPrj
	for {
		time.Sleep(config.AppConfig.DeployInterval)

		for _, svc := range currentprj.Services {
			nextPrj, err := s.ExecuteServiceDeployments.Execute(
				usecases.IArgsExecuteServiceDeployments{
					ComposeAPI:                composeAPI,
					ServiceName:               entities.ServiceName(svc.Name),
					LogDeploymentDoneEvent:    s.LogDeploymentDoneEvent,
					LogDeploymentFailureEvent: s.LogDeploymentFailureEvent,
					LogDeploymentSkippedEvent: s.LogDeploymentSkippedEvent,
				},
			)
			if err != nil {
				fmt.Println(err)
			} else {
				err := providers.StoreComposeProject(providers.IArgsStoreComposeProject{
					BackupFile: backupPath,
					OldProject: currentprj,
					TargetFile: config.AppConfig.ComposeFile,
					NewProject: nextPrj,
				})
				if err != nil {
					panic(err)
				}

				removeList, err := utils.RemoveImage(clnt, svc.Image)
				if err != nil {
					fmt.Printf("Prune old image error: %s", err.Error())
				}

				for _, r := range removeList {
					fmt.Printf("Removed an unused image %s: %s\n", svc.Image, r)
				}

				currentprj = nextPrj
			}
		}
	}
}

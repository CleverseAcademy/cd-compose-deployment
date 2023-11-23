package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RemoveImage(clnt *client.Client, imageName string) ([]string, error) {
	removeResult, err := clnt.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{
		PruneChildren: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to prune an image: %w", err)
	}

	marshalledResult := make([]string, len(removeResult))
	for idx, r := range removeResult {
		if len(r.Deleted) > 0 {
			marshalledResult[idx] = fmt.Sprintf("Deleted: %s", r.Deleted)
		} else {
			marshalledResult[idx] = fmt.Sprintf("Untagged: %s", r.Untagged)
		}
	}

	return marshalledResult, nil
}

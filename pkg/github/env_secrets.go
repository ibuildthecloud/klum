package github

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/jadolg/klum/pkg/apis/klum.cattle.io/v1alpha1"
)

func createRepositoryEnvSecret(ctx context.Context, client *github.Client, syncSpec *v1alpha1.GithubSyncSpec, secretValue []byte) error {
	var key *github.PublicKey
	var repositoryID int

	repositoryID, err := getRepoID(ctx, client, syncSpec.Owner, syncSpec.Repository)
	if err != nil {
		return err
	}
	key, _, err = client.Actions.GetEnvPublicKey(ctx, repositoryID, syncSpec.Environment)
	if err != nil {
		return err
	}

	encryptedSecret, err := encodeWithPublicKey(secretValue, key.GetKey())
	if err != nil {
		return err
	}

	secret := &github.EncryptedSecret{
		Name:           syncSpec.SecretName,
		EncryptedValue: encryptedSecret,
	}

	_, err = client.Actions.CreateOrUpdateEnvSecret(ctx, repositoryID, syncSpec.Environment, secret)
	return err
}

func deleteRepositoryEnvSecret(ctx context.Context, client *github.Client, syncSpec *v1alpha1.GithubSyncSpec) error {
	repositoryID, err := getRepoID(ctx, client, syncSpec.Owner, syncSpec.Repository)
	if err != nil {
		return err
	}

	_, err = client.Actions.DeleteEnvSecret(ctx, repositoryID, syncSpec.Environment, syncSpec.SecretName)
	return err
}

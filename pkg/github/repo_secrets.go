package github

import (
	"context"

	"github.com/google/go-github/v53/github"
	"github.com/jadolg/klum/pkg/apis/klum.cattle.io/v1alpha1"
)

func createRepositorySecret(ctx context.Context, client *github.Client, syncSpec *v1alpha1.GithubSyncSpec, secretValue []byte) error {
	var key *github.PublicKey
	key, _, err := client.Actions.GetRepoPublicKey(ctx, syncSpec.Owner, syncSpec.Repository)
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

	_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, syncSpec.Owner, syncSpec.Repository, secret)
	return err
}

func deleteRepositorySecret(ctx context.Context, client *github.Client, syncSpec *v1alpha1.GithubSyncSpec) error {
	_, err := client.Actions.DeleteRepoSecret(ctx, syncSpec.Owner, syncSpec.Repository, syncSpec.SecretName)
	return err
}

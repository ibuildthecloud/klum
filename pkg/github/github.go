package github

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v53/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

func newGithubClientWithToken(token, privateURL string) (*github.Client, error) {
	var httpClient *http.Client

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient = oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(httpClient)
	if privateURL != "" {
		baseURL, err := url.Parse(privateURL)
		if err != nil {
			return nil, err
		}
		baseURL.Path = strings.TrimSuffix(baseURL.Path, "/") + "/api/v3/"
		client.BaseURL = baseURL
	}
	return client, nil
}

func getRepoID(ctx context.Context, client *github.Client, owner string, repo string) (int, error) {
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return 0, err
	}

	repositoryID := int(*repository.ID)
	return repositoryID, nil
}

func encodeWithPublicKey(text []byte, publicKey string) (string, error) {
	publicKeyDecoded, err := decodeKeyString(publicKey)
	if err != nil {
		return "", err
	}

	encrypted, err := box.SealAnonymous(nil, text, publicKeyDecoded, rand.Reader)
	if err != nil {
		return "", err
	}

	encryptedBase64 := base64.StdEncoding.EncodeToString(encrypted)

	return encryptedBase64, nil
}

func decodeKeyString(publicKey string) (*[32]byte, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	var publicKeyDecoded [32]byte
	if copy(publicKeyDecoded[:], publicKeyBytes) < 32 {
		return nil, fmt.Errorf("not a full length key, want 32 bytes, got %d bytes: %q", len(publicKeyBytes), publicKey)
	}

	return &publicKeyDecoded, nil
}

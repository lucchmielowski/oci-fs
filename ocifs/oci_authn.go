package ocifs

import (
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/github"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"io"
)

func RegistryOpts() ([]remote.Option, []name.Option, error) {
	rOpts := make([]remote.Option, 0)
	nOpts := make([]name.Option, 0)
	keychains := make([]authn.Keychain, 0)

	keychains = append(keychains,
		authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
		google.Keychain,
		authn.DefaultKeychain,
		github.Keychain,
		authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard))),
	)

	rOpts = append(rOpts, remote.WithAuthFromKeychain(authn.NewMultiKeychain(keychains...)))
	return rOpts, nOpts, nil
}

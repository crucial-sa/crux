package auth

import "github.com/zalando/go-keyring"

const (
	service = "crucial"
	user    = "default"
)

func GetSecret() (string, error) {
	secret, err := keyring.Get(service, user)
	if err != nil {
		if err == keyring.ErrNotFound {
			return "", nil
		}

		return "", err
	}

	return secret, nil
}

func SetSecret(secret string) error {
	err := keyring.Set(service, user, secret)

	return err
}

func ClearSecret() error {
	err := keyring.Delete(service, user)

	return err
}

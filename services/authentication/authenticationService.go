package authentication

import "fmt"

func IsSecretValid(token string) bool {
	repo, err := NewAuthenticationRepository()
	if err != nil {
		fmt.Printf("Can't create secret repository: %v", err)
		return false
	}

	item, err := repo.GetItem(token)
	if err != nil {
		fmt.Printf("Can't get secret: %v", err)
		return false
	}

	return item != nil
}

func SaveSecret(token string) error {
	repo, err := NewAuthenticationRepository()
	if err != nil {
		fmt.Printf("Can't create secret repository: %v", err)
		return err
	}

	item := &AuthenticationItem{
		Secret: token,
	}

	err = repo.SaveItem(item)
	if err != nil {
		fmt.Printf("Can't save secret repository: %v", err)
		return err
	}

	return nil
}

package dungeonsec

import (
	"errors"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	"net/http"
)

func SignInternalHTTPRequest(request *http.Request) error {
	if dungeon_secrets.GetDungeonDomainSecret() == "" {
		return errors.New("Domain secret is not set")
	}

	request.Header.Set(dungeon_secrets.DOMAIN_SECRET_HEADER, dungeon_secrets.GetDungeonDomainSecret())

	return nil
}

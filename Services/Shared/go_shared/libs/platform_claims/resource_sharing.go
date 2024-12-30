package platform_claims

import (
	"errors"
	"fmt"
	"libery-dungeon-libs/dungeonsec/dungeon_secrets"
	dungeon_helpers "libery-dungeon-libs/helpers"
	dungeon_models "libery-dungeon-libs/models"
	"net/http"
	"time"
)

func WriteMediaShareTokenResponse(response http.ResponseWriter, media_identity *dungeon_models.MediaIdentity, expiration time.Time) error {
	var media_share_token string

	media_share_token, err := dungeon_models.GenerateMediaShareToken(media_identity, expiration, dungeon_secrets.GetDungeonJwtSecret())
	if err != nil {
		return errors.Join(fmt.Errorf("In libery-dungeons-libs/libs/platform_claims/resource_sharing.WriteMediaShareTokenResponse: while generating media share token\n\n%s", err))
	}

	dungeon_helpers.WriteSingleStringResponse(response, media_share_token)

	return nil
}

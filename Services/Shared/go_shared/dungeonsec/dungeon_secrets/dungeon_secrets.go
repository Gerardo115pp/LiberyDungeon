package dungeon_secrets

var dungeon_jwt_secret string    // Used to sign authorization tokens that may or may not be sent to a client machine.
var dungeon_domain_secret string // Used only for internal communication
const DOMAIN_SECRET_HEADER string = "X-Dungeon-Domain-Secret"

func SetDungeonJwtSecret(secret string) {
	if dungeon_jwt_secret != "" {
		panic("Middleware JWT secret already set")
	}

	dungeon_jwt_secret = secret
}

func SetDungeonDomainSecret(secret string) {
	if dungeon_domain_secret != "" {
		panic("Middleware domain secret already set")
	}

	dungeon_domain_secret = secret
}

func GetDungeonJwtSecret() string {
	return dungeon_jwt_secret
}

func GetDungeonDomainSecret() string {
	return dungeon_domain_secret
}

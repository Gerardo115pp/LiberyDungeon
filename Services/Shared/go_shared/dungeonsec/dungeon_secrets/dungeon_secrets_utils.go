package dungeon_secrets

func CheckSecretsReady() bool {
	return dungeon_jwt_secret != "" && dungeon_domain_secret != ""
}

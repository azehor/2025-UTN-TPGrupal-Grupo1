package utils

import "regexp"

func ValidarUUID(uuid string) bool {
	var validador = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return validador.MatchString(uuid)
}

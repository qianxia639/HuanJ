package utils

import (
	"strings"
)

func DesnsitizeEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return email
	}

	localPar := email[:at]

	if len(localPar) == 3 {
		localPar = localPar[:2] + "****" + localPar[len(localPar):]
	}

	if len(localPar) > 3 {
		localPar = localPar[:3] + "****" + localPar[len(localPar):]
	}

	return localPar + email[at:]
}

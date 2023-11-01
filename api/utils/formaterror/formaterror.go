package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "email") {
		return errors.New("email Already Taken")
	}

	if strings.Contains(err, "fullname") {
		return errors.New("fullname Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect Password")
	}
	return errors.New("incorrect Details")
}

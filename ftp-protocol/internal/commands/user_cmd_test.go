package command_test

import (
	"testing"

	command "github.com/zombieleet/ftp-protocol/internal/commands"
)

func TestValidate(t *testing.T) {

	t.Run("Should return err if params is an empty array", func(tR *testing.T) {
		user := command.UserCmd{
			Params: []string{},
		}

		if err := user.Validate(); err == nil {
			t.Fail()
			return
		}

	})

	t.Run("Should pass if user is valid", func(tR *testing.T) {
		user := command.UserCmd{
			Params: []string{"test_user"},
		}

		if err := user.Validate(); err != nil {
			t.Fail()
			return
		}

	})
}

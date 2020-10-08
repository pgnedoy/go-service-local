package flags

import (
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestNewFlagSetter(t *testing.T) {
	t.Run("panics if not given a namespace", func(t *testing.T) {
		assert.Panics(t, func() { New("", &cobra.Command{}) })
	})
}

func TestFlags(t *testing.T) {
	cmd := &cobra.Command{}
	flags := New("test", cmd)

	t.Run(".RegisterInt", func(t *testing.T) {
		t.Run("panics if not given a name", func(t *testing.T) {
			assert.Panics(t, func() { flags.RegisterInt("", "", 0, "", "") })
		})

		t.Run("sets an int", func(t *testing.T) {
			defaultValue := 42
			envValue := "43"
			flagName := "meaning_of_life"
			envVariable := "MEANING_OF_LIFE"

			flags.RegisterInt(flagName, "n", defaultValue, "Is a name", envVariable)
			resp := flags.GetInt(flagName)

			assert.Equal(t, resp, defaultValue)

			os.Setenv(envVariable, envValue)

			resp = flags.GetInt(flagName)

			assert.Equal(t, resp, 43)
		})

		flags.Reset()
		os.Clearenv()
	})

	t.Run(".RegisterString", func(t *testing.T) {
		t.Run("panics if not given a name", func(t *testing.T) {
			assert.Panics(t, func() { flags.RegisterString("", "", "", "", "") })
		})

		t.Run("sets a string", func(t *testing.T) {
			defaultValue := "shoe"
			envValue := "sock"
			flagName := "product"
			envVariable := "PRODUCT"

			flags.RegisterString(flagName, "n", defaultValue, "Is a name", envVariable)
			resp := flags.GetString(flagName)

			assert.Equal(t, resp, defaultValue)

			os.Setenv(envVariable, envValue)

			resp = flags.GetString(flagName)

			assert.Equal(t, resp, envValue)
		})

		flags.Reset()
		os.Clearenv()
	})

	t.Run(".RegisterBool", func(t *testing.T) {
		t.Run("panics if not given a name", func(t *testing.T) {
			assert.Panics(t, func() { flags.RegisterBool("", "", true, "", "") })
		})

		t.Run("sets a bool", func(t *testing.T) {
			defaultValue := true
			envValue := false
			flagName := "is_integration_test"
			envVariable := "IS_INTEGRATION_TEST"

			flags.RegisterBool(flagName, "n", defaultValue, "Is an integration test", envVariable)
			resp := flags.GetBool(flagName)

			assert.Equal(t, resp, defaultValue)

			os.Setenv(envVariable, strconv.FormatBool(envValue))

			resp = flags.GetBool(flagName)

			assert.Equal(t, resp, envValue)
		})

		flags.Reset()
		os.Clearenv()
	})

	t.Run(".Reset", func(t *testing.T) {
		flagName := "title"
		flagValue := "engineer"
		flags.RegisterString(flagName, "n", flagValue, "Is a name", "")

		t.Run("sets a bool", func(t *testing.T) {
			resp := flags.GetString(flagName)
			assert.Equal(t, resp, flagValue)

			flags.Reset()

			resp = flags.GetString(flagName)
			assert.Equal(t, resp, "")
		})

		flags.Reset()
		os.Clearenv()
	})
}

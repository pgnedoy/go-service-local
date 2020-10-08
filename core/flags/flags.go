package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Flags struct {
	namespace string
	cmd       *cobra.Command
}

func (flags *Flags) formatName(name string) string {
	return fmt.Sprintf("%s.%s", flags.namespace, name)
}

func (flags *Flags) Reset() {
	flags.cmd.ResetFlags()
	viper.Reset()
}

func (flags *Flags) GetInt(flagName string) int {
	name := flags.formatName(flagName)

	return viper.GetInt(name)
}

func (flags *Flags) GetString(flagName string) string {
	name := flags.formatName(flagName)
	return viper.GetString(name)
}

func (flags *Flags) GetBool(flagName string) bool {
	name := flags.formatName(flagName)
	return viper.GetBool(name)
}

func (flags *Flags) RegisterString(name, shorthand, defaultValue, usage, env string) {
	if name == "" {
		panic("invalid name")
	}

	cmdName := flags.formatName(name)

	if shorthand != "" {
		flags.cmd.Flags().StringP(name, shorthand, defaultValue, usage)
	} else {
		flags.cmd.Flags().String(name, defaultValue, usage)
	}

	_ = viper.BindPFlag(cmdName, flags.cmd.Flags().Lookup(name))

	if env != "" {
		_ = viper.BindEnv(cmdName, env)
	}
}

func (flags *Flags) RegisterBool(name, shorthand string, defaultValue bool, usage, env string) {
	if name == "" {
		panic("invalid name")
	}

	cmdName := flags.formatName(name)

	if shorthand != "" {
		flags.cmd.Flags().BoolP(name, shorthand, defaultValue, usage)
	} else {
		flags.cmd.Flags().Bool(name, defaultValue, usage)
	}

	_ = viper.BindPFlag(cmdName, flags.cmd.Flags().Lookup(name))

	if env != "" {
		_ = viper.BindEnv(cmdName, env)
	}
}

func (flags *Flags) RegisterInt(name, shorthand string, defaultValue int, usage, env string) {
	if name == "" {
		panic("invalid name")
	}

	cmdName := flags.formatName(name)

	if shorthand != "" {
		flags.cmd.Flags().IntP(name, shorthand, defaultValue, usage)
	} else {
		flags.cmd.Flags().Int(name, defaultValue, usage)
	}

	_ = viper.BindPFlag(cmdName, flags.cmd.Flags().Lookup(name))

	if env != "" {
		_ = viper.BindEnv(cmdName, env)
	}
}

func New(namespace string, command *cobra.Command) *Flags {
	if namespace == "" {
		panic("invalid namespace")
	}

	return &Flags{namespace: namespace, cmd: command}
}

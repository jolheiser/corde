package corde

import (
	"encoding/json"
	"fmt"
)

// CreateCommander is a command that can be registered
type CreateCommander interface {
	createCommand() CreateCommand
}

// CreateOptioner is an interface for all options
type CreateOptioner interface {
	createOption() CreateOption
}

// CreateOption is the base option type for creating any sort of option
type CreateOption struct {
	Name        string           `json:"name"`
	Type        OptionType       `json:"type"`
	Description string           `json:"description,omitempty"`
	Required    bool             `json:"required,omitempty"`
	Options     []CreateOptioner `json:"options,omitempty"`
	Choices     []Choice[any]    `json:"choices,omitempty"`
}

func (c CreateOption) createOption() CreateOption {
	return c
}

// CreateCommand is a slash command that can be registered to discord
type CreateCommand struct {
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Type        CommandType      `json:"type,omitempty"`
	Options     []CreateOptioner `json:"options,omitempty"`
}

// NewSlashCommand returns a new slash command
func NewSlashCommand(name string, description string, options ...CreateOptioner) CreateCommand {
	return CreateCommand{
		Name:        name,
		Description: description,
		Options:     options,
		Type:        COMMAND_CHAT_INPUT,
	}
}

func (c CreateCommand) createCommand() CreateCommand {
	return c
}

// CommandOptionConstraint is the constraint for CommandOption types
type CommandOptionConstraint interface {
	string | int | bool | float64 // This could be enhanced for user/mention/etc types?
}

// CommandOption is an option for a CreateCommand
type CommandOption[T CommandOptionConstraint] struct {
	Name        string
	Description string
	Required    bool
	Choices     []Choice[T]
}

// NewCommandOption is a new CommandOption of type T
func NewCommandOption[T CommandOptionConstraint](name string, description string, required bool, choices ...Choice[T]) *CommandOption[T] {
	return &CommandOption[T]{
		Name:        name,
		Description: description,
		Required:    required,
		Choices:     choices,
	}
}

func (c *CommandOption[T]) createOption() CreateOption {
	var typ OptionType
	var t T
	switch fmt.Sprintf("%T", t) {
	case "string":
		typ = OPTION_STRING
	case "int":
		typ = OPTION_INTEGER
	case "bool":
		typ = OPTION_BOOLEAN
	case "float64":
		typ = OPTION_NUMBER
	}
	var choices []Choice[any]
	for _, ch := range c.Choices {
		choices = append(choices, Choice[any]{Name: ch.Name, Value: ch.Value})
	}
	return CreateOption{
		Name:        c.Name,
		Description: c.Description,
		Required:    c.Required,
		Type:        typ,
		Choices:     choices,
	}
}

func (c *CommandOption[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.createOption())
}

// SubcommandOption is an option that is a subcommand
type SubcommandOption struct {
	Name        string
	Description string
	Options     []CreateOptioner
}

// NewSubcommand returns a new subcommand
func NewSubcommand(name string, description string, options ...CreateOptioner) *SubcommandOption {
	return &SubcommandOption{
		Name:        name,
		Description: description,
		Options:     options,
	}
}

func (o *SubcommandOption) createOption() CreateOption {
	return CreateOption{
		Options:     o.Options,
		Name:        o.Name,
		Description: o.Description,
		Type:        OPTION_SUB_COMMAND,
	}
}

func (o *SubcommandOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.createOption())
}

// SubcommandGroupOption is an option that is a subcommand group
type SubcommandGroupOption struct {
	Name        string
	Description string
	Options     []CreateOptioner
}

// NewSubcommandGroup returns a new subcommand group
func NewSubcommandGroup(name string, description string, options ...CreateOptioner) *SubcommandGroupOption {
	return &SubcommandGroupOption{
		Name:        name,
		Description: description,
		Options:     options,
	}
}

func (o *SubcommandGroupOption) createOption() CreateOption {
	return CreateOption{
		Options:     o.Options,
		Name:        o.Name,
		Description: o.Description,
		Type:        OPTION_SUB_COMMAND_GROUP,
	}
}

func (o *SubcommandGroupOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.createOption())
}

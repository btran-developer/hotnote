// Package errors defines CLI exit codes for hotnote.
package errors

const (
	ExitSuccess      = iota // ExitSuccess indicates successful execution.
	ExitGeneral             // ExitGeneral indicates a general error.
	ExitNotFound            // ExitNotFound indicates the requested resource was not found.
	ExitInvalidInput        // ExitInvalidInput indicates invalid user input.
	ExitConfigError         // ExitConfigError indicates a configuration error.
)

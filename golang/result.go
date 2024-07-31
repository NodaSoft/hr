package main

import (
	"fmt"
	"strings"
)

// Result is a struct containing an operation result
// (successful or not) and some comments on this matter.
type Result struct {
	successful bool
	messages   []string
}

// NewResult returns new Result.
func NewResult(isSuccessful bool) Result {
	return Result{successful: isSuccessful}
}

// newResult is a private constructor for Result.
func newResult(isSuccessful bool, messages []string) Result {
	return Result{successful: isSuccessful, messages: messages}
}

// WithMessages returns new Result with the provided messages.
func (r Result) WithMessages(messages ...string) Result {
	return newResult(r.successful, append(r.messages, messages...))
}

// IsSuccessful returns whether the result is successful or not.
func (r Result) IsSuccessful() bool {
	return r.successful
}

// Message joins all the messages for the result.
func (r Result) Message(delimiter string) string {
	return strings.Join(r.messages, delimiter)
}

func (r Result) String() string {
	return fmt.Sprintf("Result: {successful: %t, message: %s}", r.successful, r.Message(", "))
}

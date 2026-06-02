package tools

import "errors"

var (
	ErrTooManyFindBlocks = errors.New("more than one block was found in file for change")
	ErrNoFindBlocks      = errors.New("no find blocks were found in file for change")

	ErrUserDenied = errors.New("user dined tool usage")
)

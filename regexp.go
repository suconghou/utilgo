package utilgo

import "regexp"

var (
	urlReg          = regexp.MustCompile(`^[a-zA-z]+://[^\s]+$`)
	urlStrictReg    = regexp.MustCompile(`^(?i:https?)://[[:print:]]+$`)
	ipPortReg       = regexp.MustCompile(`^([\w\-]+\.){0,5}[\w\-]+:\d{1,5}$`)
	ipPortStrictReg = regexp.MustCompile(`^([\w\-:]+\.){3,6}[\w\-:]+:\d{1,5}$`)
)

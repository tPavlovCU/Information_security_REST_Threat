package domain

import (
	"errors"
)

var ErrThreatNotFound = errors.New("threat not found")
var ErrThreatAlreadyExists = errors.New("Threat with this MD5 already exists")

var ErrBadMD5 = errors.New("BadMD5")
var ErrEmptyName = errors.New("Empty Name")
var ErrBadType = errors.New("Type must be trojan, ransomware, worm, virus, spyware, backdoor, rootkit, adware or other")
var ErrBadSeverity = errors.New("Severity must be low, medium, high or critical")
var ErrBadConfidence = errors.New("Confidence must be between 0 and 100")

package service

import (
	"errors"
	"regexp"
	"strings"

	"threat-intelligence/internal/domain"
)

var md5Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

func validateThreat(threat *domain.Threat) error {

	if threat == nil {
		return errors.New("threat cannot be nil")
	}

	//проверка md5
	md5 := threat.MD5
	ok := md5Regex.MatchString(md5)
	if !ok {
		return domain.ErrBadMD5
	}

	//проверка Name
	name := threat.Name
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrEmptyName
	}

	//проверка Type
	threatType := threat.Type

	switch threatType {
	case "trojan", "ransomware", "worm", "virus", "spyware", "backdoor", "rootkit", "adware", "other":
		//ничеог не делаем
	default:
		return domain.ErrBadType
	}

	//проверка Severity
	severity := threat.Severity

	switch severity {
	case "low", "medium", "high", "critical":
		//ничего не делаем
	default:
		return domain.ErrBadSeverity
	}

	//проверка, confidence
	confidence := threat.Confidence
	if confidence != nil {
		if *confidence < 0 || *confidence > 100 {
			return domain.ErrBadConfidence
		}
	}
	return nil
}

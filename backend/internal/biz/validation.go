package biz

import (
	"regexp"
	"strings"
)

// Validation rules and limits. Keep in sync with the frontend schemas in frontend/src/validation.ts
// and the column sizes in db/migrations.
const (
	maxEmailLen        = 254
	maxAddressLabelLen = 30
	maxNameLen         = 100
	maxAddressLineLen  = 200
	maxCityOrStateLen  = 100
)

var (
	emailRe   = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9](?:[A-Za-z0-9\-]*[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9\-]*[A-Za-z0-9])?)*\.[A-Za-z]{2,}$`)
	codeRe    = regexp.MustCompile(`^\d{6}$`)
	mobileRe  = regexp.MustCompile(`^[6-9][0-9]{9}$`) // Indian mobile numbers: 10 digits, starting 6–9
	pincodeRe = regexp.MustCompile(`^[1-9][0-9]{5}$`) // Indian PIN codes never start with 0
)

func normalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

func validEmail(e string) bool { return len(e) <= maxEmailLen && emailRe.MatchString(e) }

func validCode(c string) bool { return codeRe.MatchString(c) }

func validPhone(p string) bool { return mobileRe.MatchString(p) }

func validName(n string) bool { return n != "" && len(n) <= maxNameLen }

func validPincode(p string) bool { return pincodeRe.MatchString(p) }

func required(s string, max int) bool { return s != "" && len(s) <= max }

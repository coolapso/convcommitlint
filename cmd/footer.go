package cmd

import (
	"regexp"
	"strings"
)

var (
	footerRegexp = regexp.MustCompile(`(?m)^(.+?):\s+(.+)$`)

	breakingChangeTypos = []string{
		"BREAKING CHANGES", "BRAKING CHANGE", "BREAKING-CHANGE",
		"BREAKING_CHANGE", "BREAKING-CHANGES", "BREAKING_CHANGES",
		"BREAKINGCHANGE", "BREAKINGCHANGES", "BREAKING CHANAGE",
		"BREAKING CHAGNE", "BREAKING CHAGE", "BRAEKING CHANGE",
		"BREACKING CHANGE", "BREKING CHANGE", "BREAKING CHNAGE",
		"BREAKING CAHNGE", "BRAKING CHANGES", "BREEAKING CHANGE",
		"BREAKING CHANNGE", "BREAKIN CHANGE", "BREAKING CHANEG",
		"BARKING CHANGE", "BARKING CHAGES",
	}
)

func lintFooter(f []string) (errs []error) {
	for i, line := range f {
		if !footerRegexp.MatchString(line) {
			if (line == "" || line == "\n") && i == len(f)-1 {
				continue
			}

			if strings.HasPrefix(line, " ") {
				continue
			}

			if line == "" {
				errs = append(errs, errInvalidFooterEmptyLine)
			}

			if line == "\n" {
				errs = append(errs, errInvalidFooterEmptyLine)
			}

			errs = append(errs, errInvalidFooterFormat)
			continue
		}

		s := strings.Split(line, ":")
		if len(s) > 0 {
			if s[0] != "BREAKING CHANGE" {
				for _, typo := range breakingChangeTypos {
					if s[0] == typo || s[0] == strings.ToLower(typo) {
						errs = append(errs, errInvalidFooterBreakingTypo)
					}
				}

				if s[0] == strings.ToUpper(s[0]) {
					errs = append(errs, errInvalidFooterUpperCase)
				}

				if strings.Contains(s[0], " ") {
					errs = append(errs, errInvalidFooterSpace)
				}
			}

			if s[1] == strings.ToUpper(s[1]) {
				errs = append(errs, errInvalidFooterUpperCase)
			}
		}
	}

	return errs
}

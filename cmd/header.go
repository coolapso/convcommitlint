package cmd

import (
	"regexp"
	"slices"
)

type header struct {
	text            string
	commitType      string
	scope           string
	scopeParentesis bool
	bang            bool
	separator       bool
	description     string
}

var (
	headerRegexp = regexp.MustCompile(`^(?P<type>\w+)(\((?P<scope>[^\)]+)\))?(?P<breaking_marker>!)?(?P<colon>:) (?P<description>.+)$`)

	featTypos = []string{
		"feta", "feat:", "feature", "features", "featur",
		"fet", "feart", "faat",
		"faet", "fear", "fat", "feat.",
	}

	fixTypos = []string{
		"fxi", "feex", "fex", "fux", "fox",
		"fiz", "fiks", "fixx", "fic", "fxix",
		"ffx", "ffix", "fiix", "fx",
		"fkx", "flx", "fjx", "foix", "fidx",
		"fixx", "fis", "tix", "dix", "vix",
		"rix", "gix", "rfix", "ffix", "cix",
	}
)

func parseHeader(h string) (header header) {
	matches := headerRegexp.FindStringSubmatch(h)
	if len(matches) > 0 {
		header.text = matches[0]
		header.commitType = matches[1]

		if matches[2] != "" {
			header.scopeParentesis = true
		}

		header.scope = matches[3]

		if matches[4] == "!" {
			header.bang = true
		}

		if matches[5] == ":" {
			header.separator = true
		}
		header.description = matches[6]
	}

	return header
}

func lintHeader(h header) (errs []error) {
	if !headerRegexp.MatchString(h.text) {
		errs = append(errs, errInvalidHeaderFormat)
	}

	if len(h.description) > 72 {
		errs = append(errs, errLongHeader)
	}

	if slices.Contains(featTypos, h.commitType) {
		errs = append(errs, errTypoHeaderFeat)
	}

	if slices.Contains(fixTypos, h.commitType) {
		errs = append(errs, errTypoHeaderFix)
	}

	return errs
}

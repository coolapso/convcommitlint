package cmd

import "errors"

var (
	errInvalidBodyBlankLines     = errors.New("invalid body: one blank line should separate body from header and footer")
	errInvalidFooterFormat       = errors.New("invalid footer format: footer should be composed of key: value pairs")
	errInvalidFooterSpace        = errors.New("invalid footer format: Spaces not allowed on footer keys")
	errInvalidFooterBreakingTypo = errors.New("invalid footer format: possible typo did you mean BREAKING CHANGES?")
	errInvalidFooterUpperCase    = errors.New("invalid footer format: Uppercase only is allowed for BREAKING CHANGE")
	errInvalidFooterEmptyLine    = errors.New("invalid footer format: empty line allowed only at end of footer")
	errTypoHeaderFix             = errors.New("invalid header format: Posssible typo, did you mean fix?")
	errTypoHeaderFeat            = errors.New("invalid header format: Posssible typo, did you mean feat?")
	errLongHeader                = errors.New("description should be <= 72 characters")
	errInvalidHeaderFormat       = errors.New("invalid header format. Expected: type(optional scope): description")
	errStop                      = errors.New("stop iterations")
	errMissingPRNum              = errors.New("missing PR number for review. Set CONVOMLINT_PR_NUMBER or use --pr-number flag when not running on GitHub Actions")
	errMissingGHToken            = errors.New("github token not found, please make sure token is available to create pull request reviews")
	errMissingRepository         = errors.New("missing repository information. Set CONVOMLINT_REPOSITORY or use --repository flag when not running on GitHub Actions")
)

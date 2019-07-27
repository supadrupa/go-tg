package cmd

import (
	"os"

	"github.com/urfave/cli"

	"github.com/mr-linch/go-tg"
)

var updateTypes = (func() []string {
	result := make([]string, len(tg.UpdateTypes))

	for i, ut := range tg.UpdateTypes {
		result[i] = ut.String()
	}

	return result
})()

func isValidUpdateType(v string) bool {
	for _, ut := range updateTypes {
		if ut == v {
			return true
		}
	}

	return false
}

func isFileAvailable(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else {
		return err
	}
}

func validate(validators ...func() error) error {
	errs := make([]error, 0, len(validators))

	for _, validate := range validators {
		if err := validate(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return cli.NewMultiError(errs...)
	}
	return nil
}

func parseAllowedUpdates(vs []string) ([]tg.UpdateType, error) {
	uts := make([]tg.UpdateType, len(vs))
	for i, v := range vs {
		ut, err := tg.ParseUpdateType(v)
		if err != nil {
			return nil, err
		}
		uts[i] = ut
	}
	return uts, nil
}

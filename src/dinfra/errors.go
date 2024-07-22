package dinfra

import "fmt"

func IsError(bad bool, err error) error {
	if bad {
		return err
	}
	return nil
}

func IsErrorf(bad bool, format string, a ...any) error {
	if bad {
		return fmt.Errorf(format, a...)
	}
	return nil
}

func HasErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

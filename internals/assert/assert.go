package assert

import "fmt"

func Equalf(expected, actual any, format string, args ...any) error {
	if expected != actual {
		return fmt.Errorf("%s: expected %q but got %q", fmt.Sprintf(format, args...), expected, actual)
	}
	return nil
}

func Equal(expected, actual any, args ...any) error {
	if expected != actual {
		return fmt.Errorf("%s: expected %q but got %q", fmt.Sprint(args...), expected, actual)
	}
	return nil
}

func AssertF(err error, format string, args ...any) error {
	if err != nil {
		return fmt.Errorf("%s: %s", fmt.Sprintf(format, args...), err)
	}
	return nil
}

func Assert(err error, args ...any) error {
	if err != nil {
		return fmt.Errorf("%s: %s", fmt.Sprint(args...), err)
	}
	return nil
}

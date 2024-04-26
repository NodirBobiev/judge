package errorutil

import (
	"fmt"
	"strings"
)

func Must(err error, args ...any) {
	if err == nil {
		return
	}
	builder := &strings.Builder{}
	for i, a := range args {
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(fmt.Sprint(a))

	}
	panic(fmt.Sprintf("%s: %s", builder.String(), err))
}

func Mustf(err error, format string, args ...any) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), err))
	}
}

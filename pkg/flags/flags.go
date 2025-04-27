package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func isFlagSet(name string) bool {
	f := flag.Lookup(name)
	if f == nil {
		return false
	}
	return f.Value.String() != f.DefValue
}

func CheckRequired() {
	required := []string{"url"}
	missing := []string{}
	for _, name := range required {
		if !isFlagSet(name) {
			missing = append(missing, "-"+name)
		}
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "Missing required flags: %s\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(2)
	}
}

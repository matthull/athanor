package cli

import "strings"

// splitArgs separates positional arguments from flag arguments, allowing
// flags to appear anywhere in the argument list (not just before positionals).
// Go's flag.Parse stops at the first non-flag arg, but our CLI uses
// "command <positional> [--flag value]" patterns throughout.
//
// Returns positional args and flag args separately. The flag args can be
// passed to flag.FlagSet.Parse().
func splitArgs(args []string) (positional, flags []string) {
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			flags = append(flags, args[i])
			// If this flag takes a value (no = in it), grab the next arg too
			if !strings.Contains(args[i], "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags = append(flags, args[i+1])
				i++
			}
		} else {
			positional = append(positional, args[i])
		}
	}
	return
}

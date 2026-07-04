package postgres
import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

// StructToNamedArgs maps an arbitrary struct to pgx.NamedArgs.
// It uses the `db` struct tag for column naming and supports the `omitempty` option.
//
// The struct argument must be a struct (not a pointer to a struct).
//
// Example usage:
//
//	type User struct {
//	    ID    int64  `db:"id,omitempty"`
//	    Name  string `db:"user_name"`
//	    Email string `db:"email,omitempty"`
//	}
//
// user := User{Name: "Alice"}
// args, err := StructToNamedArgs(user) // args will be pgx.NamedArgs{"user_name": "Alice"}
func StructToNamedArgs(s any) (pgx.NamedArgs, error) {
	v := reflect.ValueOf(s)

	// Ensure the input is a struct
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, fmt.Errorf("StructToNamedArgs received a nil pointer")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToNamedArgs expected a struct, got %s", v.Kind())
	}

	t := v.Type()
	args := make(pgx.NamedArgs)

	// Iterate over all fields of the struct
	for i := 0; i < t.NumField(); i++ {
		fieldT := t.Field(i)
		fieldV := v.Field(i)

		// Only process exported fields (names start with an uppercase letter)
		if !fieldV.CanInterface() {
			continue
		}

		// 1. Get the 'db' tag value
		tag, ok := fieldT.Tag.Lookup("db")
		if !ok {
			// Skip fields without a 'db' tag
			continue
		}

		// 2. Parse the tag for name and options (e.g., "column_name,omitempty")
		parts := strings.Split(tag, ",")
		columnName := parts[0]

		// Ignore fields explicitly tagged with `db:"-"`
		if columnName == "-" {
			continue
		}

		// Check for omitempty option
		isOmitEmpty := false
		if len(parts) > 1 && strings.Contains(parts[1], "omitempty") {
			isOmitEmpty = true
		}

		// 3. Check for zero value if omitempty is set
		if isOmitEmpty {
			// reflect.Value.IsZero() is the canonical way to check for a zero value in Go 1.13+
			if fieldV.IsZero() {
				continue // Skip the field if it's zero value and omitempty is present
			}
		}

		// 4. Add to NamedArgs
		args[columnName] = fieldV.Interface()
	}

	return args, nil
}
package cli

import (
	"flag"
	"fmt"
	"time"
)

type TimestampFlag = FlagBase[time.Time, TimestampConfig, timestampValue]

// TimestampConfig defines the config for timestamp flags
type TimestampConfig struct {
	Timezone *time.Location
	Layout   string
}

// timestampValue wrap to satisfy golang's flag interface.
type timestampValue struct {
	timestamp  *time.Time
	hasBeenSet bool
	layout     string
	location   *time.Location
}

// Below functions are to satisfy the ValueCreator interface

func (i timestampValue) Create(val time.Time, p *time.Time, c TimestampConfig) Value {
	*p = val
	return &timestampValue{
		timestamp: p,
		layout:    c.Layout,
		location:  c.Timezone,
	}
}

func (i timestampValue) ToString(b time.Time) string {
	if b.IsZero() {
		return ""
	}
	return fmt.Sprintf("%v", b)
}

// Timestamp constructor(for internal testing only)
func newTimestamp(timestamp time.Time) *timestampValue {
	return &timestampValue{timestamp: &timestamp}
}

// Below functions are to satisfy the flag.Value interface

// Parses the string value to timestamp
func (t *timestampValue) Set(value string) error {
	var timestamp time.Time
	var err error

	if t.location != nil {
		timestamp, err = time.ParseInLocation(t.layout, value, t.location)
	} else {
		timestamp, err = time.Parse(t.layout, value)
	}

	if err != nil {
		return err
	}

	if t.timestamp != nil {
		*t.timestamp = timestamp
	}
	t.hasBeenSet = true
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (t *timestampValue) String() string {
	return fmt.Sprintf("%#v", t.timestamp)
}

// Value returns the timestamp value stored in the flag
func (t *timestampValue) Value() *time.Time {
	return t.timestamp
}

// Get returns the flag structure
func (t *timestampValue) Get() any {
	return *t.timestamp
}

// Timestamp gets the timestamp from a flag name
func (cmd *Command) Timestamp(name string) *time.Time {
	if fs := cmd.lookupFlagSet(name); fs != nil {
		return lookupTimestamp(name, fs, cmd.Name)
	}
	return nil
}

// Fetches the timestamp value from the local timestampWrap
func lookupTimestamp(name string, set *flag.FlagSet, cmdName string) *time.Time {
	fl := set.Lookup(name)
	if fl != nil {
		if tv, ok := fl.Value.(*timestampValue); ok {
			v := tv.Value()

			tracef("timestamp available for flag name %[1]q with value=%[2]v (cmd=%[3]q)", name, v, cmdName)
			return v
		}
	}

	tracef("timestamp NOT available for flag name %[1]q (cmd=%[2]q)", name, cmdName)
	return nil
}

package randomizer

import "errors"

// ErrInvalidVersion is returned when a Version string does not match the
// expected format for a Version.
var ErrInvalidVersion = errors.New("invalid version")

// ErrNoDefaultPreset is returned when a Metadata file fails validation
// because there is no preset with the 'default' id.
var ErrNoDefaultPreset = errors.New("no default preset")

// ErrNoPresets is returned when a Metadata file fails validation because
// the presets array is empty.
var ErrNoPresets = errors.New("presets array is empty")

var ErrNoPatchFile = errors.New("patch file was not generated")

var ErrNoSpoilerLog = errors.New("spoiler log was not generated")

package transpiler

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

// TranspileToTOML converts the transpiled data to TOML format.
func (t *Transpiler) TranspileToTOML() ([]byte, error) {
	root, err := t.buildRootMap()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(root); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

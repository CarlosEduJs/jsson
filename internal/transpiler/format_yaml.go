package transpiler

import (
	yamlv3 "gopkg.in/yaml.v3"
)

// TranspileToYAML converts the transpiled data to YAML format.
func (t *Transpiler) TranspileToYAML() ([]byte, error) {
	root, err := t.buildRootMap()
	if err != nil {
		return nil, err
	}

	return yamlv3.Marshal(root)
}

package transpiler

import (
	"context"

	yamlv3 "gopkg.in/yaml.v3"
)

// TranspileToYAML converts the transpiled data to YAML format.
func (t *Transpiler) TranspileToYAML(ctx context.Context) ([]byte, error) {
	root, err := t.buildRootMap(ctx)
	if err != nil {
		return nil, err
	}

	return yamlv3.Marshal(root)
}

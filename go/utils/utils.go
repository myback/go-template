package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
)

func CheckErr(err error, tpls ...string) {
	if err != nil {
		tpl := "%s"
		if len(tpls) > 0 {
			tpl = tpls[0]
		}

		_, _ = fmt.Fprintf(os.Stderr, tpl, err)
		os.Exit(1)
	}
}

func MustUnmarshalValues(data string) map[string]any {
	values := map[string]any{}
	CheckErr(json.Unmarshal([]byte(data), &values))

	return values
}

func ReadValuesFile(valuesFile string) (map[string]any, error) {
	fi, err := os.Open(valuesFile)
	if err != nil {
		return nil, err
	}
	defer CheckErr(fi.Close())

	values := map[string]any{}
	ext := filepath.Ext(filepath.Base(valuesFile))

	switch ext {
	case ".yml", ".yaml":
		b, err := io.ReadAll(fi)
		if err != nil {
			return nil, err
		}

		if err = yaml.Unmarshal(b, &values); err != nil {
			return nil, err
		}
	case ".json":
		if err = json.NewDecoder(fi).Decode(&values); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown file type: %s", ext)
	}

	return values, nil
}

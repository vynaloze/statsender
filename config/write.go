package config

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"os"
	"path/filepath"
)

func AddDatasource(configDir string, filename string, datasource Datasource) error {
	path := filepath.Join(configDir, filename)
	block := struct {
		Datasource Datasource `hcl:"datasource,block"`
	}{datasource}

	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(block, f.Body())

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := file.Write(f.Bytes()); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

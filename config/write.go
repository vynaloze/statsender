package config

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/pkg/errors"
	"github.com/vynaloze/statsender/sender"
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

func AddSender(configDir string, filename string, s sender.Sender) error {
	path := filepath.Join(configDir, filename)
	f := hclwrite.NewEmptyFile()

	switch v := s.(type) {
	case sender.Sout:
		block := struct {
			Sender sender.Sout `hcl:"console,block"`
		}{v}

		gohcl.EncodeIntoBody(block, f.Body())
	case sender.Http:
		block := struct {
			Sender sender.Http `hcl:"http,block"`
		}{v}
		gohcl.EncodeIntoBody(block, f.Body())
	default:
		return errors.New("invalid sender type - valid types: 'sout', 'http'")
	}

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

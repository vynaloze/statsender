package config

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/pkg/errors"
	"github.com/vynaloze/statsender/logger"
	"github.com/vynaloze/statsender/sender"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func AddDatasource(configDir string, filename string, datasource Datasource) error {
	path := filepath.Join(configDir, filename)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
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
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
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
		return errors.New("invalid sender type - valid types: 'console', 'http'")
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

type collectorConfig struct {
	System   *System   `hcl:"system,block"`
	Postgres *Postgres `hcl:"postgres,block"`
}

func SetCollectorEnabled(configDir string, filename string, typ string, enabled bool) error {
	return replaceInFile(configDir, filename, typ, func(field *reflect.Value) {
		field.FieldByName("Enabled").SetBool(enabled)
	})
}

func SetCollectorCron(configDir string, filename string, typ string, cron string) error {
	return replaceInFile(configDir, filename, typ, func(field *reflect.Value) {
		field.FieldByName("Cron").SetString(cron)
	})
}

func replaceInFile(configDir string, filename string, typ string, setValFunc func(field *reflect.Value)) error {
	// Read file
	path := filepath.Join(configDir, filename)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return errors.New("Error parsing file. Maybe file " + path + " does not exist?")
	}
	var c collectorConfig
	moreDiags := gohcl.DecodeBody(f.Body, nil, &c)
	diags = append(diags, moreDiags...)

	// Check for errors
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,      // writer to send messages to
			parser.Files(), // the parser's file cache, for source snippets
			78,             // wrapping width
			true,           // generate colored/highlighted output
		)
		err := wr.WriteDiagnostics(diags)
		if err != nil {
			return err
		}
		return errors.New("invalid configuration - see diagnostics above")
	}

	// Change state of required collector
	c, err = setValue(c, typ, setValFunc)
	if err != nil {
		return err
	}
	// encode
	w := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(c, w.Body())
	// write
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := file.Write(w.Bytes()); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func setValue(c collectorConfig, typ string, setValFunc func(field *reflect.Value)) (collectorConfig, error) {
	log, _ := logger.New()
	config := reflect.ValueOf(c)
	for i := 0; i < config.NumField(); i++ {
		v := config.Field(i).Elem() // system or postgres
		log.Debugf("checking %+v", v)
		for j := 0; j < v.NumField(); j++ {
			tag := v.Type().Field(j).Tag.Get("hcl") // cpu,... or pg_stat...
			name := strings.Split(tag, ",")[0]
			log.Debugf("checking %s", name)
			if name == typ {
				log.Debugf("updating %s", name)
				field := v.Field(j).Elem()
				log.Debugf("field before: %+v", field)
				setValFunc(&field)
				log.Debugf("field after: %+v", field)
				return c, nil
			}
		}
	}
	return c, errors.New("collector " + typ + " not found in config")
}

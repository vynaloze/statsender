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

func writeToFile(path string, content []byte, flags int) error {
	file, err := os.OpenFile(path, flags, 0755)
	if err != nil {
		return err
	}
	if _, err := file.Write(content); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

func AddDatasource(configDir string, filename string, datasource Datasource) error {
	log, _ := logger.New()
	path := filepath.Join(configDir, filename)
	log.Debug("accessing file " + path)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
	block := struct {
		Datasource Datasource `hcl:"datasource,block"`
	}{datasource}

	log.Debugf("encoding struct %+v", block)
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(block, f.Body())

	log.Debug("appending to file " + path)
	return writeToFile(path, f.Bytes(), os.O_APPEND|os.O_CREATE|os.O_WRONLY)
}

func AddSender(configDir string, filename string, s sender.Sender) error {
	log, _ := logger.New()
	path := filepath.Join(configDir, filename)
	log.Debug("accessing file " + path)
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
		log.Debugf("encoding struct %+v", block)
		gohcl.EncodeIntoBody(block, f.Body())
	case sender.Http:
		block := struct {
			Sender sender.Http `hcl:"http,block"`
		}{v}
		log.Debugf("encoding struct %+v", block)
		gohcl.EncodeIntoBody(block, f.Body())
	default:
		return errors.New("invalid sender type - valid types: 'console', 'http'")
	}

	log.Debug("appending to file " + path)
	return writeToFile(path, f.Bytes(), os.O_APPEND|os.O_CREATE|os.O_WRONLY)
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
	log, _ := logger.New()
	// Read file
	path := filepath.Join(configDir, filename)
	log.Debug("accessing file " + path)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
	log.Debug("parsing file")
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return errors.New("Error parsing file. Maybe file " + path + " does not exist?")
	}
	log.Debug("decoding file")
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
	log.Debug("changing state of collector " + typ)
	c, err = setValue(c, typ, setValFunc)
	if err != nil {
		return err
	}
	// encode
	log.Debug("encoding file back")
	w := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(c, w.Body())
	// write
	log.Debug("replacing old file")
	return writeToFile(path, w.Bytes(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
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

func InitConfig(configDir string) error {
	log, _ := logger.New()
	// create dir if not exists
	log.Debugf("creating config dir %s (if it does not exist)", configDir)
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}
	// create collector config example
	log.Debug("creating _collectors.hcl")
	pathC := filepath.Join(configDir, "_collectors.hcl")
	err = writeToFile(pathC, []byte(defaultCollectorConfig), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	// create datasource config example
	log.Debug("creating _datasources.hcl")
	pathD := filepath.Join(configDir, "_datasources.hcl")
	err = writeToFile(pathD, []byte(defaultDatasourceConfig), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	// create sender config example
	log.Debug("creating _senders.hcl")
	pathS := filepath.Join(configDir, "_senders.hcl")
	err = writeToFile(pathS, []byte(defaultSenderConfig), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	log.Debug("done...")
	return nil
}

package config

import (
	"github.com/fatih/color"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// ConfigError represents all parsing errors, which occurred while reading the configuration file
type ConfigError struct {
	err   error
	diags hcl.Diagnostics
	files map[string]*hcl.File
}

// ReadConfig finds, parses and merges all .hcl and .json files in a given directory,
// then it decodes the data into Config struct or returns ConfigError
func ReadConfig(dir string) (*Config, *ConfigError) {
	// Find, parse, merge and decode all files
	parser := hclparse.NewParser()
	var files []*hcl.File
	var diags hcl.Diagnostics

	// find all files
	allFiles, err := allHclFiles(dir)
	if err != nil {
		return nil, &ConfigError{err: err}
	}

	// parse them
	for _, fn := range allFiles {
		f, moreDiags := parser.ParseHCLFile(fn)
		files = append(files, f)
		diags = append(diags, moreDiags...)
	}
	// Check for errors
	if diags.HasErrors() {
		return nil, &ConfigError{diags: diags, files: parser.Files()}
	}

	// merge them
	body := hcl.MergeFiles(files)
	// decode them
	var c Config
	moreDiags := gohcl.DecodeBody(body, nil, &c)
	diags = append(diags, moreDiags...)

	// Check for errors
	if diags.HasErrors() {
		return nil, &ConfigError{diags: diags, files: parser.Files()}
	}
	return &c, nil
}

// Print pretty-prints all the parsing errors
func (ce ConfigError) Print() error {
	if ce.err != nil {
		redBold := color.New(color.FgRed).Add(color.Bold)
		red := color.New(color.FgRed)
		_, _ = redBold.Print("FAIL: ")
		_, _ = red.Println(ce.err.Error())
	}
	if ce.diags != nil {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout, // writer to send messages to
			ce.files,  // the parser's file cache, for source snippets
			78,        // wrapping width
			true,      // generate colored/highlighted output
		)
		return wr.WriteDiagnostics(ce.diags)
	}
	return nil
}

func allHclFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".hcl" || filepath.Ext(path) == ".json" {
				files = append(files, filepath.Join(dir, f.Name()))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("no files found in the given directory " + dir)
	}
	return files, nil
}

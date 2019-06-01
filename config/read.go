package config

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func ReadConfig(dir string) (*Config, error) {
	// Find, parse, merge and decode all files
	parser := hclparse.NewParser()
	var files []*hcl.File
	var diags hcl.Diagnostics

	allFiles, err := allHclFiles(dir)
	if err != nil {
		return nil, err
	}

	for _, fn := range allFiles {
		f, moreDiags := parser.ParseHCLFile(fn)
		files = append(files, f)
		diags = append(diags, moreDiags...)
	}

	body := hcl.MergeFiles(files)
	var c Config
	moreDiags := gohcl.DecodeBody(body, nil, &c)
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
			return nil, err
		}
		return nil, errors.New("invalid configuration - see diagnostics above")
	}
	return &c, nil
}

func allHclFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".hcl" {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, errors.New("No files found in the given directory " + dir)
	}
	return files, nil
}

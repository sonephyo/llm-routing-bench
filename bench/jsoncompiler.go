package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// <load_pattern>_<token_size>_<prompt_type>.json
func WriteResult(result ExperimentResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.Chmod(outputDir, 0777); err != nil {
		return fmt.Errorf("chmod output dir: %w", err)
	}

	filename := fmt.Sprintf("%s_%d_%s.json",
		result.Metadata.LoadPattern,
		result.Metadata.TokenSize,
		result.Metadata.PromptType,
	)
	path := filepath.Join(outputDir, filename)

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal result: %w", err)
	}

	if err := os.WriteFile(path, data, 0666); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	if err := os.Chmod(path, 0666); err != nil {
		return fmt.Errorf("chmod result file: %w", err)
	}
	return nil
}

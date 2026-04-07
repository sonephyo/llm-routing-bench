package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// <load_pattern>_<token_size>_<prompt_type>.json
func WriteResult(result ExperimentResult, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
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

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

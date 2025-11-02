package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func BaseConfig(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)

	cacheFile := "base_config.cache.yaml"

	if err != nil {
		if cached, cerr := readCache(cacheFile); cerr == nil {
			if validateYAML([]byte(cached)) == nil {
				return cached, nil
			}
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if cached, cerr := readCache(cacheFile); cerr == nil {
			if validateYAML([]byte(cached)) == nil {
				return cached, nil
			}
		}
		return "", &httpError{status: resp.Status}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		if cached, cerr := readCache(cacheFile); cerr == nil {
			if validateYAML([]byte(cached)) == nil {
				return cached, nil
			}
		}
		return "", err
	}
	// Validate remote YAML; if invalid, try cached valid YAML
	if vErr := validateYAML(data); vErr != nil {
		if cached, cerr := readCache(cacheFile); cerr == nil {
			if validateYAML([]byte(cached)) == nil {
				return cached, nil
			}
		}
		return "", fmt.Errorf("invalid YAML from remote: %v", vErr)
	}

	// Write latest valid content to cache (best-effort)
	_ = writeCache(cacheFile, data)

	return string(data), nil
}

type httpError struct{ status string }

func (e *httpError) Error() string { return e.status }

func readCache(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func writeCache(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// validateYAML checks if the provided data is syntactically valid YAML.
func validateYAML(data []byte) error {
	var n yaml.Node
	return yaml.Unmarshal(data, &n)
}

package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"zeo-api/internal/config"
	"zeo-api/internal/utils/file"
)

type ZeoRunner struct {
	config *config.ZeoConfig
}

type ZeoResult struct {
	Success     bool
	ExitCode    int
	Stdout      string
	Stderr      string
	OutputFiles map[string][]byte
}

func NewZeoRunner(cfg *config.ZeoConfig) *ZeoRunner {
	return &ZeoRunner{config: cfg}
}

func (zr *ZeoRunner) RunCommand(ctx context.Context, structureFile string, args []string, outputFiles []string) (*ZeoResult, error) {
	// Ensure workspace exists
	if err := os.MkdirAll(zr.config.Workdir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Copy structure file to workspace
	workspaceFile := filepath.Join(zr.config.Workdir, filepath.Base(structureFile))
	if err := zr.copyFile(structureFile, workspaceFile); err != nil {
		return nil, fmt.Errorf("failed to copy structure file: %w", err)
	}
	defer func() {
		_ = os.Remove(workspaceFile)
	}()

	// Prepare command arguments
	fullArgs := append([]string{}, args...)
	fullArgs = append(fullArgs, workspaceFile)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, zr.config.Timeout)
	defer cancel()

	// Execute command
	cmd := exec.CommandContext(ctx, zr.config.ExecutablePath, fullArgs...)
	cmd.Dir = zr.config.Workdir

	stdout, err := cmd.CombinedOutput()

	result := &ZeoResult{
		Success:     err == nil,
		Stdout:      string(stdout),
		Stderr:      "",
		OutputFiles: make(map[string][]byte),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		result.Stderr = err.Error()
	}

	// Collect output files
	maxFileSize := int64(100 * 1024 * 1024) // 100MB limit
	for _, outputFile := range outputFiles {
		outputPath := filepath.Join(zr.config.Workdir, outputFile)

		// Ensure output file is within workspace
		absOutputPath, err := filepath.Abs(filepath.Clean(outputPath))
		absWorkdir, err2 := filepath.Abs(filepath.Clean(zr.config.Workdir))
		if err != nil || err2 != nil || !strings.HasPrefix(absOutputPath, absWorkdir) {
			continue
		}

		if file.FileExists(outputPath) {
			// Check file size
			info, err := os.Stat(outputPath)
			if err != nil || info.Size() > maxFileSize {
				_ = os.Remove(outputPath)
				continue
			}

			content, err := file.GetFileContent(outputPath)
			if err != nil {
				continue // Skip files that can't be read
			}
			result.OutputFiles[outputFile] = content
			_ = os.Remove(outputPath) // Clean up after reading
		}
	}

	return result, nil
}

func (zr *ZeoRunner) copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func (zr *ZeoRunner) ValidateZeoExecutable() error {
	_, err := exec.LookPath(zr.config.ExecutablePath)
	return err
}

func BuildZeoArgs(analysisType string, params map[string]interface{}) ([]string, error) {
	var args []string

	// Add high accuracy flag if specified
	if ha, ok := params["ha"].(bool); ok && ha {
		args = append(args, "-ha")
	}

	switch analysisType {
	case "pore_diameter":
		args = append(args, "-res", "output.res")
	case "surface_area":
		probeRadius := getFloatParam(params, "probe_radius", 1.21)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		samples := getIntParam(params, "samples", 2000)
		if err := validateIntParam(samples, 100, 1000000, "samples"); err != nil {
			return nil, err
		}
		args = append(args, "-sa", fmt.Sprintf("%.2f", probeRadius), fmt.Sprintf("%d", samples), "output.sa")
	case "accessible_volume":
		probeRadius := getFloatParam(params, "probe_radius", 1.21)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		chanRadius := getFloatParam(params, "chan_radius", 1.21)
		if err := validateFloatParam(chanRadius, 0.1, 10.0, "chan_radius"); err != nil {
			return nil, err
		}
		samples := getIntParam(params, "samples", 2000)
		if err := validateIntParam(samples, 100, 1000000, "samples"); err != nil {
			return nil, err
		}
		args = append(args, "-vol", fmt.Sprintf("%.2f", probeRadius), fmt.Sprintf("%.2f", chanRadius), fmt.Sprintf("%d", samples), "output.vol")
	case "probe_volume":
		probeRadius := getFloatParam(params, "probe_radius", 1.21)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		chanRadius := getFloatParam(params, "chan_radius", 1.21)
		if err := validateFloatParam(chanRadius, 0.1, 10.0, "chan_radius"); err != nil {
			return nil, err
		}
		samples := getIntParam(params, "samples", 2000)
		if err := validateIntParam(samples, 100, 1000000, "samples"); err != nil {
			return nil, err
		}
		args = append(args, "-volpo", fmt.Sprintf("%.2f", probeRadius), fmt.Sprintf("%.2f", chanRadius), fmt.Sprintf("%d", samples), "output.volpo")
	case "channel_analysis":
		probeRadius := getFloatParam(params, "probe_radius", 1.21)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		args = append(args, "-chan", fmt.Sprintf("%.2f", probeRadius), "output.chan")
	case "framework_info":
		args = append(args, "-strinfo", "output.strinfo")
	case "pore_size_dist":
		probeRadius := getFloatParam(params, "probe_radius", 1.21)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		chanRadius := getFloatParam(params, "chan_radius", probeRadius)
		if err := validateFloatParam(chanRadius, 0.1, 10.0, "chan_radius"); err != nil {
			return nil, err
		}
		samples := getIntParam(params, "samples", 50000)
		if err := validateIntParam(samples, 1000, 1000000, "samples"); err != nil {
			return nil, err
		}
		args = append(args, "-psd", fmt.Sprintf("%.2f", probeRadius), fmt.Sprintf("%.2f", chanRadius), "100", fmt.Sprintf("%d", samples), "output.psd")
	case "blocking_spheres":
		probeRadius := getFloatParam(params, "probe_radius", 1.86)
		if err := validateFloatParam(probeRadius, 0.1, 10.0, "probe_radius"); err != nil {
			return nil, err
		}
		args = append(args, "-block", fmt.Sprintf("%.2f", probeRadius), "output.block")
	case "open_metal_sites":
		args = append(args, "-oms", "output.oms")
	default:
		return nil, fmt.Errorf("unsupported analysis type: %s", analysisType)
	}

	return args, nil
}

func validateFloatParam(value float64, min, max float64, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %.2f and %.2f", name, min, max)
	}
	return nil
}

func validateIntParam(value int, min, max int, name string) error {
	if value < min || value > max {
		return fmt.Errorf("%s must be between %d and %d", name, min, max)
	}
	return nil
}

func getFloatParam(params map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := params[key].(float64); ok {
		return val
	}
	if val, ok := params[key].(float32); ok {
		return float64(val)
	}
	return defaultValue
}

func getIntParam(params map[string]interface{}, key string, defaultValue int) int {
	if val, ok := params[key].(int); ok {
		return val
	}
	if val, ok := params[key].(float64); ok {
		return int(val)
	}
	return defaultValue
}

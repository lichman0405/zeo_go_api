package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type PoreDiameterResult struct {
	IncludedDiameter  float64 `json:"included_diameter"`
	FreeDiameter      float64 `json:"free_diameter"`
	IncludedAlongFree float64 `json:"included_along_free"`
}

type SurfaceAreaResult struct {
	ASAUnitcell  float64 `json:"asa_unitcell"`
	ASAVolume    float64 `json:"asa_volume"`
	ASAMass      float64 `json:"asa_mass"`
	NASAUnitcell float64 `json:"nasa_unitcell"`
	NASAVolume   float64 `json:"nasa_volume"`
	NASAMass     float64 `json:"nasa_mass"`
}

type AccessibleVolumeResult struct {
	UnitcellVolume float64            `json:"unitcell_volume"`
	Density        float64            `json:"density"`
	AV             map[string]float64 `json:"av"`
	NAV            map[string]float64 `json:"nav"`
}

type ProbeVolumeResult struct {
	POAVUnitcell  float64 `json:"poav_unitcell"`
	POAVFraction  float64 `json:"poav_fraction"`
	POAVMass      float64 `json:"poav_mass"`
	PONAVUnitcell float64 `json:"ponav_unitcell"`
	PONAVFraction float64 `json:"ponav_fraction"`
	PONAVMass     float64 `json:"ponav_mass"`
}

type ChannelAnalysisResult struct {
	Dimension         int     `json:"dimension"`
	IncludedDiameter  float64 `json:"included_diameter"`
	FreeDiameter      float64 `json:"free_diameter"`
	IncludedAlongFree float64 `json:"included_along_free"`
}

type FrameworkInfoResult struct {
	Filename           string                   `json:"filename"`
	Formula            string                   `json:"formula"`
	Segments           int                      `json:"segments"`
	NumberOfFrameworks int                      `json:"number_of_frameworks"`
	NumberOfMolecules  int                      `json:"number_of_molecules"`
	Frameworks         []map[string]interface{} `json:"frameworks"`
}

type BlockingSpheresResult struct {
	Channels      []interface{} `json:"channels"`
	Pockets       []interface{} `json:"pockets"`
	NodesAssigned []interface{} `json:"nodes_assigned"`
	Raw           string        `json:"raw"`
}

type OpenMetalSitesResult struct {
	OpenMetalSitesCount int `json:"open_metal_sites_count"`
}

// ParsePoreDiameter parses Zeo++ -res output
func ParsePoreDiameter(data string) (*PoreDiameterResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	parts := strings.Fields(lines[len(lines)-1])
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid format: expected 3 values, got %d", len(parts))
	}

	included, err1 := strconv.ParseFloat(parts[0], 64)
	free, err2 := strconv.ParseFloat(parts[1], 64)
	includedAlong, err3 := strconv.ParseFloat(parts[2], 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("failed to parse values: %v, %v, %v", err1, err2, err3)
	}

	return &PoreDiameterResult{
		IncludedDiameter:  included,
		FreeDiameter:      free,
		IncludedAlongFree: includedAlong,
	}, nil
}

// ParseSurfaceArea parses Zeo++ -sa output
func ParseSurfaceArea(data string) (*SurfaceAreaResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	lastLine := lines[len(lines)-1]
	parts := strings.Fields(lastLine)
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid format: expected 6 values, got %d", len(parts))
	}

	asaUnitcell, _ := strconv.ParseFloat(parts[0], 64)
	asaVolume, _ := strconv.ParseFloat(parts[1], 64)
	asaMass, _ := strconv.ParseFloat(parts[2], 64)
	nasaUnitcell, _ := strconv.ParseFloat(parts[3], 64)
	nasaVolume, _ := strconv.ParseFloat(parts[4], 64)
	nasaMass, _ := strconv.ParseFloat(parts[5], 64)

	return &SurfaceAreaResult{
		ASAUnitcell:  asaUnitcell,
		ASAVolume:    asaVolume,
		ASAMass:      asaMass,
		NASAUnitcell: nasaUnitcell,
		NASAVolume:   nasaVolume,
		NASAMass:     nasaMass,
	}, nil
}

// ParseAccessibleVolume parses Zeo++ -vol output
func ParseAccessibleVolume(data string) (*AccessibleVolumeResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	lastLine := lines[len(lines)-1]
	parts := strings.Fields(lastLine)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid format: expected 3 values, got %d", len(parts))
	}

	unitcellVolume, _ := strconv.ParseFloat(parts[0], 64)
	density, _ := strconv.ParseFloat(parts[1], 64)
	av, _ := strconv.ParseFloat(parts[2], 64)

	return &AccessibleVolumeResult{
		UnitcellVolume: unitcellVolume,
		Density:        density,
		AV:             map[string]float64{"value": av},
		NAV:            map[string]float64{"value": 0}, // Default if not provided
	}, nil
}

// ParseProbeVolume parses Zeo++ -volpo output
func ParseProbeVolume(data string) (*ProbeVolumeResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	lastLine := lines[len(lines)-1]
	parts := strings.Fields(lastLine)
	if len(parts) < 6 {
		return nil, fmt.Errorf("invalid format: expected 6 values, got %d", len(parts))
	}

	poavUnitcell, _ := strconv.ParseFloat(parts[0], 64)
	poavFraction, _ := strconv.ParseFloat(parts[1], 64)
	poavMass, _ := strconv.ParseFloat(parts[2], 64)
	ponavUnitcell, _ := strconv.ParseFloat(parts[3], 64)
	ponavFraction, _ := strconv.ParseFloat(parts[4], 64)
	ponavMass, _ := strconv.ParseFloat(parts[5], 64)

	return &ProbeVolumeResult{
		POAVUnitcell:  poavUnitcell,
		POAVFraction:  poavFraction,
		POAVMass:      poavMass,
		PONAVUnitcell: ponavUnitcell,
		PONAVFraction: ponavFraction,
		PONAVMass:     ponavMass,
	}, nil
}

// ParseChannelAnalysis parses Zeo++ -chan output
func ParseChannelAnalysis(data string) (*ChannelAnalysisResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	lastLine := lines[len(lines)-1]
	parts := strings.Fields(lastLine)
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid format: expected 4 values, got %d", len(parts))
	}

	dimension, _ := strconv.Atoi(parts[0])
	includedDiameter, _ := strconv.ParseFloat(parts[1], 64)
	freeDiameter, _ := strconv.ParseFloat(parts[2], 64)
	includedAlongFree, _ := strconv.ParseFloat(parts[3], 64)

	return &ChannelAnalysisResult{
		Dimension:         dimension,
		IncludedDiameter:  includedDiameter,
		FreeDiameter:      freeDiameter,
		IncludedAlongFree: includedAlongFree,
	}, nil
}

// ParseFrameworkInfo parses Zeo++ -strinfo output
func ParseFrameworkInfo(data string) (*FrameworkInfoResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	// Simple parsing - in real implementation, this would be more sophisticated
	result := &FrameworkInfoResult{
		Filename:           "",
		Formula:            "",
		Segments:           1,
		NumberOfFrameworks: 1,
		NumberOfMolecules:  0,
		Frameworks:         []map[string]interface{}{},
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Framework") {
			result.NumberOfFrameworks++
		}
	}

	return result, nil
}

// ParseBlockingSpheres parses Zeo++ -block output
func ParseBlockingSpheres(data string) (*BlockingSpheresResult, error) {
	return &BlockingSpheresResult{
		Channels:      []interface{}{},
		Pockets:       []interface{}{},
		NodesAssigned: []interface{}{},
		Raw:           data,
	}, nil
}

// ParseOpenMetalSites parses Zeo++ -oms output
func ParseOpenMetalSites(data string) (*OpenMetalSitesResult, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty output")
	}

	lastLine := strings.TrimSpace(lines[len(lines)-1])
	count, err := strconv.ParseInt(lastLine, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse open metal sites count: %w", err)
	}

	return &OpenMetalSitesResult{
		OpenMetalSitesCount: int(count),
	}, nil
}

// ParseOutputFile parses the specified output file based on analysis type
func ParseOutputFile(analysisType string, data string) (interface{}, error) {
	switch analysisType {
	case "pore_diameter":
		return ParsePoreDiameter(data)
	case "surface_area":
		return ParseSurfaceArea(data)
	case "accessible_volume":
		return ParseAccessibleVolume(data)
	case "probe_volume":
		return ParseProbeVolume(data)
	case "channel_analysis":
		return ParseChannelAnalysis(data)
	case "framework_info":
		return ParseFrameworkInfo(data)
	case "blocking_spheres":
		return ParseBlockingSpheres(data)
	case "open_metal_sites":
		return ParseOpenMetalSites(data)
	default:
		return nil, fmt.Errorf("unsupported analysis type: %s", analysisType)
	}
}

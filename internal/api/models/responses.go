package models

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Cached  bool        `json:"cached"`
}

type PoreDiameterResponse struct {
	IncludedDiameter  float64 `json:"included_diameter"`
	FreeDiameter      float64 `json:"free_diameter"`
	IncludedAlongFree float64 `json:"included_along_free"`
	Cached            bool    `json:"cached"`
}

type SurfaceAreaResponse struct {
	ASAUnitcell  float64 `json:"asa_unitcell"`
	ASAVolume    float64 `json:"asa_volume"`
	ASAMass      float64 `json:"asa_mass"`
	NASAUnitcell float64 `json:"nasa_unitcell"`
	NASAVolume   float64 `json:"nasa_volume"`
	NASAMass     float64 `json:"nasa_mass"`
	Cached       bool    `json:"cached"`
}

type AccessibleVolumeResponse struct {
	UnitcellVolume float64            `json:"unitcell_volume"`
	Density        float64            `json:"density"`
	AV             map[string]float64 `json:"av"`
	NAV            map[string]float64 `json:"nav"`
	Cached         bool               `json:"cached"`
}

type ProbeVolumeResponse struct {
	POAVUnitcell  float64 `json:"poav_unitcell"`
	POAVFraction  float64 `json:"poav_fraction"`
	POAVMass      float64 `json:"poav_mass"`
	PONAVUnitcell float64 `json:"ponav_unitcell"`
	PONAVFraction float64 `json:"ponav_fraction"`
	PONAVMass     float64 `json:"ponav_mass"`
	Cached        bool    `json:"cached"`
}

type ChannelAnalysisResponse struct {
	Dimension         int     `json:"dimension"`
	IncludedDiameter  float64 `json:"included_diameter"`
	FreeDiameter      float64 `json:"free_diameter"`
	IncludedAlongFree float64 `json:"included_along_free"`
	Cached            bool    `json:"cached"`
}

type FrameworkInfoResponse struct {
	Filename           string                   `json:"filename"`
	Formula            string                   `json:"formula"`
	Segments           int                      `json:"segments"`
	NumberOfFrameworks int                      `json:"number_of_frameworks"`
	NumberOfMolecules  int                      `json:"number_of_molecules"`
	Frameworks         []map[string]interface{} `json:"frameworks"`
	Cached             bool                     `json:"cached"`
}

type BlockingSpheresResponse struct {
	Channels      []interface{} `json:"channels"`
	Pockets       []interface{} `json:"pockets"`
	NodesAssigned []interface{} `json:"nodes_assigned"`
	Raw           string        `json:"raw"`
	Cached        bool          `json:"cached"`
}

type OpenMetalSitesResponse struct {
	OpenMetalSitesCount int  `json:"open_metal_sites_count"`
	Cached              bool `json:"cached"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

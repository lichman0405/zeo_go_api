package models

import (
	"mime/multipart"
)

type AnalysisRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
	ChanRadius    float64               `form:"chan_radius" binding:"omitempty,min=0.1,max=10"`
	Samples       int                   `form:"samples" binding:"omitempty,min=100,max=1000000"`
	Bins          int                   `form:"bins" binding:"omitempty,min=10,max=1000"`
}

type PoreDiameterRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
}

type SurfaceAreaRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
	Samples       int                   `form:"samples" binding:"omitempty,min=100,max=1000000"`
}

type VolumeRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
	ChanRadius    float64               `form:"chan_radius" binding:"omitempty,min=0.1,max=10"`
	Samples       int                   `form:"samples" binding:"omitempty,min=100,max=1000000"`
}

type ChannelRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
}

type PSDRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
	ChanRadius    float64               `form:"chan_radius" binding:"omitempty,min=0.1,max=10"`
	Samples       int                   `form:"samples" binding:"omitempty,min=1000,max=1000000"`
}

type BlockingSpheresRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
	ProbeRadius   float64               `form:"probe_radius" binding:"omitempty,min=0.1,max=10"`
}

type FrameworkInfoRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
}

type OpenMetalSitesRequest struct {
	StructureFile *multipart.FileHeader `form:"structure_file" binding:"required"`
	HA            bool                  `form:"ha" binding:"omitempty"`
}

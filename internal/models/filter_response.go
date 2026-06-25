package models

type IFPFiltersResponse struct {
	Sizes           []string `json:"sizes"`
	Processors      []string `json:"processors"`
	ProcessorSpeeds []string `json:"processorSpeeds"`
	Storages        []string `json:"storages"`
	SmartFeatures   []string `json:"smartFeatures"`
}

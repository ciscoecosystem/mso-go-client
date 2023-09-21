package models

type SiteServiceGraph struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteServiceGraph(ops, path string, siteServiceGraphList []interface{}) *PatchPayloadList {

	return &PatchPayloadList{
		Ops:   ops,
		Path:  path,
		Value: siteServiceGraphList,
	}

}

package models

type SchemaSiteExternalEpg struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteExternalEpg(ops, path string, externalEpgRef map[string]interface{}, l3outRef map[string]interface{}) *SchemaSiteExternalEpg {
	var externalepgMap map[string]interface{}
	externalepgMap = map[string]interface{}{
		"externalEpgRef": externalEpgRef,
	}
	if l3outRef != nil {
		externalepgMap["l3outRef"] = l3outRef
	}
	return &SchemaSiteExternalEpg{
		Ops:   ops,
		Path:  path,
		Value: externalepgMap,
	}
}

func (schemaSiteExternalEpg *SchemaSiteExternalEpg) ToMap() (map[string]interface{}, error) {
	schemaSiteExternalEpgMap := make(map[string]interface{})

	A(schemaSiteExternalEpgMap, "op", schemaSiteExternalEpg.Ops)
	A(schemaSiteExternalEpgMap, "path", schemaSiteExternalEpg.Path)
	if schemaSiteExternalEpg.Value != nil {
		A(schemaSiteExternalEpgMap, "value", schemaSiteExternalEpg.Value)
	}

	return schemaSiteExternalEpgMap, nil
}

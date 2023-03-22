package models

type SchemaSiteAnpEpgBulkStaticPort struct {
	Ops   string        `json:",omitempty"`
	Path  string        `json:",omitempty"`
	Value []interface{} `json:",omitempty"`
}

func NewSchemaSiteAnpEpgBulkStaticPort(ops, path string, staticPortsList []interface{}) *SchemaSiteAnpEpgBulkStaticPort {

	return &SchemaSiteAnpEpgBulkStaticPort{
		Ops:   ops,
		Path:  path,
		Value: staticPortsList,
	}

}

func (anpAttributes *SchemaSiteAnpEpgBulkStaticPort) ToMap() (map[string]interface{}, error) {
	anpAttributesMap := make(map[string]interface{})
	A(anpAttributesMap, "op", anpAttributes.Ops)
	A(anpAttributesMap, "path", anpAttributes.Path)
	if anpAttributes.Value != nil {
		A(anpAttributesMap, "value", anpAttributes.Value)
	}

	return anpAttributesMap, nil
}

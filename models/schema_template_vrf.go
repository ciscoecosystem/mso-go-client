package models

type SchemaTemplateVrf struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaTemplateVrf(ops, path, Name, displayName, ipDataPlaneLearning string, l3m, vzany, preferredGroup bool) *SchemaSite {
	var VrfMap map[string]interface{}

	VrfMap = map[string]interface{}{
		"displayName":         displayName,
		"name":                Name,
		"l3MCast":             l3m,
		"vzAnyEnabled":        vzany,
		"ipDataPlaneLearning": ipDataPlaneLearning,
		"preferredGroup":      preferredGroup,
	}

	return &SchemaSite{
		Ops:   ops,
		Path:  path,
		Value: VrfMap,
	}

}

func (schematemplatevrfAttributes *SchemaTemplateVrf) ToMap() (map[string]interface{}, error) {
	schematemplatevrfAttributeMap := make(map[string]interface{})
	A(schematemplatevrfAttributeMap, "op", schematemplatevrfAttributes.Ops)
	A(schematemplatevrfAttributeMap, "path", schematemplatevrfAttributes.Path)
	if schematemplatevrfAttributes.Value != nil {
		A(schematemplatevrfAttributeMap, "value", schematemplatevrfAttributes.Value)
	}

	return schematemplatevrfAttributeMap, nil
}

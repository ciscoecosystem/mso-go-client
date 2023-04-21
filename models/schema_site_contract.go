package models

type SiteContract struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func NewSchemaSiteContract(ops, path string, contractRef map[string]interface{}) *SiteContract {
	contractMap := map[string]interface{}{
		"contractRef": contractRef,
	}

	return &SiteContract{
		Ops:   ops,
		Path:  path,
		Value: contractMap,
	}
}

func (contractAttributes *SiteContract) ToMap() (map[string]interface{}, error) {
	contractAttributesMap := make(map[string]interface{})
	A(contractAttributesMap, "op", contractAttributes.Ops)
	A(contractAttributesMap, "path", contractAttributes.Path)
	if contractAttributes.Value != nil {
		A(contractAttributesMap, "value", contractAttributes.Value)
	}

	return contractAttributesMap, nil
}

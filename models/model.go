package models

type Model interface {
	ToMap() (map[string]interface{}, error)
}

type patchPayload struct {
	Ops   string                 `json:",omitempty"`
	Path  string                 `json:",omitempty"`
	Value map[string]interface{} `json:",omitempty"`
}

func (patchPayloadAttributes *patchPayload) ToMap() (map[string]interface{}, error) {
	patchPayloadAttributesMap := make(map[string]interface{})
	A(patchPayloadAttributesMap, "op", patchPayloadAttributes.Ops)
	A(patchPayloadAttributesMap, "path", patchPayloadAttributes.Path)
	if patchPayloadAttributes.Value != nil {
		A(patchPayloadAttributesMap, "value", patchPayloadAttributes.Value)
	}
	return patchPayloadAttributesMap, nil
}

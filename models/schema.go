package models

type Schema struct {
	Id          string `json:",omitempty"`
	DisplayName string `json:",omitempty"`

	Templates []map[string]interface{} `json:",omitempty"`

	Sites []map[string]interface{} `json:",omitempty"`
}

type SchemaReplace struct {
	Ops   string                   `json:",omitempty"`
	Path  string                   `json:",omitempty"`
	Value []map[string]interface{} `json:",omitempty"`
}

func NewSchema(id, displayName, templateName, tenantId string, template []interface{}) (*Schema, *SchemaReplace) {
	result := []map[string]interface{}{}
	if templateName != "" {
		templateMap := map[string]interface{}{
			"name":          templateName,
			"tenantId":      tenantId,
			"displayName":   templateName,
			"anps":          []interface{}{},
			"contracts":     []interface{}{},
			"vrfs":          []interface{}{},
			"bds":           []interface{}{},
			"filters":       []interface{}{},
			"externalEpgs":  []interface{}{},
			"serviceGraphs": []interface{}{},
		}
		result = []map[string]interface{}{
			templateMap,
		}
	} else {
		templateMap := make(map[string]interface{})
		for _, map_values := range template {
			map_template := make(map[string]interface{})
			map_template_values := map_values.(map[string]interface{})
			map_template["name"] = map_template_values["name"]
			map_template["displayName"] = map_template_values["displayName"]
			map_template["tenantId"] = map_template_values["tenantId"]

			templateMap = map[string]interface{}{
				"anps":          []interface{}{},
				"contracts":     []interface{}{},
				"vrfs":          []interface{}{},
				"bds":           []interface{}{},
				"filters":       []interface{}{},
				"externalEpgs":  []interface{}{},
				"serviceGraphs": []interface{}{},
			}
			result = append(result, mergeMaps(map_template, templateMap))
		}
	}
	if id == "" {
		return &Schema{
			Id:          id,
			DisplayName: displayName,
			Templates:   result,
			Sites:       []map[string]interface{}{},
		}, nil
	} else {
		return nil, &SchemaReplace{
			Ops:   "replace",
			Path:  "/templates",
			Value: result,
		}
	}
}

func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, individualMap := range maps {
		for k, v := range individualMap {
			result[k] = v
		}
	}
	return result
}

func (schema *Schema) ToMap() (map[string]interface{}, error) {
	schemaAttributeMap := make(map[string]interface{})
	A(schemaAttributeMap, "id", schema.Id)
	A(schemaAttributeMap, "displayName", schema.DisplayName)
	A(schemaAttributeMap, "templates", schema.Templates)
	A(schemaAttributeMap, "sites", schema.Sites)

	return schemaAttributeMap, nil
}

func (schemaReplace *SchemaReplace) ToMap() (map[string]interface{}, error) {
	schemaReplaceMap := make(map[string]interface{})
	A(schemaReplaceMap, "op", schemaReplace.Ops)
	A(schemaReplaceMap, "path", schemaReplace.Path)
	if schemaReplace.Value != nil {
		A(schemaReplaceMap, "value", schemaReplace.Value)
	}

	return schemaReplaceMap, nil
}

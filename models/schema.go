package models

type SchemaAttributes struct {
	Schema string `json:",omitempty"`

	Templates []interface{} `json:",omitempty"`

	Sites []string `json:",omitempty"`
}

type Templates struct {
	Name          string   `json:",omitempty"`
	TenantId      string   `json:",omitempty"`
	DisplayName   string   `json:",omitempty"`
	Anps          []string `json:",omitempty"`
	Contracts     []string `json:",omitempty"`
	Vrfs          []string `json:",omitempty"`
	Bds           []string `json:",omitempty"`
	Filter        []string `json:",omitempty"`
	ExternalEpgs  []string `json:",omitempty"`
	ServiceGraphs []string `json:",omitempty"`
}

func NewSchemacontainer(schemaAttr SchemaAttributes) *SchemaAttributes {

	SchemacontainerAttributes := schemaAttr
	return &SchemacontainerAttributes
}

func (schemaAttributes *SchemaAttributes) ToMap() (map[string]interface{}, error) {
	schemaAttributeMap := make(map[string]interface{})
	A(schemaAttributeMap, "displayName", schemaAttributes.Schema)
	A(schemaAttributeMap, "templates", schemaAttributes.Templates)
	A(schemaAttributeMap, "sites", schemaAttributes.Sites)

	return schemaAttributeMap, nil
}

// func TenantFromContainerList(cont *container.Container, index int) *Tenant {

// 	TenantCont := cont.S("imdata").Index(index).S(FvtenantClassName, "attributes")
// 	return &Tenant{
// 		BaseAttributes{
// 			DistinguishedName: G(TenantCont, "dn"),
// 			Description:       G(TenantCont, "descr"),
// 			Status:            G(TenantCont, "status"),
// 			ClassName:         FvtenantClassName,
// 			Rn:                G(TenantCont, "rn"),
// 		},

// 		TenantAttributes{

// 			Name : G(TenantCont, "name"),

// 	        Annotation : G(TenantCont, "annotation"),

// 	        NameAlias : G(TenantCont, "nameAlias"),

//         },

// 	}
// }

// func TenantFromContainer(cont *container.Container) *Tenant {

// 	return TenantFromContainerList(cont, 0)
// }

// func TenantListFromContainer(cont *container.Container) []*Tenant {
// 	length, _ := strconv.Atoi(G(cont, "totalCount"))

// 	arr := make([]*Tenant, length)

// 	for i := 0; i < length; i++ {

// 		arr[i] = TenantFromContainerList(cont, i)
// 	}

// 	return arr
// }

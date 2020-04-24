package models

type Location struct {
	Latitude  float64 `json:",omitempty"`
	Longitude float64 `json:",omitempty"`
}

type SiteAttributes struct {
	Name         string        `json:",omitempty"`
	ApicUsername string        `json:",omitempty"`
	ApicPassword string        `json:",omitempty"`
	ApicSiteId   string        `json:",omitempty"`
	Site         string        `json:",omitempty"`
	Labels       []interface{} `json:",omitempty"`
	Location     *Location     `json:",omitempty"`
	Url          []interface{} `json:",omitempty"`
}

func NewSite(siteAttr SiteAttributes) *SiteAttributes {

	SiteAttributes := siteAttr
	return &SiteAttributes
}

func (siteAttributes *SiteAttributes) ToMap() (map[string]interface{}, error) {
	siteAttributeMap := make(map[string]interface{})
	A(siteAttributeMap, "name", siteAttributes.Name)
	A(siteAttributeMap, "username", siteAttributes.ApicUsername)
	A(siteAttributeMap, "password", siteAttributes.ApicPassword)
	A(siteAttributeMap, "apic_site_id", siteAttributes.ApicSiteId)
	A(siteAttributeMap, "site", siteAttributes.Site)
	A(siteAttributeMap, "labels", siteAttributes.Labels)
	A(siteAttributeMap, "location", siteAttributes.Location)
	A(siteAttributeMap, "urls", siteAttributes.Url)

	return siteAttributeMap, nil
}

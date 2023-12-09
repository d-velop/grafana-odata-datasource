package odata

import "encoding/xml"

const (
	EdmString         = "Edm.String"
	EdmBoolean        = "Edm.Boolean"
	EdmSingle         = "Edm.Single"
	EdmDouble         = "Edm.Double"
	EdmDecimal        = "Edm.Decimal"
	EdmSByte          = "Edm.SByte"
	EdmByte           = "Edm.Byte"
	EdmInt16          = "Edm.Int16"
	EdmInt32          = "Edm.Int32"
	EdmInt64          = "Edm.Int64"
	EdmDateTimeOffset = "Edm.DateTimeOffset"
	EdmGuid           = "Edm.Guid"
	EdmTime           = "Edm.Time"
	EdmDate           = "Edm.Date"
	EdmDateTime       = "Edm.DateTime"

	Metadata = "$metadata"
	Filter   = "$filter"
	Select   = "$select"
)

type Response struct {
	D *struct {
		Results []map[string]interface{} `json:"results"`
	} `json:"d"`
	Results []map[string]interface{} `json:"results"`
	Value   []map[string]interface{} `json:"value"`
}

type Edmx struct {
	XMLName      xml.Name        `xml:"Edmx"`
	Version      string          `xml:"Version,attr"`
	XmlNs        string          `xml:"edmx,attr"`
	DataServices []*DataServices `xml:"DataServices"`
}

type DataServices struct {
	XMLName xml.Name  `xml:"DataServices"`
	Schemas []*Schema `xml:"Schema"`
}

type Schema struct {
	XMLName          xml.Name           `xml:"Schema"`
	Namespace        string             `xml:"Namespace,attr"`
	XmlNs            string             `xml:"xmlns,attr"`
	EntityTypes      []*EntityType      `xml:"EntityType"`
	EntityContainers []*EntityContainer `xml:"EntityContainer"`
}

type EntityType struct {
	XMLName    xml.Name    `xml:"EntityType"`
	Name       string      `xml:"Name,attr"`
	Key        []*Key      `xml:"Key"`
	Properties []*Property `xml:"Property"`
}

type Key struct {
	XMLName     xml.Name       `xml:"Key"`
	PropertyRef []*PropertyRef `xml:"PropertyRef"`
}

type PropertyRef struct {
	XMLName xml.Name `xml:"PropertyRef"`
	Name    string   `xml:"Name,attr"`
}

type Property struct {
	XMLName  xml.Name `xml:"Property"`
	Name     string   `xml:"Name,attr"`
	Type     string   `xml:"Type,attr"`
	Nullable string   `xml:"Nullable,attr"`
}

type EntityContainer struct {
	XMLName   xml.Name     `xml:"EntityContainer"`
	Name      string       `xml:"Name,attr"`
	EntitySet []*EntitySet `xml:"EntitySet"`
}

type EntitySet struct {
	XMLName    xml.Name `xml:"EntitySet"`
	Name       string   `xml:"Name,attr"`
	EntityType string   `xml:"EntityType,attr"`
}

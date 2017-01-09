package xml

import (
	"bytes"
	corexml "encoding/xml"
	"fmt"
)

type Node struct {
	XMLName corexml.Name
	Attrs   map[string]string `xml:"-"`
	Content []byte            `xml:",innerxml"`
	Nodes   []Node            `xml:",any"`
}

func (n *Node) UnmarshalXML(d *corexml.Decoder, start corexml.StartElement) error {
	if n.Attrs == nil {
		n.Attrs = map[string]string{}
	}
	for _, a := range start.Attr {
		n.Attrs[a.Name.Local] = a.Value
	}

	type node Node
	return d.DecodeElement((*node)(n), &start)
}

func Unmarshal(xml string) (*Node, error) {
	buf := bytes.NewBuffer([]byte(xml))
	dec := corexml.NewDecoder(buf)

	var parsedXML Node
	if err := dec.Decode(&parsedXML); err != nil {
		return nil, fmt.Errorf("Failed to parse XML: %s", err.Error())
	}

	return &parsedXML, nil
}

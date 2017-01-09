package xml

import "fmt"

func Equal(xml1 string, xml2 string) error {
	parsedXML1, err := Unmarshal(xml1)
	if err != nil {
		return err
	}

	parsedXML2, err := Unmarshal(xml2)
	if err != nil {
		return err
	}

	_, err = nodeEq(parsedXML1, parsedXML2)
	return err
}

func nodeEq(node *Node, node2 *Node) (int, error) {
	// Returned on failure - How close the match was
	var depth = 0

	// Check that all names match
	if node.XMLName != node2.XMLName {
		return depth, fmt.Errorf("Could not find a node with the same name as %q", node.XMLName.Local)
	}
	depth++

	// Check that all attributes match
	if err := attrsEq(node, node2); err != nil {
		return depth, err
	}
	depth++

	// For each node, check that there's a matching one on the other side
	for _, n := range node.Nodes {
		var notFoundErr error
		var notFoundDepth = 0
		for i := len(node2.Nodes) - 1; i >= 0; i-- {
			if depth, err := nodeEq(&n, &node2.Nodes[i]); err == nil {
				// We've found a match - delete this node to stop it being match against again
				node2.Nodes = append(node2.Nodes[:i], node2.Nodes[i+1:]...)
				notFoundErr = nil
				break
			} else if depth > notFoundDepth {
				notFoundDepth = depth
				notFoundErr = err
			}
		}
		if notFoundErr != nil {
			return notFoundDepth, notFoundErr
		}
	}

	if len(node2.Nodes) != 0 {
		return depth, fmt.Errorf(
			"Node %q has a different number of child nodes than %q",
			node.XMLName.Local,
			node2.XMLName.Local,
		)
	}

	// There are no extra nodes on either side - increase the depth
	depth++

	// We're at the bottom (no more children) - Check if the content matches
	content1 := string(node.Content)
	content2 := string(node2.Content)
	if len(node.Nodes) == 0 && len(node2.Nodes) == 0 && content1 != content2 {
		return depth, fmt.Errorf("Content %q does not match %q", content1, content2)
	}

	return 0, nil
}

func attrsEq(node *Node, node2 *Node) error {
	for k, v := range node.Attrs {
		if v != node2.Attrs[k] {
			return fmt.Errorf("Attribute mismatch - %q != %q", v, node2.Attrs[k])
		}
	}

	if len(node.Attrs) != len(node2.Attrs) {
		return fmt.Errorf(
			"Number of attributes in node %q (%d) != node %q (%d)",
			node.XMLName.Local,
			len(node.Attrs),
			node2.XMLName.Local,
			len(node2.Attrs),
		)
	}

	return nil
}

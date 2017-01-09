package xml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItIgnoresSpacing(t *testing.T) {
	t.Parallel()

	xml1 := `<outer><inner>abc</inner></outer>`
	xml2 := ` <outer>    <inner>abc</inner>
	</outer>  `

	err := Equal(xml1, xml2)

	require.NoError(t, err)
}

func TestItIgnoresAttributeOrder(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
		</outer>
	`
	xml2 := `
		<outer c="d" a="b">
			<inner h="i" f="g">abc</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.NoError(t, err)
}

func TestItDetectsAttributeMismatch(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
		</outer>
	`
	xml2 := `
		<outer c="d" a="d">
			<inner f="g" h="WRONG">abc</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, `Attribute mismatch - "b" != "d"`)
}

func TestItDetectsValueMismatch(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
		</outer>
	`
	xml2 := `
		<outer a="b" c="d">
			<inner f="g" h="i">WRONG</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, `Content "abc" does not match "WRONG"`)
}

func TestItDetectsExtraAttribute(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
		</outer>
	`
	xml2 := `
		<outer a="b" c="d">
			<inner f="g" h="i" j="k">abc</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, `Number of attributes in node "inner" (2) != node "inner" (3)`)
}

func TestItDetectsExtraNode(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
		</outer>
	`
	xml2 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
			<inner f="g" h="i">abc</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, `Node "outer" has a different number of child nodes than "outer"`)
}

func TestItChoosesClosestMatchForErrorMessage(t *testing.T) {
	t.Parallel()

	xml1 := `
		<outer a="b" c="d">
			<inner f="g" h="i">abc</inner>
			<inner x="y" z="z">xyz</inner>
		</outer>
	`
	xml2 := `
		<outer a="b" c="d">
			<inner f="not this one" h="because the attrs are wrong">abc</inner>
			<inner f="g" h="i">this one is a closer match</inner>
		</outer>
	`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, `Content "abc" does not match "this one is a closer match"`)
}

func TestItErrorsOnBadInput(t *testing.T) {
	t.Parallel()

	xml1 := `{ "this_is": "Not XML" }`
	xml2 := `lol`

	err := Equal(xml1, xml2)

	require.EqualError(t, err, "Failed to parse first XML: EOF")
}

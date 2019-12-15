package markdown

import "testing"

func TestMarkdown_OrderedList(t *testing.T) {
	m := Markdown{}
	things := []interface{}{
		"thing",
		"other",
		[]interface{}{
			"hurrah",
		},
		"happy newyears!:",
	}
	m.OrderedList(things)
	t.Log(m.String())
}

func TestMarkdown_Table(t *testing.T) {
	table := Table{Headers: []Renderer{&Text{"Hello"}, &Text{"World"}}, Alignments: []Renderer{&LeftAligned{}, &CenterAligned{}}, Rows: [][]Renderer{
		{&Italics{"YO"}, &Bold{"YO"}},
	}}
	m := Markdown{}
	m.Table(table)
	t.Log(m.String())
}

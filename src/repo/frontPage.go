package repo

import (
	"github.com/christopher-henderson/OneCRL-Viewer/git"
	"github.com/christopher-henderson/OneCRL-Viewer/markdown"
	"io/ioutil"
	"path/filepath"
)

func NewFrontPage(db []CanonicalEntry) error {
	var TABLE = markdown.Table{}
	TABLE.Headers = []markdown.Renderer{
		&markdown.Text{"Issuer"}, &markdown.Text{"Serial"}, &markdown.Text{"Fingerprint"}, &markdown.Text{"Last Modified"},
	}
	TABLE.Alignments = []markdown.Renderer{
		&markdown.LeftAligned{}, &markdown.CenterAligned{}, &markdown.CenterAligned{}, &markdown.RightAligned{},
	}
	TABLE.Rows = make([][]markdown.Renderer, 0)
	g := git.NewRepo(Repo)
	head, err := g.HEAD()
	if err != nil {
		return err
	}
	for _, entry := range db {
		TABLE.Rows = append(TABLE.Rows, entry.FrontPageRow(head))
	}
	return ioutil.WriteFile(filepath.Join(Repo, "README.MD"), []byte((&markdown.Markdown{}).Table(TABLE).String()), 0666)
}

package main // import "github.com/christopher-henderson/OneCRL-Viewer"
import (
	"github.com/christopher-henderson/OneCRL-Viewer/crtSh"
	"github.com/christopher-henderson/OneCRL-Viewer/kinto"
	"github.com/christopher-henderson/OneCRL-Viewer/repo"
	"log"
)

const certDir = `H:\TestRepo`

func main() {
	repo.Init(certDir)
	if err := kinto.Init(certDir); err != nil {
		panic(err)
	}
	changes, err := kinto.Changes()
	if err != nil {
		log.Println(err)
		return
	}
	if !changes.Changed() {
		log.Println("no changes")
		return
	}
	certs, err := crtSh.GetCerts(changes.Serials())
	if err != nil {
		log.Println(err)
		return
	}
	changes.Associate(certs)
	log.Println(repo.Update(changes.State()))
}

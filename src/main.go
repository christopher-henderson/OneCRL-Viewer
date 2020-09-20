package main // import "github.com/christopher-henderson/OneCRL-Viewer"
import (
	"fmt"
	"github.com/christopher-henderson/OneCRL-Viewer/crtSh"
	"github.com/christopher-henderson/OneCRL-Viewer/kinto"
	"github.com/christopher-henderson/OneCRL-Viewer/repo"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run main.go </path/to/repo>")
		os.Exit(1)
	}
	certDir := os.Args[1]
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
	err = repo.Update(changes.State())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

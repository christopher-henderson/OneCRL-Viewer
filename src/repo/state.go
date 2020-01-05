package repo

import (
	"github.com/christopher-henderson/OneCRL-Viewer/kinto"
)

type CertRepo struct {
}

type Cert struct {
	Kinto       kinto.KintoEntry
	Certificate []byte
	Dirname     string
}

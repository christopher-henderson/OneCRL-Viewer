package main // import "github.com/christopher-henderson/OneCRL-Viewer"
import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/christopher-henderson/OneCRL-Viewer/markdown"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"
)

var db *sql.DB

func init() {
	// https://github.com/lib/pq/issues/389
	// binary_parameters=yes is because of "pq: unnamed prepared statement does not exist"
	connStr := "postgres://guest@crt.sh/certwatch?sslmode=verify-full&binary_parameters=yes"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

const Kinto = `https://firefox.settings.services.mozilla.com/v1/buckets/blocklists/collections/certificates/records`

type KintoData struct {
	Data []KintoEntry `json:"data"`
}

type KintoEntry struct {
	Schema       int64   `json:"schema"`
	Details      Details `json:"details"`
	Enabled      bool    `json:"enabled"`
	IssuerName   string  `json:"issuerName"`
	SerialNumber string  `json:"serialNumber"`
	ID           string  `json:"id"`
	LastModified int64   `json:"lastModified"`
}

type Details struct {
	Bug     string `json:"bug"`
	Who     string `json:"who"`
	Why     string `json:"why"`
	Name    string `json:"name"`
	Created string `json:"created"`
}

func (k *KintoEntry) ReadableIssuer() string {
	i, err := base64.StdEncoding.DecodeString(k.IssuerName)
	if err != nil {
		panic(err)
	}
	var name pkix.RDNSequence
	_, err = asn1.Unmarshal(i, &name)
	if err != nil {
		panic(err)
	}
	return name.String()
}

func (k *KintoEntry) HexSerial() string {
	i, err := base64.StdEncoding.DecodeString(k.SerialNumber)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", i)
}

func (k *KintoEntry) Timestamp() string {
	return time.Unix(0, k.LastModified).String()
}

const certDir = `H:\OneCRL-Viewer\certs`

func deser(entry KintoEntry) {

}

//SELECT c.ISSUER_CA_ID,
//        NULL::text ISSUER_NAME,
//        encode(x509_serialNumber(c.CERTIFICATE), 'hex') NAME_VALUE,
//        min(c.ID) MIN_CERT_ID,
//        min(ctle.ENTRY_TIMESTAMP) MIN_ENTRY_TIMESTAMP,
//        x509_notBefore(c.CERTIFICATE) NOT_BEFORE,
//        x509_notAfter(c.CERTIFICATE) NOT_AFTER
//    FROM ct_log_entry ctle,
//        certificate c
//    WHERE c.ID = ctle.CERTIFICATE_ID
//        AND x509_serialNumber(c.CERTIFICATE) = decode('016d05b10de8d1d0e3f660560a6a9b', 'hex')
//    GROUP BY c.ID, c.ISSUER_CA_ID, ISSUER_NAME, NAME_VALUE
//    ORDER BY MIN_ENTRY_TIMESTAMP DESC, NAME_VALUE, ISSUER_NAME;

// select certificate from certificate c where x509_serialNumber(c.certificate) = decode('6488b3ffd2c6bfb39d3bf05a9fc054500a8d7723', 'hex');

var TABLE = markdown.Table{}

//https://groups.google.com/forum/#!topic/crtsh/sUmV0mBz8bQ
func main() {
	TABLE.Headers = []markdown.Renderer{
		&markdown.Text{"Issuer"}, &markdown.Text{"Serial"}, &markdown.Text{"Fingerprint"}, &markdown.Text{"Last Modified"},
	}
	TABLE.Alignments = []markdown.Renderer{
		&markdown.LeftAligned{}, &markdown.CenterAligned{}, &markdown.CenterAligned{}, &markdown.RightAligned{},
	}
	TABLE.Rows = make([][]markdown.Renderer, 0)
	resp, err := http.Get(Kinto)
	if err != nil {
		panic(err)
	}
	var kinto KintoData
	err = json.NewDecoder(resp.Body).Decode(&kinto)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	defer func() {
		readme := markdown.Markdown{}
		readme.Table(TABLE)
		ioutil.WriteFile(`H:\OneCRL-Viewer\README.md`, []byte(readme.String()), 0666)
	}()
	for _, k := range kinto.Data {
		download(k)
	}
	//t := markdown.Table{}
	//t.Headers = []markdown.Renderer{
	//	&markdown.Text{"Issuer"}, &markdown.Text{"Serial"}, &markdown.Text{"Last Modified"},
	//}
	//t.Alignments = []markdown.Renderer{
	//	&markdown.LeftAligned{}, &markdown.CenterAligned{}, &markdown.RightAligned{},
	//}
	//t.Rows = make([][]markdown.Renderer, 0)
	//for _, cert := range kinto.Data {
	//	row := []markdown.Renderer{
	//		&markdown.Text{cert.ReadableIssuer()},
	//		&markdown.Text{cert.HexSerial()},
	//		&markdown.Text{cert.Timestamp()},
	//	}
	//	t.Rows = append(t.Rows, row)
	//}
	//content := (&markdown.Markdown{}).Table(t).String()
	//err = ioutil.WriteFile(`H:\OneCRL-Viewer\thing.md`, []byte(content), 0666)
	//if err != nil {
	//	panic(err)
	//}
}

func appendRow(k KintoEntry, fingerprint, url string) {
	row := []markdown.Renderer{
		&markdown.Text{k.ReadableIssuer()},
		&markdown.Text{k.HexSerial()},
		&markdown.Text{(&markdown.Link{fingerprint, url}).String()},
		&markdown.Text{k.Timestamp()},
	}
	TABLE.Rows = append(TABLE.Rows, row)
}

func download(k KintoEntry) {
	rows, err := db.Query(`select certificate from certificate c where x509_serialNumber(c.certificate) = decode($1, 'hex')`, k.HexSerial())
	if err != nil {
		panic(err)
	}
	thisIssuer := k.ReadableIssuer()
	var found bool
	for rows.Next() {
		var cert []byte
		rows.Scan(&cert)
		c, err := x509.ParseCertificate(cert)
		if err != nil {
			panic(err)
		}
		if c.Issuer.ToRDNSequence().String() != thisIssuer {
			fmt.Printf("WARN: duplicate serial from different CA %s \n", k.HexSerial())
			continue
		}
		if found {
			fmt.Printf("WARN: duplicate serial from SAME CA!!!! %s \n", k.HexSerial())
			continue
		}
		save(k, c)
	}
}

const certPath = `H:\OneCRL-Viewer\certs`

func save(k KintoEntry, c *x509.Certificate) {
	// Cant's use b64 because it includes filesystem characters
	//hasher := sha256.New()
	//hasher.Write(c.Raw)
	//fingerprintB64 := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	hasher := sha256.New()
	hasher.Write(c.Raw)
	fingerprintHex := fmt.Sprintf("%X", hasher.Sum(nil))
	certDir := path.Join(certPath, fingerprintHex)
	err := os.MkdirAll(certDir, 0777)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(certDir, "cert.pem"), pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw}), 0666)
	if err != nil {
		panic(err)
	}
	kintoJson, err := json.MarshalIndent(k, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(certDir, "kinto.json"), kintoJson, 0666)
	if err != nil {
		panic(err)
	}

	crtSh := fmt.Sprintf("https://crt.sh/?q=%s", fingerprintHex)
	resp, err := http.DefaultClient.Get(crtSh)
	if err != nil {
		panic(err)
	}
	lr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	liveReport := trimCrtShHeader(lr)
	resp.Body.Close()
	m := markdown.Markdown{}
	m.H1().WriteNL(c.Subject.CommonName).
		H3().WriteNL("Snapshot of crt.sh").
		H5().Write("Click ").Link("here", crtSh).WriteNL(" for a live crt.sh report").
		WriteNL().Divider().
		WriteNL(string(liveReport))
	err = ioutil.WriteFile(path.Join(certDir, "README.md"), []byte(m.String()), 0666)
	if err != nil {
		panic(err)
	}
	appendRow(k, fingerprintHex, fmt.Sprintf("https://github.com/christopher-henderson/OneCRL-Viewer/tree/master/certs/%s", fingerprintHex))
}

func s() {
	connStr := "postgres://guest@crt.sh/certwatch?sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`select certificate from certificate c where x509_serialNumber(c.certificate) = decode($1, 'hex')`, "6488b3ffd2c6bfb39d3bf05a9fc054500a8d7723")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var cert []byte
		rows.Scan(&cert)
		c, err := x509.ParseCertificate(cert)
		if err != nil {
			panic(err)
		}
		fmt.Println(c.Issuer)
	}
}

const header = `  <A style="text-decoration:none" href="/"><SPAN class="title">crt.sh</SPAN></A>&nbsp; <SPAN class="whiteongrey">Certificate Search</SPAN>
<BR><BR>`

func trimCrtShHeader(report []byte) []byte {
	b := bytes.ReplaceAll(report, []byte(header), []byte(""))
	return nope.ReplaceAll(b, []byte(""))
}

var nope = regexp.MustCompile(`(?s)<HEAD>.*</HEAD>`)

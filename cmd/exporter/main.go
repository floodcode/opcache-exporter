package main

import (
	"io/ioutil"
	"log"

	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"gopkg.in/alecthomas/kingpin.v2"

	fcgiclient "github.com/tomasen/fcgi_client"
)

// Exporter collects OPcache status from the given FastCGI URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	client *fcgiclient.FCGIClient
}

func (e *Exporter) getOpcacheStatus() {
	env := make(map[string]string)
	env["SCRIPT_FILENAME"] = "/var/www/html/test.php"

	resp, err := e.client.Get(env)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Status code:", resp.StatusCode)

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	log.Println(string(content))
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri string) (*Exporter, error) {
	client, err := fcgiclient.Dial("tcp", uri)
	if err != nil {
		return nil, err
	}

	exporter := &Exporter{
		client: client,
	}

	return exporter, nil
}

func main() {
	var (
		fcgiURI = kingpin.Flag("opcache.fcgi-uri", "Connection string to FastCGI server.").Default("127.0.0.1:9000").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	exporter, err := NewExporter(*fcgiURI)
	if err != nil {
		panic(err.Error())
	}

	exporter.getOpcacheStatus()
}

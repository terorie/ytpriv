package cmd

import (
	"archive/tar"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/terorie/ytwrk/api"
	"github.com/terorie/ytwrk/data"
)

var videoParseRawCmd = cobra.Command{
	Use:   "parseraw",
	Short: "Parse a tar archive generated by dumpraw",
	Args:  cobra.NoArgs,
	Run:   cmdFunc(doVideoParseRaw),
}

func doVideoParseRaw(_ *cobra.Command, _ []string) (err error) {
	enc := json.NewEncoder(os.Stdout)

	rd := tar.NewReader(os.Stdin)
	for {
		header, err := rd.Next()
		if err != nil {
			logrus.WithError(err).Fatal("Tar stream failed")
		}

		if !strings.HasSuffix(header.Name, ".json") {
			logrus.WithField("file", header.Name).
				Warn("Ignoring file")
			continue
		}

		id := header.Name[:len(header.Name)-len(".json")]
		if len(id) != 11 {
			logrus.WithField("file", header.Name).
				Warn("Ignoring file")
			continue
		}

		fileBuf, err := ioutil.ReadAll(rd)
		if err != nil {
			logrus.WithError(err).Fatal("Tar stream failed")
		}

		var v data.Video
		if err := api.ParseVideoBody(&v, fileBuf, nil); err != nil {
			logrus.WithField("id", id).WithError(err).
				Error("Failed to parse video")
			continue
		}

		if err := enc.Encode(&v); err != nil {
			logrus.Fatal(err)
		}

		logrus.Debugf("Got video %s", id)
	}
}

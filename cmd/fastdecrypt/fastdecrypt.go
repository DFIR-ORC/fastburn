package main

/*

fastdecrypt command line tool to decrypt and decode

Build

	make decrypt

Test
   ./fastdecrypt -input share/samples/encrypted_ff_sample/ORC_WorkStation_W11-22000-51_FastFind.7z.p7b -output share/samples/encrypted_ff_sample/ORC_WorkStation_W11-22000-51_FastFind.7zs -key share/samples/encrypted_ff_sample/contoso.com.key -trace


**/

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"fastburn/internal/fastfind"
	"fastburn/internal/utils"

	_ "fastburn/cmd/fastburn/rsrc"

	log "github.com/sirupsen/logrus"
)

func Version() string {
	var (
		time     string
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.time":
				time = s.Value
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if modified {
		return fmt.Sprintf("%s-%s-dirty", time, revision)
	}
	return fmt.Sprintf("%s-%s", time, revision)
}

func PrintUsage() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}

func main() {

	var outputFlag string
	var inputFlag string
	var keyFlag string

	flag.Usage = PrintUsage

	infoFlag := flag.Bool("info", false, "Enable info log level")
	debugFlag := flag.Bool("debug", false, "Enable debug log level")
	traceFlag := flag.Bool("trace", false, "Enable trace log level")
	versionFlag := flag.Bool("version", false, "Show version and exit")

	flag.StringVar(&inputFlag, "input", "", "Specify input filename")
	flag.StringVar(&outputFlag, "output", "", "Specify output filename")
	flag.StringVar(&keyFlag, "key", "", "Specify key filename (PEM encoded non encrypted PKCS8 format)")

	flag.Parse()

	if *versionFlag {
		version := Version()
		fmt.Printf("Fastburnt - version: %s\n", version)
		os.Exit(0)
	}

	utils.SetLogLevel(*infoFlag, *debugFlag, *traceFlag)

	clearText, err := fastfind.DecryptCMSData(keyFlag, inputFlag)
	if err != nil {
		log.Fatalf("PKCS7 decryption of '%s' with key '%s' failed: %v",
			inputFlag, keyFlag, err)
	}

	err = fastfind.UnstreamBuffer(clearText, outputFlag)
	if err != nil {
		log.Fatalf("Stream decoding to '%s' failed: %v", outputFlag, err)
	}

	log.Infof("'%s' decrypted and streamed of  to '%s'", inputFlag, outputFlag)
}

//eof

package main

/*

fastdecrypt command line tool to decrypt and decode

Build

	make decrypt

Test
   ./fastdecrypt -input share/samples/encrypted_ff_sample/ORC_WorkStation_W11-22000-51_FastFind.7z.p7b -output share/samples/encrypted_ff_sample/ORC_WorkStation_W11-22000-51_FastFind.7zs -cert share/samples/encrypted_ff_sample/contoso.com.pem -key share/samples/encrypted_ff_sample/contoso.com.key -trace


**/

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"fastburn/internal/fastfind"
	"fastburn/internal/utils"

	_ "fastburn/cmd/fastburn/rsrc"

	"github.com/cloudflare/cfssl/log"
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
	var certFlag string
	var keyFlag string

	flag.Usage = PrintUsage

	infoFlag := flag.Bool("info", false, "Enable debug mode")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	traceFlag := flag.Bool("trace", false, "Enable trace mode")
	versionFlag := flag.Bool("version", false, "Show version and exit")

	flag.StringVar(&inputFlag, "input", "", "Specify input filename")
	flag.StringVar(&outputFlag, "output", "", "Specify output filename")
	flag.StringVar(&certFlag, "cert", "", "Specify certificate filename (PEM format)")
	flag.StringVar(&keyFlag, "key", "", "Specify key filename (PEM encoded non encrypted PKCS8 format)")

	flag.Parse()

	//args := flag.Args()

	if *versionFlag {
		version := Version()
		fmt.Printf("Fastburnt - version: %s\n", version)
		os.Exit(0)
	}
	/*
		if len(args) != 0 {
			PrintUsage()
			os.Exit(0)
		}*/

	utils.SetLogLevel(*infoFlag, *debugFlag, *traceFlag)

	streamFile := outputFlag + "s" // TODO generate a proper temporary path

	err := fastfind.DecryptPKCS7Container(certFlag, keyFlag, inputFlag, streamFile)
	if err != nil {
		log.Errorf("PKCS7 decryption of '%s' with key '%s' and certificate '%s' to '%s' failed: %v", inputFlag, keyFlag, certFlag, streamFile, err)
		os.Exit(-1)
	}

	err = fastfind.Unstream(streamFile, outputFlag)
	if err != nil {
		log.Errorf("Stream deccoding of '%s' to '%s' failed: %v", inputFlag, streamFile, err)
		os.Exit(-1)
	}

	log.Info("'%s' decrypted and streamed of  to '%s'", inputFlag, outputFlag)
}

//eof

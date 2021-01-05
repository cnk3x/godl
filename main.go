package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)

	var (
		version = filepath.Base(os.Args[0])
	)

	if version == "godl" {
		if len(os.Args) > 2 {
			version = os.Args[2]
		}
	}

	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}

	root, err := goroot(version)
	if err != nil {
		log.Fatalf("%s: %v", version, err)
	}

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "new":
		newTool(root, version)
	case "download":
		download(root, version)
	case "link":
		link(root, version)
	case "-h":
		help(version)
	default:
		run(root, version)
	}
}

func help(version string) {
	log.Printf("%s [new <version>| download | link | version]", version)
	os.Exit(0)
}

func checkDownload(root, version string) {
	if version == "godl" {
		help(version)
	}
	if _, err := os.Stat(filepath.Join(root, unpackedOkay)); err != nil {
		log.Fatalf("%s: not downloaded. Run '%s download' to download to %v", version, version, root)
	}
}

func newTool(root, version string) {
	src, _ := filepath.Abs(os.Args[0])
	godl := "/usr/local/bin/godl"

	if src != godl {
		if _, err := os.Stat(godl); version == "godl" || err != nil {
			if err := fileCopy(godl, src); err != nil {
				log.Fatal(err)
			}
		}
	}

	if version == "godl" {
		os.Exit(0)
	}

	dst := "/usr/local/bin/" + version
	if err := linkDo(godl, dst); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func download(root, version string) {
	if err := install(root, version); err != nil {
		log.Fatalf("%s: download failed: %v", version, err)
	}
	os.Exit(0)
}

func link(root, version string) {
	checkDownload(root, version)
	goos := getOS()
	if goos == "linux" || goos == "darwin" {
		gobin := filepath.Join(root, "bin")
		targetDIR := "/usr/local/bin"

		if err := linkDo(gobin+"/go", targetDIR+"/go"); err != nil {
			log.Fatal(err)
		}

		if err := linkDo(gobin+"/gofmt", targetDIR+"/gofmt"); err != nil {
			log.Fatal(err)
		}
		log.Println("link finished")
		os.Exit(0)
	}
	log.Fatalf("install command only support linux and darwin")
}

func run(root, version string) {
	checkDownload(root, version)
	gobin := filepath.Join(root, "bin", "go"+exe())
	cmd := exec.Command(gobin, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	newPath := filepath.Join(root, "bin")
	if p := os.Getenv("PATH"); p != "" {
		newPath += string(filepath.ListSeparator) + p
	}
	cmd.Env = dedupEnv(caseInsensitiveEnv, append(os.Environ(), "GOROOT="+root, "PATH="+newPath))
	if err := cmd.Run(); err != nil {
		// TODO: return the same exit status maybe.
		os.Exit(1)
	}
	os.Exit(0)
}

package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/vault/shamir"
)

func main() {
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
	parts := splitCmd.Int("parts", 5, "Number of parts to split into")
	threshold := splitCmd.Int("threshold", 3, "Number of parts needed to restore")

	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'split' or 'restore' subcommands")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "split":
		err = splitCmd.Parse(os.Args[2:])
		if err != nil {
			break
		}
		err = split(splitCmd.Arg(0), *parts, *threshold)
	case "restore":
		err = restoreCmd.Parse(os.Args[2:])
		if err != nil {
			break
		}
		err = restore(restoreCmd.Arg(0))
	default:
		err = fmt.Errorf("expected 'split' or 'restore' subcommands")
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func split(filename string, parts int, threshold int) error {
	src, err := source(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	secret, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	byteParts, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return err
	}

	for _, bytePart := range byteParts {
		fmt.Println(base64.StdEncoding.EncodeToString(bytePart))
	}

	return nil
}

func restore(filename string) error {
	src, err := source(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}
	strParts := strings.Split(string(data), "\n")

	byteParts := [][]byte{}
	for _, strPart := range strParts {
		if strPart == "" {
			continue
		}

		bytePart, err := base64.StdEncoding.DecodeString(strPart)
		if err != nil {
			return err
		}

		byteParts = append(byteParts, bytePart)
	}

	secret, err := shamir.Combine(byteParts)
	if err != nil {
		return err
	}

	fmt.Println(string(secret))

	return nil
}

func source(filename string) (io.ReadCloser, error) {
	if filename == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("you have an error in stdin:%s", err)
		}

		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			return nil, fmt.Errorf("nothing to read from stdin")
		}

		return os.Stdin, nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %q, %v", filename, err)
	}
	return f, nil
}

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
	shares := splitCmd.Int("shares", 5, "Number of shares to split into")
	threshold := splitCmd.Int("threshold", 3, "Number of shares needed to restore")

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
		err = split(splitCmd.Arg(0), *shares, *threshold)
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

func split(filename string, shares int, threshold int) error {
	src, err := source(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	secret, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	byteShares, err := shamir.Split(secret, shares, threshold)
	if err != nil {
		return err
	}

	for _, byteShare := range byteShares {
		fmt.Println(base64.StdEncoding.EncodeToString(byteShare))
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
	strShares := strings.Split(string(data), "\n")

	byteShares := [][]byte{}
	for _, strShare := range strShares {
		if strShare == "" {
			continue
		}

		byteShare, err := base64.StdEncoding.DecodeString(strShare)
		if err != nil {
			return err
		}

		byteShares = append(byteShares, byteShare)
	}

	secret, err := shamir.Combine(byteShares)
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
			return nil, fmt.Errorf("you have an error in stdin: %w", err)
		}

		if (stat.Mode() & os.ModeNamedPipe) == 0 {
			return nil, fmt.Errorf("nothing to read from stdin")
		}

		return os.Stdin, nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %q: %w", filename, err)
	}
	return f, nil
}

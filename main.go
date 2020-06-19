package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path"
	"strings"
)

var (
	in        = flag.String("in", "", "The file to be processed.")
	out       = flag.String("out", "", "The name of the output file (will append '_replaced' if none specified).")
	from      = flag.String("from", "", "The variable(s) to replace (comma separated).")
	to        = flag.String("to", "", "The replacement variable(s) (comma separated).")
	overwrite = flag.Bool("overwrite", false, "Specify whether or not to overwrite the output file if it already exists.")
)

func main() {
	flag.Parse()
	doLotsOfShit()
}

func doLotsOfShit() {

	if *in == "" {
		log.Fatalln("You must specify a file to be used as input.")
	}

	if *from == "" {
		log.Fatalln("You must specify at least one 'from' variable.")
	}

	if *to == "" {
		log.Fatalln("You must specify at least one 'to' variable..")
	}

	froms := strings.Split(*from, ",")
	tos := strings.Split(*to, ",")

	if len(froms) != len(tos) {
		log.Fatalln("You must specify the same number of 'from' replacements as 'to' replacements.")
	}

	fin, err := os.Stat(*in)

	if os.IsNotExist(err) {
		log.Fatalf("Could not find input file: %s", *in)
	}

	if fin.Mode().IsDir() {
		log.Fatalln("Specified input file is a directory.")
	}

	infile, err := os.Open(*in)
	defer func() {
		if er := infile.Close(); er != nil {
			log.Fatalf("Error closing input file: %s", er)
		}
	}()

	if err != nil {
		log.Fatalf("Unable to read input file: %s", err)
	}

	outpath := *out

	if *out == "" {
		outpath = strings.TrimSuffix(*in, path.Ext(*in)) + "_replaced" + path.Ext(*in)
	}

	fout, err := os.Stat(outpath)

	if !os.IsNotExist(err) {
		if !*overwrite {
			log.Fatalf("Output file %q already exists and overwrite is set to false. (try --overwrite)", outpath)
		} else {
			if fout.Mode().IsDir() {
				log.Fatalf("Generated output file %q exists and is a directory.", outpath)
			}
		}
	}

	outfile, err := os.OpenFile(outpath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("Could not open or create output file %q: %s", outpath, err)
	}

	defer func() {
		if er := outfile.Close(); er != nil {
			log.Fatalf("Error closing output file: %s", er)
		}
	}()

	scanner := bufio.NewScanner(infile)
	outw := bufio.NewWriter(outfile)

	for scanner.Scan() {
		_, er := outw.WriteString(replace(scanner.Text(), froms, tos))
		if er != nil {
			log.Fatalf("Error writing to output file: %s", er)
		}
	}

	if err = outw.Flush(); err != nil {
		log.Fatalf("Error flushing write buffer to output file: %s", err)
	}

	if err = scanner.Err(); err != nil {
		log.Fatalf("Error reading from input file: %s", err)
	}
}

func replace(in string, froms []string, tos []string) (out string) {
	out = in
	for i, f := range froms {
		out = strings.Replace(out, f, tos[i], -1)
	}
	out += "\n"
	return
}

/*
Copyright 2017 Akira Midorikawa

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func help() {
	fmt.Fprintf(os.Stderr, "Usage: %s [set-resolution | help]\n", os.Args[0])
}

func setResolution(name string, args []string) {
	var input, output string
	var ppi uint

	flagSet := flag.NewFlagSet(name, flag.ExitOnError)

	flagSet.StringVar(&input, "input", "", "Path for input")
	flagSet.StringVar(&output, "output", "", "Path for output")
	flagSet.UintVar(&ppi, "ppi", 0, "Pixel per inch")

	flagSet.Parse(args)

	if input == "" || output == "" || ppi == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		flagSet.PrintDefaults()
		os.Exit(99)
	}

	image, err := Load(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	ppm := math.Floor((float64(ppi) / 0.0254) + 0.5)
	err = image.SetPhysChunk(&PhysChunk{uint32(ppm), uint32(ppm), 1})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", input, err.Error())
		os.Exit(1)
	}

	if err := image.Write(output); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		help()
		os.Exit(99)
	}

	name := os.Args[1]
	args := os.Args[2:]

	switch name {
	case "set-resolution":
		setResolution(name, args)

	case "help":
		help()

	default:
		help()
		os.Exit(99)
	}
}

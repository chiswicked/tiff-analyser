// Copyright 2019 Norbert Metz. All rights reserved.
// Use of this source code is governed by a GNU GPLv3-style
// license that can be found in the LICENSE file.
//
// Simple command-line tool to determine if a given file is
// a valid TIFF file with a single flattened image layer.
//
// See also:
// https://printplanet.com/forum/prepress-and-workflow/adobe/248476-identifying-layered-tifs-in-indesign?p=248649#post248649
// https://www.awaresystems.be/imaging/tiff/tifftags/imagesourcedata.html

package main

import (
	"fmt"
	"os"

	"github.com/google/tiff"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s filename1 [filename2 filename3 ...]\n", os.Args[0])
		os.Exit(0)
	}

	for _, fileName := range os.Args[1:] {
		file, err := os.Open(fileName)
		defer file.Close()

		if err != nil {
			fmt.Println("Cannot open file:", fileName)
			os.Exit(1)
		}

		sym := "✔ "

		if !IsFlattenedTIFF(file) {
			sym = "✘ "
		}

		fmt.Printf("%s %s\n", sym, fileName)
	}

}

// IsFlattenedTIFF Returns whether a given file is a flattened TIFF.
// Returns true if the file is a valid TIFF file and has a single flattened image layer.
// Returns false if the file is an invalid TIFF file or has multiple image layers.
func IsFlattenedTIFF(file *os.File) bool {

	t, err := tiff.Parse(tiff.NewReadAtReadSeeker(file), nil, nil)
	if err != nil {
		// Not a valid TIFF file
		return false
	}

	for _, ifd := range t.IFDs() {
		for _, field := range ifd.Fields() {
			// Not flattened
			if field.Tag().ID() == 37724 && field.Tag().Name() == "ImageSourceData" {
				return false
			}

		}
	}

	return true
}

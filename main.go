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
	"errors"
	"fmt"
	"os"

	"github.com/google/tiff"
)

// Potential errors that indicate incompatibility with Exstream importer
var (
	ErrNotTIFFFile        = errors.New("The file does not appear to be a TIFF file")
	ErrTIFFLayers         = errors.New("The file contains layers")
	ErrCompressedTIFFFile = errors.New("The file appears to be compressed")
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
		errMsg := ""

		if ok, errs := IsExstreamCompatible(file); !ok {
			sym = "✘ "
			for _, msg := range errs {
				errMsg += msg.Error() + "\n"
			}
		}

		fmt.Printf("%s %s\n", sym, fileName)
		if errMsg != "" {
			fmt.Printf("Error:\n%s", errMsg)
		}
	}

}

// IsExstreamCompatible Returns whether a given file is a flattened TIFF.
// Returns true if the file is a valid TIFF file and has a single flattened image layer.
// Returns false if the file is an invalid TIFF file or has multiple image layers.
func IsExstreamCompatible(file *os.File) (bool, []error) {
	errs := []error{}
	t, err := tiff.Parse(tiff.NewReadAtReadSeeker(file), nil, nil)
	if err != nil {
		errs := append(errs, ErrNotTIFFFile)
		return false, errs
	}

	for _, ifd := range t.IFDs() {
		for _, field := range ifd.Fields() {
			// fmt.Println(field.Tag().ID(), ":", field.Tag().Name())

			// Not flattened (has layers)
			if field.Tag().ID() == 37724 && field.Tag().Name() == "ImageSourceData" {
				errs = append(errs, ErrTIFFLayers)
			}

			// Compressed
			if field.Tag().ID() == 317 && field.Tag().Name() == "Predictor" {
				errs = append(errs, ErrCompressedTIFFFile)
			}

			if len(errs) > 0 {
				return false, errs
			}
		}
	}

	return true, nil
}

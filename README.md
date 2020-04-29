# TIFF Analyser
Simple command-line tool to determine if a given file is a valid TIFF file with a single flattened image layer.

## Usage
```
git clone git@github.com:chiswicked/tiff-analyser.git
cd tiff-analyser
go get -v -t -d ./...
go run main.go file1.tif file2.tif
```
# go-picture-cleanup

This is a small utility to resolve the issue of Adobe Lightroom importing not only RAW images but also the accompanying JPEG files (in case your camera writes both format, i.e. you configured it to do so).

I have configured my camera to do so because the JPG files are handy for faster 1:1 magnification (on the camera) - and the RAW files are obviously want you want for further processing. 

One could just run `rm *.jpg` - risking to delete JPEGs for which no RAW file exists (like pictures received from someone else, from your phone or drone).

This simple tool lists all JPEGs for which RAW files exist. 
Matching is done based on the name.
By default `nef`, `arw` and `dng` are considered RAW files - but this can be overridden using the `-raw-exts` flag providing a comma-separated list.

## Installation

```bash
go install github.com/florianloch/go-picture-cleanup@latest
```

## Usage

```bash
# Run in the current directory, the tool checks subdirectories recursively
# By default, it runs in dry-mode and does not delete anything
go-picture-cleanup .

# Once you confirmed only files you want to be deleted are listed run in non-dry-mode
go-picture-cleanup -d .
```

## Disclaimer

This piece of software is provided as-is, no guarantees regarding functionality etc. are given.
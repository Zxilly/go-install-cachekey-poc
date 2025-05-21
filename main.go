package main

import (
	"cmp"
	"crypto/sha256"
	"debug/buildinfo"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"slices"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(("Usage: go run main.go binary"))
	}

	bin := os.Args[1]
	var err error
	if !filepath.IsAbs(bin) {
		// found in path
		bin, err = exec.LookPath(bin)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Get the build info
	bi, err := buildinfo.ReadFile(bin)
	if err != nil {
		log.Fatal(err)
	}

	key := fmt.Sprintf("%s-%s", bi.Path, bi.GoVersion)

	// make hash reproducible
	slices.SortFunc(bi.Deps, func(a, b *debug.Module) int {
		if a.Path != b.Path {
			return cmp.Compare(a.Path, b.Path)
		}
		if a.Version != b.Version {
			return cmp.Compare(a.Version, b.Version)
		}
		return cmp.Compare(a.Sum, b.Sum)
	})

	slices.SortFunc(bi.Settings, func(a, b debug.BuildSetting) int {
		if a.Key != b.Key {
			return cmp.Compare(a.Key, b.Key)
		}
		return cmp.Compare(a.Value, b.Value)
	})

	// Hash the build info
	h := sha256.New().Sum([]byte(bi.String()))
	
	key += fmt.Sprintf("%x", h)
	fmt.Printf("Hash: %x\n", h)
	fmt.Printf("Buildinfo: %s\n", bi.String())
}	

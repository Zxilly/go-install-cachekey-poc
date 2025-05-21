package main

import (
	"cmp"
	"crypto/md5"
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
	h := md5.New()
	h.Write([]byte(bi.String()))
	result := h.Sum(nil)
	
	key += fmt.Sprintf("-%x", result)
	fmt.Printf("Buildinfo: %s\n\n", bi.String())
	fmt.Printf("Hash: %s\n", key)
}	

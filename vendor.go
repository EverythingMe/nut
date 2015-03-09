package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/jingweno/nut/vendor/_nuts/github.com/codegangsta/cli"
)

var vendorCommand = cli.Command{
	Name:   "vendor",
	Usage:  "change all vendored imports in the current project to the vendored paths",
	Action: runVendor,
}

// vendorDir traverses all the go files in a directory and rewrites all the desp in vendoredDeps
// to their vendored path
func vendorDir(pth, root string, vendoredDeps []string) error {

	os.Chdir(pth)

	lister := pkgLister{
		Env: os.Environ(),
	}

	pks, err := lister.List("")
	if err != nil {
		return err
	}

	// we should find only one package here, or this dir might be empty
	if len(pks) == 1 {
		pk := pks[0]
		if root == "" {
			root = pk.ImportPath
		}

		for _, f := range pk.AllGoFiles() {
			fmt.Println("Vendoring deps in", f)
			err := rewriteGoFile(f, root, vendoredDeps)
			if err != nil {
				return err
			}
		}

	}

	files, err := ioutil.ReadDir(pth)

	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), ".") || strings.HasPrefix(file.Name(), "_") || file.Name() == "internal" {
			continue
		}

		err := vendorDir(path.Join(pth, file.Name()), root, vendoredDeps)
		if err != nil {
			fmt.Println("error traversing", file.Name(), ": ", err)
			return err
		}

	}

	return nil

}

func runVendor(c *cli.Context) {

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read the vendored deps from the manifest
	vendored := make([]string, 0)
	for d := range setting.Manifest().Deps {
		vendored = append(vendored, d)
	}

	vendorDir(cwd, "", vendored)

}

//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/janimo/goclick/utils"
)

var buildsourceCommand = clickCommand{"buildsource", "Build a Click source package.", runBuildsource}

var usageBuildsource = "Usage: goclick buildsource [options] DIRECTORY"

var ignorePatternsBuildsource = []string{
	"*.a",
	"*.click",
	"*.la",
	"*.o",
	"*.so",
	".*.sw?",
	"*~",
	",,*",
	".[#~]*",
	".arch-ids",
	".arch-inventory",
	".be",
	".bzr",
	".bzr-builddeb",
	".bzr.backup",
	".bzr.tags",
	".bzrignore",
	".cvsignore",
	".deps",
	".git",
	".gitattributes",
	".gitignore",
	".gitmodules",
	".hg",
	".hgignore",
	".hgsigs",
	".hgtags",
	".shelf",
	".svn",
	"CVS",
	"DEADJOE",
	"RCS",
	"_MTN",
	"_darcs",
	"{arch}",
}

func buildSource(dir, manifestPath string) (string, error) {
	tmpdir, err := ioutil.TempDir("", "click")
	if err != nil {
		utils.ExitError(err)
	}
	defer os.RemoveAll(tmpdir)

	os.Chmod(tmpdir, 0755)

	rootPath := filepath.Join(tmpdir, "source")

	utils.Copytree(dir, rootPath, ignorePatternsBuildsource)

	realManifestPath := filepath.Join(rootPath, "manifest.json")
	manifest := utils.ReadManifest(realManifestPath)

	packageName := fmt.Sprintf("%s_%s.tar.gz", manifest["name"].(string), utils.EpochlessVersion(manifest["version"].(string)))
	packagePath := filepath.Join(dir, packageName)
	utils.CreateFakerootTarball(rootPath, packagePath, nil)

	return packagePath, nil
}

func runBuildsource(args []string) {
	if len(args) == 0 {
		utils.ExitMsg(usageBuildsource)
	}
	directory := args[0]
	fi, err := os.Stat(directory)
	if os.IsNotExist(err) {
		utils.ExitError(err)
	}
	if !fi.IsDir() {
		utils.ExitMsg(directory + " is not a directory")
	}
	manifest := "."
	manifest = filepath.Join(manifest, "manifest.json")
	if !utils.Exists(filepath.Join(directory, manifest)) {
		utils.ExitMsg("directory " + directory + " does not contain manifest file " + manifest)
	}
	path, err := buildSource(directory, manifest)
	if err != nil {
		utils.ExitError(err)
	}
	fmt.Printf("Successfully built source package in '%s'.\n", path)
}

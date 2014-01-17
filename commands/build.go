//
// Copyright (c) 2014 Canonical Ltd.
//
// Author: Jani Monoses <jani@ubuntu.com>
//
package commands

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/janimo/goclick/utils"
)

var buildCommand = clickCommand{"build", "Build a Click package.", runBuild}

var usage = "Usage: goclick build [options] DIRECTORY"

var ignorePatternsBuild = []string{
	"*.click",
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

func writeControl(path string, m utils.Manifest) {
	s := fmt.Sprintf(`Package: %s
Version: %s
Click-Version: %s
Architecture: %s
Maintainer: %s
Installed-Size: %s
Description: %s
`,
		m["name"].(string),
		m["version"].(string),
		utils.SpecVersion,
		m["architecture"].(string),
		m["maintainer"].(string),
		m["installed-size"].(string),
		m["description"].(string))
	err := ioutil.WriteFile(path, []byte(s), 0664)
	if err != nil {
		utils.ExitError(err)
	}
}

func writeMD5Sums(out, dir string) {
	//This is ugly. Make sure dir has trailing "/" so it works well in the strings.Replace below
	dir = utils.AppendSlash(dir)

	m, err := os.Create(out)
	if err != nil {
		utils.ExitError(err)
	}
	defer m.Close()

	h := md5.New()
	utils.Walk(dir,
		func(filename string, fi os.FileInfo) {
			if fi.IsDir() {
				return
			}
			f, err := os.Open(filename)
			if err != nil {
				utils.ExitError(err)
			}
			h.Reset()
			_, err = io.Copy(h, f)
			if err != nil {
				utils.ExitError(err)
			}
			sum := h.Sum(nil)
			fmt.Fprintf(m, "%x  %s\n", sum, strings.Replace(filename, dir, "", 1))
		},
		nil)
}

//Pack all package binary assets into tarballs in an ar archive
func build(dir, manifestPath string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "click")
	if err != nil {
		utils.ExitError(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chmod(tmpDir, 0755)

	dataDir := filepath.Join(tmpDir, "data")

	utils.Copytree(dir, dataDir, ignorePatternsBuild)

	fullManifestpath := filepath.Join(dataDir, "manifest.json")

	manifest := utils.ReadManifest(fullManifestpath)
	manifest["installed-size"] = utils.GetDirSize(dataDir)

	//Do not ship manifest in the data dir, it gets written in the control dir
	os.Remove(fullManifestpath)

	// Make control dir
	controlDir := filepath.Join(tmpDir, "DEBIAN")
	os.Mkdir(controlDir, 0777)

	//Create control file
	controlPath := filepath.Join(controlDir, "control")
	writeControl(controlPath, manifest)

	//Create manifest file
	realManifestPath := filepath.Join(controlDir, "manifest")
	utils.WriteManifest(realManifestPath, manifest)

	//Create preinst file
	preinstPath := filepath.Join(controlDir, "preinst")
	ioutil.WriteFile(preinstPath, []byte(utils.StaticPreinst), 0664)

	//Create md5sums file
	md5sumsPath := filepath.Join(controlDir, "md5sums")
	writeMD5Sums(md5sumsPath, dataDir)

	packageName := fmt.Sprintf("%s_%s_%s.click", manifest["name"].(string), utils.EpochlessVersion(manifest["version"].(string)), manifest["architecture"].(string))

	//Pack them all in an ar file
	packagePath := filepath.Join(dir, packageName)
	pack(tmpDir, controlDir, dataDir, packagePath)
	return packagePath, nil
}

func pack(tmpDir, controlDir, dataDir, packagePath string) {
	dataTarPath := filepath.Join(tmpDir, "data.tar.gz")
	utils.CreateFakerootTarball(dataDir, dataTarPath, nil)
	controlTarPath := filepath.Join(tmpDir, "control.tar.gz")
	utils.CreateTarball(controlDir, controlTarPath, nil)

	arw, _ := utils.NewArWriter(packagePath)
	arw.AddMagic()
	arw.AddData("debian-binary", []byte("2.0\n"))
	arw.AddData("_click-binary", []byte(utils.SpecVersion+"\n"))
	arw.AddFile("control.tar.gz", controlTarPath)
	arw.AddFile("data.tar.gz", dataTarPath)
	arw.Close()
}

func runBuild(args []string) {
	if len(args) == 0 {
		utils.ExitMsg(usage)
	}
	directory := args[0]
	fi, err := os.Stat(directory)
	if os.IsNotExist(err) {
		utils.ExitError(err)
	}
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
	path, err := build(directory, manifest)
	if err != nil {
		utils.ExitError(err)
	}
	fmt.Printf("Successfully built package in '%s'.\n", path)
}

// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

//go:build none
// +build none

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/klaytn/klaytn/utils/build"
)

var (

	// Files that end up in the klay*.zip archive.
	klayArchiveFiles = []string{
		"COPYING",
		executablePath("klay"),
	}

	// Files that end up in the klay-alltools*.zip archive.
	allToolsArchiveFiles = []string{
		"COPYING",
		executablePath("klay"),
	}

	// A debian package is created for all executables listed here.
	debExecutables = []debExecutable{
		{
			Name:        "klay",
			Description: "Klaytn CLI client.",
		},
	}

	// Distros for which packages are created.
	// Note: vivid is unsupported because there is no golang-1.6 package for it.
	// Note: wily is unsupported because it was officially deprecated on lanchpad.
	// Note: yakkety is unsupported because it was officially deprecated on lanchpad.
	// Note: zesty is unsupported because it was officially deprecated on lanchpad.
	debDistros = []string{"trusty", "xenial", "artful", "bionic"}
)

var GOBIN, _ = filepath.Abs(filepath.Join("build", "bin"))

func executablePath(name string) string {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(GOBIN, name)
}

func main() {
	log.SetFlags(log.Lshortfile)

	if _, err := os.Stat(filepath.Join("build", "ci.go")); os.IsNotExist(err) {
		log.Fatal("this script must be run from the root of the repository")
	}
	if len(os.Args) < 2 {
		log.Fatal("need subcommand as first argument")
	}
	switch os.Args[1] {
	case "install":
		doInstall(os.Args[2:])
	case "test":
		doTest(os.Args[2:])
	case "cover":
		doCover(os.Args[2:])
	case "fmt":
		doFmt(os.Args[2:])
	case "lint":
		doLint(os.Args[2:], true)
	case "lint-try":
		doLint(os.Args[2:], false)
	case "archive":
		doArchive(os.Args[2:])
	case "debsrc":
		doDebianSource(os.Args[2:])
	case "nsis":
		doWindowsInstaller(os.Args[2:])
	case "aar":
		doAndroidArchive(os.Args[2:])
	case "xcode":
		doXCodeFramework(os.Args[2:])
	case "xgo":
		doXgo(os.Args[2:])
	default:
		log.Fatal("unknown command ", os.Args[1])
	}
}

// Compiling

func doInstall(cmdline []string) {
	var (
		arch = flag.String("arch", "", "Architecture to cross build for")
		cc   = flag.String("cc", "", "C compiler to cross build with")
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	// Check Go version. People regularly open issues about compilation
	// failure with outdated Go. This should save them the trouble.
	if !strings.Contains(runtime.Version(), "devel") {
		// Figure out the minor version number since we can't textually compare (1.10 < 1.9)
		var minor int
		fmt.Sscanf(strings.TrimPrefix(runtime.Version(), "go1."), "%d", &minor)

		if minor < 9 {
			log.Println("You have Go version", runtime.Version())
			log.Println("klaytn requires at least Go version 1.9 and cannot")
			log.Println("be compiled with an earlier version. Please upgrade your Go installation.")
			os.Exit(1)
		}
	}
	// Compile packages given as arguments, or everything if there are no arguments.
	packages := []string{"./..."}
	if flag.NArg() > 0 {
		packages = flag.Args()
	}

	if *arch == "" || *arch == runtime.GOARCH {
		goinstall := goTool("install", buildFlags(env)...)
		goinstall.Args = append(goinstall.Args, "-v")
		//goinstall.Args = append(goinstall.Args, "-race")
		goinstall.Args = append(goinstall.Args, packages...)
		build.MustRun(goinstall)
		return
	}
	// If we are cross compiling to ARMv5 ARMv6 or ARMv7, clean any previous builds
	if *arch == "arm" {
		os.RemoveAll(filepath.Join(runtime.GOROOT(), "pkg", runtime.GOOS+"_arm"))
		for _, path := range filepath.SplitList(build.GOPATH()) {
			os.RemoveAll(filepath.Join(path, "pkg", runtime.GOOS+"_arm"))
		}
	}
	// Seems we are cross compiling, work around forbidden GOBIN
	goinstall := goToolArch(*arch, *cc, "install", buildFlags(env)...)
	goinstall.Args = append(goinstall.Args, "-v")
	goinstall.Args = append(goinstall.Args, []string{"-buildmode", "archive"}...)
	goinstall.Args = append(goinstall.Args, packages...)
	build.MustRun(goinstall)

	if cmds, err := ioutil.ReadDir("cmd"); err == nil {
		for _, cmd := range cmds {
			pkgs, err := parser.ParseDir(token.NewFileSet(), filepath.Join(".", "cmd", cmd.Name()), nil, parser.PackageClauseOnly)
			if err != nil {
				log.Fatal(err)
			}
			for name := range pkgs {
				if name == "main" {
					gobuild := goToolArch(*arch, *cc, "build", buildFlags(env)...)
					gobuild.Args = append(gobuild.Args, "-v")
					gobuild.Args = append(gobuild.Args, []string{"-o", executablePath(cmd.Name())}...)
					gobuild.Args = append(gobuild.Args, "."+string(filepath.Separator)+filepath.Join("cmd", cmd.Name()))
					build.MustRun(gobuild)
					break
				}
			}
		}
	}
}

func buildFlags(env build.Environment) (flags []string) {
	var ld []string
	if env.Commit != "" {
		ld = append(ld, "-X", "main.gitCommit="+env.Commit)
		ld = append(ld, "-X", "github.com/klaytn/klaytn/cmd/utils/nodecmd.gitCommit="+env.Commit)
	}
	if env.Tag != "" {
		ld = append(ld, "-X", "github.com/klaytn/klaytn/cmd/utils/nodecmd.gitTag="+env.Tag)
	}
	if runtime.GOOS == "darwin" {
		ld = append(ld, "-s")
	}

	if env.IsDisabledSymbolTable {
		ld = append(ld, "-s")
	}
	if env.IsStaticLink {
		// Pass the static link flag to the external linker.
		// By default, cmd/link will use external linking mode when non-standard cgo packages are involved.
		ld = append(ld, "-linkmode", "external", "-extldflags", "-static")
	}
	if env.IsKlaytnRaceDetectionOn {
		flags = append(flags, "-race")
	}
	if len(ld) > 0 {
		flags = append(flags, "-ldflags", strings.Join(ld, " "))
	}
	return flags
}

func goTool(subcmd string, args ...string) *exec.Cmd {
	return goToolArch(runtime.GOARCH, os.Getenv("CC"), subcmd, args...)
}

func goToolArch(arch string, cc string, subcmd string, args ...string) *exec.Cmd {
	cmd := build.GoTool(subcmd, args...)
	cmd.Env = []string{"GOPATH=" + build.GOPATH()}
	if arch == "" || arch == runtime.GOARCH {
		cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
	} else {
		cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
		cmd.Env = append(cmd.Env, "GOARCH="+arch)
	}
	if cc != "" {
		cmd.Env = append(cmd.Env, "CC="+cc)
	}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	return cmd
}

// Running The Tests
//
// "tests" also includes static analysis tools such as vet.

func doTest(cmdline []string) {
	var (
		parallel = flag.Int("p", 0, "The number of parallel test executions (default: the number of CPUs available)")
		excludes = flag.String("exclude", "", "Comma-separated top-level directories to be excluded in test")
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}

	if *excludes != "" {
		packages = build.ExcludePackages(packages, strings.Split(*excludes, ","))
	}

	// Run analysis tools before the tests.
	build.MustRun(goTool("vet", packages...))

	// Run the actual tests.
	gotest := goTool("test", buildFlags(env)...)
	if *parallel != 0 {
		gotest.Args = append(gotest.Args, "-p", strconv.Itoa(*parallel))
	}
	gotest.Args = append(gotest.Args, "--timeout=30m")
	gotest.Args = append(gotest.Args, packages...)
	build.MustRun(gotest)
}

func doCover(cmdline []string) {
	var (
		parallel   = flag.Int("p", 0, "The number of parallel coverage test executions (default: the number of CPUs available)")
		excludes   = flag.String("exclude", "", "Comma-separated top-level directories to be excluded in coverage test")
		outputFile = flag.String("coverprofile", "coverage.out", "The coverage profile file will be generated by coverage test")
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}

	if *excludes != "" {
		packages = build.ExcludePackages(packages, strings.Split(*excludes, ","))
	}

	coverPackages := []string{"./..."}
	coverExcludes := []string{
		"/tests",
		"/metric",
		"/build",
		"/client",
		"/contracts",
		"/simulations",
		"/api",
		"/fork",
		"/mocks",
	}

	coverPackages = build.ExcludePackages(coverPackages, coverExcludes)
	coverPackagesString := strings.Join(coverPackages, ",")

	// Run analysis tools before the tests.
	build.MustRun(goTool("vet", packages...))

	// Generate a coverage output file.
	build.MustRunCommand("sh", "-c", "echo 'mode: atomic' > "+*outputFile)

	// Run the actual tests.
	gotest := goTool("test", buildFlags(env)...)
	if *parallel != 0 {
		gotest.Args = append(gotest.Args, "-p", strconv.Itoa(*parallel))
	}

	gotest.Args = append(gotest.Args, "-cover", "-covermode=atomic", "-coverprofile="+*outputFile)
	gotest.Args = append(gotest.Args, "-coverpkg", coverPackagesString)
	gotest.Args = append(gotest.Args, "--timeout=30m")
	gotest.Args = append(gotest.Args, packages...)
	build.MustRun(gotest)
}

func doFmt(cmdline []string) {
	// runs gometalinter on requested packages
	flag.CommandLine.Parse(cmdline)

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}

	lintBin := installLinter()

	// Run fast linters batched together
	configs := []string{
		"run",
		"--tests",
		"--disable-all",
		"--enable=gofmt",
		"--timeout=2m",
	}
	build.MustRunCommand(lintBin, append(configs, packages...)...)
}

// runs gometalinter on requested packages and exits immediately when linter warning observed if exitOnError is true
func doLint(cmdline []string, exitOnError bool) {
	flag.CommandLine.Parse(cmdline)

	packages := []string{"./..."}
	if len(flag.CommandLine.Args()) > 0 {
		packages = flag.CommandLine.Args()
	}

	lintBin := installLinter()

	// Prepare a report file for linters
	fname := "linter_report.txt"
	fileOut, err := os.Create(fname)
	defer fileOut.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generating a linter report %s using above linters.\n", fname)

	oldStdout := os.Stdout
	os.Stdout = fileOut

	// Run fast linters batched together
	configs := []string{
		"run",
		"--tests",
		"--disable-all",
		"--enable=varcheck",
		"--enable=misspell",
		"--enable=goconst",
	}
	args := append(configs, packages...)
	if exitOnError {
		build.MustRunCommand(lintBin, args...)
	} else {
		build.TryRunCommand(lintBin, args...)
	}

	// Run fast linters batched together
	configs = []string{
		"run",
		"--tests",
		"--disable-all",
		"--enable=deadcode",
		"--enable=dupl",
		"--enable=errcheck",
		"--enable=ineffassign",
		"--enable=interfacer",
		"--enable=unparam",
		"--enable=unused",
	}
	args = append(configs, packages...)
	if exitOnError {
		build.MustRunCommand(lintBin, args...)
	} else {
		build.TryRunCommand(lintBin, args...)
	}

	// Run slow linters one by one
	for _, linter := range []string{"unconvert", "gosimple", "staticcheck", "gocyclo"} {
		configs = []string{"run", "--tests", "--deadline=10m", "--disable-all", "--enable=" + linter}
		args = append(configs, packages...)
		if exitOnError {
			build.MustRunCommand(lintBin, args...)
		} else {
			build.TryRunCommand(lintBin, args...)
		}
	}

	// Restore stdout
	os.Stdout = oldStdout

	fmt.Printf("Succefully generating %s.\n", fname)
}

// Release Packaging

func doArchive(cmdline []string) {
	var (
		arch   = flag.String("arch", runtime.GOARCH, "Architecture cross packaging")
		atype  = flag.String("type", "zip", "Type of archive to write (zip|tar)")
		signer = flag.String("signer", "", `Environment variable holding the signing key (e.g. LINUX_SIGNING_KEY)`)
		upload = flag.String("upload", "", `Destination to upload the archives`)
		ext    string
	)
	flag.CommandLine.Parse(cmdline)
	switch *atype {
	case "zip":
		ext = ".zip"
	case "tar":
		ext = ".tar.gz"
	default:
		log.Fatal("unknown archive type: ", atype)
	}

	var (
		env      = build.Env()
		base     = archiveBasename(*arch, env)
		klaybin  = "klay-" + base + ext
		alltools = "klay-alltools-" + base + ext
	)
	maybeSkipArchive(env)
	if err := build.WriteArchive(klaybin, klayArchiveFiles); err != nil {
		log.Fatal(err)
	}
	if err := build.WriteArchive(alltools, allToolsArchiveFiles); err != nil {
		log.Fatal(err)
	}
	for _, archive := range []string{klaybin, alltools} {
		if err := archiveUpload(archive, *upload, *signer); err != nil {
			log.Fatal(err)
		}
	}
}

func archiveBasename(arch string, env build.Environment) string {
	platform := runtime.GOOS + "-" + arch
	if arch == "arm" {
		platform += os.Getenv("GOARM")
	}
	if arch == "android" {
		platform = "android-all"
	}
	if arch == "ios" {
		platform = "ios-all"
	}
	return platform + "-" + archiveVersion(env)
}

func archiveVersion(env build.Environment) string {
	version := build.VERSION()
	if isUnstableBuild(env) {
		version += "-unstable"
	}
	if env.Commit != "" {
		version += "-" + env.Commit[:8]
	}
	return version
}

func archiveUpload(archive string, blobstore string, signer string) error {
	// If signing was requested, generate the signature files
	if signer != "" {
		pgpkey, err := base64.StdEncoding.DecodeString(os.Getenv(signer))
		if err != nil {
			return fmt.Errorf("invalid base64 %s", signer)
		}
		if err := build.PGPSignFile(archive, archive+".asc", string(pgpkey)); err != nil {
			return err
		}
	}
	return nil
}

// skips archiving for some build configurations.
func maybeSkipArchive(env build.Environment) {
	if env.IsPullRequest {
		log.Printf("skipping because this is a PR build")
		os.Exit(0)
	}
	if env.IsCronJob {
		log.Printf("skipping because this is a cron job")
		os.Exit(0)
	}
	if env.Branch != "master" && !strings.HasPrefix(env.Tag, "v1.") {
		log.Printf("skipping because branch %q, tag %q is not on the whitelist", env.Branch, env.Tag)
		os.Exit(0)
	}
}

// Debian Packaging

func doDebianSource(cmdline []string) {
	var (
		signer  = flag.String("signer", "", `Signing key name, also used as package author`)
		upload  = flag.String("upload", "", `Where to upload the source package (usually "ppa:klaytn/klaytn")`)
		workdir = flag.String("workdir", "", `Output directory for packages (uses temp dir if unset)`)
		now     = time.Now()
	)
	flag.CommandLine.Parse(cmdline)
	*workdir = makeWorkdir(*workdir)
	env := build.Env()
	maybeSkipArchive(env)

	// Import the signing key.
	if b64key := os.Getenv("PPA_SIGNING_KEY"); b64key != "" {
		key, err := base64.StdEncoding.DecodeString(b64key)
		if err != nil {
			log.Fatal("invalid base64 PPA_SIGNING_KEY")
		}
		gpg := exec.Command("gpg", "--import")
		gpg.Stdin = bytes.NewReader(key)
		build.MustRun(gpg)
	}

	// Create the packages.
	for _, distro := range debDistros {
		meta := newDebMetadata(distro, *signer, env, now)
		pkgdir := stageDebianSource(*workdir, meta)
		debuild := exec.Command("debuild", "-S", "-sa", "-us", "-uc")
		debuild.Dir = pkgdir
		build.MustRun(debuild)

		changes := fmt.Sprintf("%s_%s_source.changes", meta.Name(), meta.VersionString())
		changes = filepath.Join(*workdir, changes)
		if *signer != "" {
			build.MustRunCommand("debsign", changes)
		}
		if *upload != "" {
			build.MustRunCommand("dput", *upload, changes)
		}
	}
}

func makeWorkdir(wdflag string) string {
	var err error
	if wdflag != "" {
		err = os.MkdirAll(wdflag, 0744)
	} else {
		wdflag, err = ioutil.TempDir("", "klay-build-")
	}
	if err != nil {
		log.Fatal(err)
	}
	return wdflag
}

func isUnstableBuild(env build.Environment) bool {
	if env.Tag != "" {
		return false
	}
	return true
}

type debMetadata struct {
	Env build.Environment

	// klaytn version being built. Note that this
	// is not the debian package version. The package version
	// is constructed by VersionString.
	Version string

	Author       string // "name <email>", also selects signing key
	Distro, Time string
	Executables  []debExecutable
}

type debExecutable struct {
	Name, Description string
}

func newDebMetadata(distro, author string, env build.Environment, t time.Time) debMetadata {
	if author == "" {
		// No signing key, use default author.
		author = "Klaytn Builds <infra@groundx.xyz>"
	}
	return debMetadata{
		Env:         env,
		Author:      author,
		Distro:      distro,
		Version:     build.VERSION(),
		Time:        t.Format(time.RFC1123Z),
		Executables: debExecutables,
	}
}

// Name returns the name of the metapackage that depends
// on all executable packages.
func (meta debMetadata) Name() string {
	if isUnstableBuild(meta.Env) {
		return "klaytn-unstable"
	}
	return "klaytn"
}

// VersionString returns the debian version of the packages.
func (meta debMetadata) VersionString() string {
	vsn := meta.Version
	if meta.Env.Buildnum != "" {
		vsn += "+build" + meta.Env.Buildnum
	}
	if meta.Distro != "" {
		vsn += "+" + meta.Distro
	}
	return vsn
}

// ExeList returns the list of all executable packages.
func (meta debMetadata) ExeList() string {
	names := make([]string, len(meta.Executables))
	for i, e := range meta.Executables {
		names[i] = meta.ExeName(e)
	}
	return strings.Join(names, ", ")
}

// ExeName returns the package name of an executable package.
func (meta debMetadata) ExeName(exe debExecutable) string {
	if isUnstableBuild(meta.Env) {
		return exe.Name + "-unstable"
	}
	return exe.Name
}

// ExeConflicts returns the content of the Conflicts field
// for executable packages.
func (meta debMetadata) ExeConflicts(exe debExecutable) string {
	if isUnstableBuild(meta.Env) {
		// Set up the conflicts list so that the *-unstable packages
		// cannot be installed alongside the regular version.
		//
		// https://www.debian.org/doc/debian-policy/ch-relationships.html
		// is very explicit about Conflicts: and says that Breaks: should
		// be preferred and the conflicting files should be handled via
		// alternates. We might do this eventually but using a conflict is
		// easier now.
		return "klaytn, " + exe.Name
	}
	return ""
}

func stageDebianSource(tmpdir string, meta debMetadata) (pkgdir string) {
	pkg := meta.Name() + "-" + meta.VersionString()
	pkgdir = filepath.Join(tmpdir, pkg)
	if err := os.Mkdir(pkgdir, 0755); err != nil {
		log.Fatal(err)
	}

	// Copy the source code.
	build.MustRunCommand("git", "checkout-index", "-a", "--prefix", pkgdir+string(filepath.Separator))

	// Put the debian build files in place.
	debian := filepath.Join(pkgdir, "debian")
	build.Render("build/deb.rules", filepath.Join(debian, "rules"), 0755, meta)
	build.Render("build/deb.changelog", filepath.Join(debian, "changelog"), 0644, meta)
	build.Render("build/deb.control", filepath.Join(debian, "control"), 0644, meta)
	build.Render("build/deb.copyright", filepath.Join(debian, "copyright"), 0644, meta)
	build.RenderString("8\n", filepath.Join(debian, "compat"), 0644, meta)
	build.RenderString("3.0 (native)\n", filepath.Join(debian, "source/format"), 0644, meta)
	for _, exe := range meta.Executables {
		install := filepath.Join(debian, meta.ExeName(exe)+".install")
		docs := filepath.Join(debian, meta.ExeName(exe)+".docs")
		build.Render("build/deb.install", install, 0644, exe)
		build.Render("build/deb.docs", docs, 0644, exe)
	}

	return pkgdir
}

// Windows installer

func doWindowsInstaller(cmdline []string) {
	// Parse the flags and make skip installer generation on PRs
	var (
		arch    = flag.String("arch", runtime.GOARCH, "Architecture for cross build packaging")
		signer  = flag.String("signer", "", `Environment variable holding the signing key (e.g. WINDOWS_SIGNING_KEY)`)
		upload  = flag.String("upload", "", `Destination to upload the archives`)
		workdir = flag.String("workdir", "", `Output directory for packages (uses temp dir if unset)`)
	)
	flag.CommandLine.Parse(cmdline)
	*workdir = makeWorkdir(*workdir)
	env := build.Env()
	maybeSkipArchive(env)

	// Aggregate binaries that are included in the installer
	var (
		devTools []string
		allTools []string
		klayTool string
	)
	for _, file := range allToolsArchiveFiles {
		if file == "COPYING" { // license, copied later
			continue
		}
		allTools = append(allTools, filepath.Base(file))
		if filepath.Base(file) == "klay.exe" {
			klayTool = file
		} else {
			devTools = append(devTools, file)
		}
	}

	// Render NSIS scripts: Installer NSIS contains two installer sections,
	// first section contains the klaytn binary, second section holds the dev tools.
	templateData := map[string]interface{}{
		"License":  "COPYING",
		"Klay":     klayTool,
		"DevTools": devTools,
	}
	build.Render("build/nsis.klay.nsi", filepath.Join(*workdir, "klay.nsi"), 0644, nil)
	build.Render("build/nsis.install.nsh", filepath.Join(*workdir, "install.nsh"), 0644, templateData)
	build.Render("build/nsis.uninstall.nsh", filepath.Join(*workdir, "uninstall.nsh"), 0644, allTools)
	build.Render("build/nsis.pathupdate.nsh", filepath.Join(*workdir, "PathUpdate.nsh"), 0644, nil)
	build.Render("build/nsis.envvarupdate.nsh", filepath.Join(*workdir, "EnvVarUpdate.nsh"), 0644, nil)
	build.CopyFile(filepath.Join(*workdir, "SimpleFC.dll"), "build/nsis.simplefc.dll", 0755)
	build.CopyFile(filepath.Join(*workdir, "COPYING"), "COPYING", 0755)

	// Build the installer. This assumes that all the needed files have been previously
	// built (don't mix building and packaging to keep cross compilation complexity to a
	// minimum).
	version := strings.Split(build.VERSION(), ".")
	if env.Commit != "" {
		version[2] += "-" + env.Commit[:8]
	}
	installer, _ := filepath.Abs("klay-" + archiveBasename(*arch, env) + ".exe")
	build.MustRunCommand("makensis.exe",
		"/DOUTPUTFILE="+installer,
		"/DMAJORVERSION="+version[0],
		"/DMINORVERSION="+version[1],
		"/DBUILDVERSION="+version[2],
		"/DARCH="+*arch,
		filepath.Join(*workdir, "klay.nsi"),
	)

	// Sign and publish installer.
	if err := archiveUpload(installer, *upload, *signer); err != nil {
		log.Fatal(err)
	}
}

// Android archives

func doAndroidArchive(cmdline []string) {
	var (
		local  = flag.Bool("local", false, `Flag whether we're only doing a local build (skip Maven artifacts)`)
		signer = flag.String("signer", "", `Environment variable holding the signing key (e.g. ANDROID_SIGNING_KEY)`)
		deploy = flag.String("deploy", "", `Destination to deploy the archive (usually "https://oss.sonatype.org")`)
		upload = flag.String("upload", "", `Destination to upload the archive`)
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	// Sanity check that the SDK and NDK are installed and set
	if os.Getenv("ANDROID_HOME") == "" {
		log.Fatal("Please ensure ANDROID_HOME points to your Android SDK")
	}
	if os.Getenv("ANDROID_NDK") == "" {
		log.Fatal("Please ensure ANDROID_NDK points to your Android NDK")
	}
	// Build the Android archive and Maven resources
	build.MustRun(goTool("get", "golang.org/x/mobile/cmd/gomobile", "golang.org/x/mobile/cmd/gobind"))
	build.MustRun(gomobileTool("init", "--ndk", os.Getenv("ANDROID_NDK")))
	build.MustRun(gomobileTool("bind", "-ldflags", "-s -w", "--target", "android", "--javapkg", "org.klaytn", "-v", "github.com/klaytn/klaytn/mobile"))

	if *local {
		// If we're building locally, copy bundle to build dir and skip Maven
		os.Rename("klay.aar", filepath.Join(GOBIN, "klay.aar"))
		return
	}
	meta := newMavenMetadata(env)
	build.Render("build/mvn.pom", meta.Package+".pom", 0755, meta)

	// Skip Maven deploy and Azure upload for PR builds
	maybeSkipArchive(env)

	// Sign and upload the archive to Azure
	archive := "klay-" + archiveBasename("android", env) + ".aar"
	os.Rename("klay.aar", archive)

	if err := archiveUpload(archive, *upload, *signer); err != nil {
		log.Fatal(err)
	}
	// Sign and upload all the artifacts to Maven Central
	os.Rename(archive, meta.Package+".aar")
	if *signer != "" && *deploy != "" {
		// Import the signing key into the local GPG instance
		b64key := os.Getenv(*signer)
		key, err := base64.StdEncoding.DecodeString(b64key)
		if err != nil {
			log.Fatalf("invalid base64 %s", *signer)
		}
		gpg := exec.Command("gpg", "--import")
		gpg.Stdin = bytes.NewReader(key)
		build.MustRun(gpg)

		keyID, err := build.PGPKeyID(string(key))
		if err != nil {
			log.Fatal(err)
		}
		// Upload the artifacts to Sonatype and/or Maven Central
		repo := *deploy + "/service/local/staging/deploy/maven2"
		if meta.Develop {
			repo = *deploy + "/content/repositories/snapshots"
		}
		build.MustRunCommand("mvn", "gpg:sign-and-deploy-file", "-e", "-X",
			"-settings=build/mvn.settings", "-Durl="+repo, "-DrepositoryId=ossrh",
			"-Dgpg.keyname="+keyID,
			"-DpomFile="+meta.Package+".pom", "-Dfile="+meta.Package+".aar")
	}
}

func gomobileTool(subcmd string, args ...string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(GOBIN, "gomobile"), subcmd)
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = []string{
		"GOPATH=" + build.GOPATH(),
		"PATH=" + GOBIN + string(os.PathListSeparator) + os.Getenv("PATH"),
	}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "PATH=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	return cmd
}

type mavenMetadata struct {
	Version      string
	Package      string
	Develop      bool
	Contributors []mavenContributor
}

type mavenContributor struct {
	Name  string
	Email string
}

func newMavenMetadata(env build.Environment) mavenMetadata {
	// Collect the list of authors from the repo root
	contribs := []mavenContributor{}
	if authors, err := os.Open("AUTHORS"); err == nil {
		defer authors.Close()

		scanner := bufio.NewScanner(authors)
		for scanner.Scan() {
			// Skip any whitespace from the authors list
			line := strings.TrimSpace(scanner.Text())
			if line == "" || line[0] == '#' {
				continue
			}
			// Split the author and insert as a contributor
			re := regexp.MustCompile("([^<]+) <(.+)>")
			parts := re.FindStringSubmatch(line)
			if len(parts) == 3 {
				contribs = append(contribs, mavenContributor{Name: parts[1], Email: parts[2]})
			}
		}
	}
	// Render the version and package strings
	version := build.VERSION()
	if isUnstableBuild(env) {
		version += "-SNAPSHOT"
	}
	return mavenMetadata{
		Version:      version,
		Package:      "klay-" + version,
		Develop:      isUnstableBuild(env),
		Contributors: contribs,
	}
}

// XCode frameworks

func doXCodeFramework(cmdline []string) {
	var (
		local  = flag.Bool("local", false, `Flag whether we're only doing a local build (skip Maven artifacts)`)
		signer = flag.String("signer", "", `Environment variable holding the signing key (e.g. IOS_SIGNING_KEY)`)
		deploy = flag.String("deploy", "", `Destination to deploy the archive (usually "trunk")`)
		upload = flag.String("upload", "", `Destination to upload the archives`)
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	// Build the iOS XCode framework
	build.MustRun(goTool("get", "golang.org/x/mobile/cmd/gomobile", "golang.org/x/mobile/cmd/gobind"))
	build.MustRun(gomobileTool("init"))
	bind := gomobileTool("bind", "-ldflags", "-s -w", "--target", "ios", "-v", "github.com/klaytn/klaytn/mobile")

	if *local {
		// If we're building locally, use the build folder and stop afterwards
		bind.Dir, _ = filepath.Abs(GOBIN)
		build.MustRun(bind)
		return
	}
	archive := "klay-" + archiveBasename("ios", env)
	if err := os.Mkdir(archive, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	bind.Dir, _ = filepath.Abs(archive)
	build.MustRun(bind)
	build.MustRunCommand("tar", "-zcvf", archive+".tar.gz", archive)

	// Skip CocoaPods deploy and Azure upload for PR builds
	maybeSkipArchive(env)

	// Sign and upload the framework to Azure
	if err := archiveUpload(archive+".tar.gz", *upload, *signer); err != nil {
		log.Fatal(err)
	}
	// Prepare and upload a PodSpec to CocoaPods
	if *deploy != "" {
		meta := newPodMetadata(env, archive)
		build.Render("build/pod.podspec", "Klaytn.podspec", 0755, meta)
		build.MustRunCommand("pod", *deploy, "push", "Klaytn.podspec", "--allow-warnings", "--verbose")
	}
}

type podMetadata struct {
	Version      string
	Commit       string
	Archive      string
	Contributors []podContributor
}

type podContributor struct {
	Name  string
	Email string
}

func newPodMetadata(env build.Environment, archive string) podMetadata {
	// Collect the list of authors from the repo root
	contribs := []podContributor{}
	if authors, err := os.Open("AUTHORS"); err == nil {
		defer authors.Close()

		scanner := bufio.NewScanner(authors)
		for scanner.Scan() {
			// Skip any whitespace from the authors list
			line := strings.TrimSpace(scanner.Text())
			if line == "" || line[0] == '#' {
				continue
			}
			// Split the author and insert as a contributor
			re := regexp.MustCompile("([^<]+) <(.+)>")
			parts := re.FindStringSubmatch(line)
			if len(parts) == 3 {
				contribs = append(contribs, podContributor{Name: parts[1], Email: parts[2]})
			}
		}
	}
	version := build.VERSION()
	if isUnstableBuild(env) {
		version += "-unstable." + env.Buildnum
	}
	return podMetadata{
		Archive:      archive,
		Version:      version,
		Commit:       env.Commit,
		Contributors: contribs,
	}
}

// Cross compilation

func doXgo(cmdline []string) {
	var (
		alltools = flag.Bool("alltools", false, `Flag whether we're building all known tools, or only on in particular`)
	)
	flag.CommandLine.Parse(cmdline)
	env := build.Env()

	subCmd := "get"
	if strings.HasPrefix(runtime.Version(), "go1.18") {
		subCmd = "install"
	}

	// Make sure xgo is available for cross compilation
	gogetxgo := goTool(subCmd, "github.com/klaytn/xgo")
	build.MustRun(gogetxgo)

	// If all tools building is requested, build everything the builder wants
	args := append(buildFlags(env), flag.Args()...)

	if *alltools {
		args = append(args, []string{"--dest", GOBIN}...)
		for _, res := range allToolsArchiveFiles {
			if strings.HasPrefix(res, GOBIN) {
				// Binary tool found, cross build it explicitly
				args = append(args, "./"+filepath.Join("cmd", filepath.Base(res)))
				xgo := xgoTool(args)
				build.MustRun(xgo)
				args = args[:len(args)-1]
			}
		}
		return
	}
	// Otherwise xxecute the explicit cross compilation
	path := args[len(args)-1]
	args = append(args[:len(args)-1], []string{"--dest", GOBIN, path}...)

	xgo := xgoTool(args)
	build.MustRun(xgo)
}

func xgoTool(args []string) *exec.Cmd {
	cmd := exec.Command(filepath.Join(GOBIN, "xgo"), args...)
	cmd.Env = []string{
		"GOPATH=" + build.GOPATH(),
		"GOBIN=" + GOBIN,
	}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	return cmd
}

func installLinter() string {
	lintBin := filepath.Join(build.GOPATH(), "bin", "golangci-lint")

	_, err := exec.LookPath(lintBin)
	if err != nil {
		fmt.Println("Installing golangci-lint.")

		cmdCurl := exec.Command("curl", "-sSfL", "https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh")
		cmdSh := exec.Command("sh", "-s", "--", "-b", filepath.Join(build.GOPATH(), "bin"), "v1.24.0")
		cmdSh.Stdin, err = cmdCurl.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(">>>", strings.Join(cmdCurl.Args, " "))
		if err := cmdCurl.Start(); err != nil {
			log.Fatal(err)
		}

		build.MustRun(cmdSh)

		if err := cmdCurl.Wait(); err != nil {
			log.Fatal(err)
		}
	}
	return lintBin
}

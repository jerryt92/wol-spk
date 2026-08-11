package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var dsmPackageVersionPattern = regexp.MustCompile(`^[0-9]+(?:[._-][0-9]+)*$`)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	packageName, version, err := packageMetadata(filepath.Join(root, "synology", "INFO"))
	if err != nil {
		fatal(err)
	}

	buildDir := filepath.Join(root, "build")
	stageDir := filepath.Join(buildDir, "stage")
	targetDir := filepath.Join(stageDir, "target")
	spkDir := filepath.Join(buildDir, "spk")
	spkPath := filepath.Join(buildDir, fmt.Sprintf("%s-%s-x86_64.spk", packageName, version))

	if err := os.RemoveAll(buildDir); err != nil {
		fatal(err)
	}
	for _, directory := range []string{
		filepath.Join(targetDir, "bin"),
		filepath.Join(targetDir, "ui"),
		spkDir,
	} {
		if err := os.MkdirAll(directory, 0755); err != nil {
			fatal(err)
		}
	}

	fmt.Println("Building Go binary...")
	run(root, nil, "go", "run", "./tools/i18ngen")
	run(root, nil, "go", "run", "./tools/icongen")
	run(root, reproducibleGoEnv(),
		"go", "build", "-trimpath", "-mod=readonly", "-buildvcs=false", "-ldflags=-s -w -buildid=wolmanager -X main.packageVersion="+version, "-o", filepath.Join(targetDir, "bin", "wolmanager"), "./cmd/wolmanager")

	fmt.Println("Preparing DSM package payload...")
	must(copyDir(filepath.Join(root, "synology", "ui"), filepath.Join(targetDir, "ui")))
	run(root, nil, "go", "run", "./tools/tarpack",
		"-source", targetDir,
		"-output", filepath.Join(spkDir, "package.tgz"),
		"-gzip",
		"-executable", "bin/wolmanager",
		"-executable", "ui/index.cgi",
		"-normalize-lf", "ui/config",
		"-normalize-lf", "ui/index.cgi",
		"-normalize-lf", "ui/wolmanager.js",
		"-normalize-lf", "ui/images/icon.svg",
	)
	must(copyFile(filepath.Join(root, "synology", "PACKAGE_ICON.PNG"), filepath.Join(spkDir, "PACKAGE_ICON.PNG")))
	must(copyFile(filepath.Join(root, "synology", "PACKAGE_ICON_256.PNG"), filepath.Join(spkDir, "PACKAGE_ICON_256.PNG")))
	must(copyFile(filepath.Join(root, "synology", "INFO"), filepath.Join(spkDir, "INFO")))
	must(copyFile(filepath.Join(root, "synology", "conf", "privilege"), filepath.Join(spkDir, "conf", "privilege")))
	must(copyDir(filepath.Join(root, "synology", "scripts"), filepath.Join(spkDir, "scripts")))
	run(root, nil, "go", "run", "./tools/tarpack",
		"-source", spkDir,
		"-output", spkPath,
		"-executable", "scripts/postinst",
		"-executable", "scripts/postupgrade",
		"-executable", "scripts/postuninst",
		"-executable", "scripts/preuninst",
		"-executable", "scripts/preupgrade",
		"-executable", "scripts/start-stop-status",
		"-normalize-lf", "INFO",
		"-normalize-lf", "conf/privilege",
		"-normalize-lf", "scripts/postinst",
		"-normalize-lf", "scripts/postupgrade",
		"-normalize-lf", "scripts/postuninst",
		"-normalize-lf", "scripts/preuninst",
		"-normalize-lf", "scripts/preupgrade",
		"-normalize-lf", "scripts/start-stop-status",
	)
	fmt.Printf("Created %s\n", spkPath)
}

func packageMetadata(path string) (string, string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	read := func(name string) string {
		match := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(name) + `="([^"]+)"\r?$`).FindSubmatch(contents)
		if len(match) != 2 {
			return ""
		}
		return string(match[1])
	}
	packageName, version := read("package"), read("version")
	if packageName == "" || version == "" {
		return "", "", fmt.Errorf("failed to read package metadata from %s", path)
	}
	if !dsmPackageVersionPattern.MatchString(version) {
		return "", "", fmt.Errorf("invalid DSM package version %q in %s: use numeric components separated by dots, underscores, or hyphens, for example 1.0.2", version, path)
	}
	return packageName, version, nil
}

func reproducibleGoEnv() map[string]string {
	return map[string]string{
		"CGO_ENABLED":       "0",
		"GOARCH":            "amd64",
		"GOAMD64":           "v1",
		"GOENV":             "off",
		"GOEXPERIMENT":      "",
		"GOFLAGS":           "",
		"GOTOOLCHAIN":       "local",
		"GOWORK":            "off",
		"GOOS":              "linux",
		"SOURCE_DATE_EPOCH": "0",
	}
}

func run(directory string, overrides map[string]string, command string, arguments ...string) {
	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append([]string{}, os.Environ()...)
	for key, value := range overrides {
		prefix := key + "="
		for i := len(cmd.Env) - 1; i >= 0; i-- {
			if len(cmd.Env[i]) >= len(prefix) && cmd.Env[i][:len(prefix)] == prefix {
				cmd.Env = append(cmd.Env[:i], cmd.Env[i+1:]...)
			}
		}
		cmd.Env = append(cmd.Env, prefix+value)
	}
	must(cmd.Run())
}

func copyDir(source, destination string) error {
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return os.MkdirAll(destination, 0755)
		}
		output := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(output, 0755)
		}
		return copyFile(path, output)
	})
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}
	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func must(err error) {
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

package ableton

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charlievieth/fastwalk"
)

const (
	defaultInstallDirWindows = "C:\\ProgramData\\Ableton"
	defaultInstallDirDarwin  = "/Applications"
)

type Installation struct {
	Path string
	Name string
}

type InstallationData struct {
	Path string
	Name string
}

var fwConfig = &fastwalk.Config{}

func FindInstallations() ([]Installation, error) {
	var installations []Installation

	var err error
	switch runtime.GOOS {
	case "windows":
		installDir := defaultInstallDirWindows
		if _, statErr := os.Stat(installDir); os.IsNotExist(statErr) {
			// Directory does not exist, prompt user for input
			println("Default installation directory not found, input the path manually:")
			var userInput string
			_, scanErr := fmt.Scanln(&userInput)
			if scanErr != nil {
				return installations, scanErr
			}
			installDir = strings.TrimSpace(userInput)
			if _, statErr2 := os.Stat(installDir); os.IsNotExist(statErr2) {
				return installations, os.ErrNotExist
			}
		}
		err = fastwalk.Walk(fwConfig,
			installDir, func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() && strings.Contains(info.Name(), "Live") {
					binDir := filepath.Join(path, "/Program")
					instName := info.Name()

					err := filepath.Walk(binDir, func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return nil
						}

						if !info.IsDir() && strings.Contains(info.Name(), "Live") && strings.Contains(info.Name(), ".exe") {
							installations = append(installations, Installation{path, instName})
						}
						return nil
					})
					if err != nil {
						return err
					}

					return fastwalk.SkipDir
				}
				return nil
			})
	case "darwin":
		installDir := defaultInstallDirDarwin
		if _, statErr := os.Stat(installDir); os.IsNotExist(statErr) {
			// Directory does not exist, prompt user for input
			println("Default installation directory not found, input the path manually:")
			var userInput string
			_, scanErr := fmt.Scanln(&userInput)
			if scanErr != nil {
				return installations, scanErr
			}
			installDir = strings.TrimSpace(userInput)
			if _, statErr2 := os.Stat(installDir); os.IsNotExist(statErr2) {
				return installations, os.ErrNotExist
			}
		}
		err = fastwalk.Walk(fwConfig, installDir, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			entryName := info.Name()
			if info.IsDir() && strings.Contains(entryName, "Ableton Live") {
				binPath := filepath.Join(path, "/Contents/MacOS/Live")
				instName := strings.TrimPrefix(entryName, "Ableton ")
				instName = strings.TrimSuffix(instName, ".app")

				if _, err := os.Stat(binPath); err == nil {
					installations = append(installations, Installation{binPath, instName})
				}

				return fastwalk.SkipDir
			}
			return nil
		})
	}

	if err != nil {
		return installations, err
	}

	return installations, err
}

func FindInstallationData() ([]InstallationData, error) {
	var data []InstallationData

	appData, err := os.UserConfigDir()
	if err != nil {
		return data, err
	}
	defaultDataLocation := filepath.Join(appData, "/Ableton")
	err = filepath.Walk(defaultDataLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && strings.Contains(info.Name(), "Live") {
			dataPath := path
			dataName := info.Name()
			err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				if info.IsDir() && strings.Contains(info.Name(), "Unlock") {
					data = append(data, InstallationData{dataPath, dataName})
				}
				return nil
			})
			if err != nil {
				return err
			}

			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return data, err
	}

	return data, err
}

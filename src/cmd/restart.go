package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/khanhas/spicetify-cli/src/utils"
)

func DownloadNoControl() {
	fmt.Println("\n SpotifyNoControl not found, downloading... \n ")
	DownloadAHK()
	fmt.Println("\nTidying up... \n ")
	time.Sleep(utils.INTERVAL)
	TidyUp(spicetifyFolder + `\SpotifyNoControl.ahk`)
	TidyUp(spicetifyFolder + `\SpotifyNoControl.ico`)
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadAHK() {

	mainUrl := "https://raw.githubusercontent.com/SaifAqqad/AHK_SpotifyNoControl/master/"
	dlPath := spicetifyFolder + `\`

	fileUrl := mainUrl + "SpotifyNoControl.ico"
	err := DownloadFile(dlPath+"SpotifyNoControl.ico", fileUrl)
	if err != nil {
		fmt.Println("Error downloading icon file.")
		panic(err)
	}

	fileUrl = mainUrl + "SpotifyNoControl.ahk"
	err = DownloadFile(dlPath+"SpotifyNoControl.ahk", fileUrl)
	if err != nil {
		fmt.Println("Error downloading AHK file.")
		panic(err)
	}

	err = exec.Command("C:\\Program Files\\AutoHotkey\\Compiler\\Ahk2Exe.exe", "/in", dlPath+"SpotifyNoControl.ahk", "/icon", dlPath+"SpotifyNoControl.ico").Run()
	if err != nil {
		fmt.Println("AHK Compile Error.")
		panic(err)
	}
	fmt.Println("Compiled SpotifyNoControl.")
}

func TidyUp(filepath string) error {
	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// RestartSpotify .
func RestartSpotify(flags ...string) {
	launchFlag := settingSection.Key("spotify_launch_flags").Strings("|")
	if len(launchFlag) > 0 {
		flags = append(flags, launchFlag...)
	}

	switch runtime.GOOS {
	case "windows":
		exec.Command("taskkill", "/F", "/IM", "spotify.exe").Run()
		if isAppX {
			ps, _ := exec.LookPath("powershell.exe")
			exe := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "WindowsApps", "Spotify.exe")
			flags = append([]string{"-NoProfile", "-NonInteractive", `& "` + exe + `" --app-directory="` + appDestPath + `"`}, flags...)
			exec.Command(ps, flags...).Start()
		} else {
			exec.Command(filepath.Join(spotifyPath, "spotify.exe"), flags...).Start()
		}

	case "linux":
		exec.Command("pkill", "spotify").Run()
		exec.Command(filepath.Join(spotifyPath, "spotify"), flags...).Start()
	case "darwin":
		exec.Command("pkill", "Spotify").Run()
		flags = append([]string{"-a", "/Applications/Spotify.app"}, flags...)
		exec.Command("open", flags...).Start()
	}
}

// RestartNoControl .
func RestartNoControl(flags ...string) {
	launchFlag := settingSection.Key("spotify_launch_flags").Strings("|")
	if len(launchFlag) > 0 {
		flags = append(flags, launchFlag...)
	}

	switch runtime.GOOS {
	case "windows":
		ncMissing := true
		exec.Command("taskkill", "/F", "/IM", "spotify.exe").Run()
		for ncMissing {
			time.Sleep(utils.INTERVAL)
			ncMissing = false
			if isAppX {
				ps, _ := exec.LookPath("powershell.exe")
				exe := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "WindowsApps", "spotify.exe")
				flags = append([]string{"-NoProfile", "-NonInteractive", `& "` + exe + `" --app-directory="` + appDestPath + `"`}, flags...)
				exec.Command(ps, flags...).Start()

			} else {
				noControlPath := spicetifyFolder + `\SpotifyNoControl.exe`
				err := exec.Command(filepath.Join(noControlPath), flags...).Start()

				if err != nil {
					DownloadNoControl()
					ncMissing = true
					continue
				}
			}
		}

	case "linux":
		exec.Command("pkill", "spotify").Run()
		exec.Command(filepath.Join(spotifyPath, "spotify"), flags...).Start()
	case "darwin":
		exec.Command("pkill", "Spotify").Run()
		flags = append([]string{"-a", "/Applications/Spotify.app"}, flags...)
		exec.Command("open", flags...).Start()
	}
}

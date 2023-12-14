package main

import (
	"os"
	"strings"
)

type StatusbarInfo struct {
	CWD      string
	Hostname string
	User     string
}

func render_statusline() (string, func()) {
	info := get_info()

	statusline_text := ""

	color_positions := []float64{}

	// split path by slash
	folders := strings.Split(info.CWD, "/")

	// for _, folder := range folders {
	//   statusline_text += folder + "/"
	//   color_positions = append(color_positions, float64(len(statusline_text) + (len(folder) / 2)) / 10)
	// }

	for _, folder := range folders {
		statusline_text += folder + " "
		position := float64(len(statusline_text) + (len(folder) / 2))
		normalizedPosition := position / float64(len(info.CWD))
		color_positions = append(color_positions, normalizedPosition)
	}

        bg_length := len(statusline_text)

        renderer := func () {
	  render_gradient(bg_length, color_positions)
        }

        statusline_text = statusline_text[:len(statusline_text)-1] + "  "

        return statusline_text, renderer

}

func get_info() StatusbarInfo {
	return StatusbarInfo{
		CWD:      get_cwd(),
		Hostname: get_hostname(),
		User:     get_user(),
	}
}

func get_cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return cwd
}

func get_hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return hostname
}

func get_user() string {
	return os.Getenv("USER")
}

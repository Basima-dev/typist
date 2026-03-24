package main 

import (
	"fmt"
	"os"
	"os/exec"
     "runtime"
	tea "github.com/charmbracelet/bubbletea"
)

func main (){
	// web flag launch browser ui 
	
	for _, arg := range os.Args[1:]{
		if arg == "--web" || arg == "web" {
			addr, err := startWebServer()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to start web server: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("typist web UI -> %s\n", addr)
			fmt.Printf("press ctrl+c to stop")
			openBrowser(addr)
			//block forever 
			select {}
		}
	}
        // here we set the defualt tui mode 
       p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	   if _, err := p.Run(); err != nil {
		   fmt.Fprintf(os.Stderr, "error: %v\n", err) 
		   os.Exit(1)
	   }
}

func openBrowser(url string) {
	var cmd string 
	var args []string 
	switch runtime.GOOS {
	case "windows":
		cmd = "open"
		args = []string{"/c", "strat", url}
	case "darwin": 
	    cmd = "open"
		args = []string{url}
	default: // Linux (arch btw)
	    cmd = "xdg-open"
		args = []string{url} 
	}
	exec.Command(cmd, args...).Start()
}


package main
  
// Importing fmt, io, and bytes
import (
    "github.com/amrragab8080/tail"
    "strings"
    "os"
    "log"
)
  
// Calling main
func main() {
    t, _ := tail.TailFile("/var/log/dmesg", tail.Config{Follow: true})
    for line := range t.Lines {
        if strings.Contains(line.Text, "NVRM:") {
            f, err := os.OpenFile("/var/log/gpuevent.log",
	       os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
            if err != nil {
	       log.Println(err)
            }
            defer f.Close()
            if _, err := f.WriteString(line.Text + "\n"); err != nil {
	        log.Println(err)
            }
        }
    }
}
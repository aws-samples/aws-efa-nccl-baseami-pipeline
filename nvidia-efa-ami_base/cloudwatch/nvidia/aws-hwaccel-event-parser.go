package main
  
// Importing fmt, io, and bytes
import (
    "fmt"
    "github.com/amrragab8080/tail"
    "strings"
)
  
// Calling main
func main() {
     t, _ := tail.TailFile("/var/log/dmesg", tail.Config{Follow: true})
     for line := range t.Lines {
         if strings.Contains(line.Text, "NVRM:") {
          fmt.Println(line.Text)
       }
     }

}

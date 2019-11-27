// Gets optioncode.md file from https://github.com/timdorr/tesla-api.git and 
// parses it to build out a json file of optioncodes

package main

import (
  "fmt"
  "os/exec"
  "bufio"
  "os"
  "strings"

)

func getrepo(){
  cmd := exec.Command("git","clone","https://github.com/timdorr/tesla-api.git")
  err := cmd.Run()

  if(err != nil){
    fmt.Printf("Unable to get repo: %v\n", err)
    return
  }

}

func getoptionsfile(){
  cmd := exec.Command("cp","tesla-api/docs/vehicle/optioncodes.md", "./")
  err := cmd.Run()

  if(err != nil){
    fmt.Printf("unable to get optioncodes.md: %v\n", err)
    return
  }

}

func rmrepo(){

  cmd := exec.Command("rm","-rf","tesla-api")
  cmd.Run()
}

func main(){

  fmt.Println("Building json file")
  getrepo()
  getoptionsfile()
  rmrepo()

  file, err := os.Open("optioncodes.md")

  outfile, _ := os.Create("file.json")
  outhandle := bufio.NewWriter(outfile)

  defer outfile.Close()

  if( err != nil){
    fmt.Printf("Unable to open optioncodes.md: %s\n", err)
  }

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

  var textLine string
  var code string
  var title string
  //var description string

  var first bool

  first = true

  fmt.Fprintf(outhandle,"{\n")

  for scanner.Scan() {
    textLine = scanner.Text()

    if( strings.HasPrefix(textLine, "|") ){

      if(first){
        first = false
      }else{
        fmt.Fprintf(outhandle,",\n") 
      }

      s := strings.Split(textLine, "|")
   
     code = strings.Trim(strings.TrimSpace(s[1]), "\"")

     title = strings.Replace(strings.TrimSpace(s[2]), "\"", "", -1)

    // description = strings.TrimSpace(s[3])

     fmt.Fprintf(outhandle,"  \"%s\": \"%s\"",code, title)

    } // if has a | in front 

  } // for loop

  fmt.Fprintf(outhandle,"\n")

  fmt.Fprintf(outhandle,"}\n")

  fmt.Fprintf(outhandle,"\n")

  outhandle.Flush()

  file.Close()

}

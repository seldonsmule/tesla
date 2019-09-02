package main

import (
	"os"
	"fmt"
        "bufio"
        "syscall"
        "strings"
        "github.com/seldonsmule/logmsg"
        "golang.org/x/crypto/ssh/terminal"
)

type MyLogin struct {

}


func (et *MyLogin) init(){

  logmsg.Print(logmsg.Info, "In MyLogin init");

}

func (et *MyLogin) prompt(promptString string, hideText bool) (resultString string, worked bool){

  var response1 string
  var response2 string
  var reader *bufio.Reader
  var ok bool

  for ok = true; ok;  {

      if(hideText){
        fmt.Print("Enter ", promptString, ": ")
        bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
        response1 = string(bytePassword)
        fmt.Println() // it's necessary to add a new line after user's input

        fmt.Print("Re-Enter ", promptString, ": ")
        bytePassword, _ = terminal.ReadPassword(int(syscall.Stdin))
        response2 = string(bytePassword)
        fmt.Println() // it's necessary to add a new line after user's input

        if(strings.Compare(response1,response2) != 0){
          fmt.Println(promptString, "did not match\n")
          ok = true
        }else{
          ok = false
        }
      }else{
        reader = bufio.NewReader(os.Stdin)
        fmt.Print("Enter ", promptString, ": ")
        response1, _ = reader.ReadString('\n')
        response1 = strings.TrimSuffix(response1,"\n")
        fmt.Print("Re-Enter ", promptString, ": ")
        response2, _ = reader.ReadString('\n');
        response2 = strings.TrimSuffix(response2,"\n")
      }

      if(strings.Compare(response1,response2) != 0){
        fmt.Println(promptString, "did not match\n")
        ok = true
      }else{
        ok = false
      }
    
  } 

  return response1,true
}

func (et *MyLogin) TerminalLogin(userIDName string, authName string) (userid string, passwd string, worked bool){

  var email string
  var passwd1 string
  var passwd2 string
  var ok bool

  email, ok = et.prompt("Email Address", false)

  for ok = true; ok;  {

      fmt.Print("Enter ", authName, ": ")
      bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
      passwd1 = string(bytePassword)
      fmt.Println() // it's necessary to add a new line after user's input

      fmt.Print("Re-Enter ", authName, ": ")
      bytePassword, _ = terminal.ReadPassword(int(syscall.Stdin))
      passwd2 = string(bytePassword)
      fmt.Println() // it's necessary to add a new line after user's input

      if(strings.Compare(passwd1,passwd2) != 0){
        fmt.Println(authName, "did not match\n")
        ok = true
      }else{
        ok = false
      }
    
  } 

  return email, passwd1, true
}

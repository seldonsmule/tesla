package main

import (
	"os"
	"fmt"
	"flag"
        "strings"
//        "math"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/tesla"
)

func help(){
  

  fmt.Println("tesla_admin - Admin command to setup the tesla.db for use with other commands")
  fmt.Println()
  flag.PrintDefaults()
  fmt.Println()
  fmt.Println("cmds:")
  fmt.Println("     Addsecrets - Add Tesla API secrets - needed for all APIs")
  fmt.Println("     getsecrets - Displays Tesla API secrets - needed for all APIs")
  fmt.Println("     updatesecrets - Updates Tesla API secrets - needed for all APIs")
  fmt.Println("     login - Sets up login info")
  fmt.Println("     refreshtoke - Refreshes token created via login")
  fmt.Println("     getowner - Shows owner details")
  fmt.Println("     delowner - Deletes owner details")
  fmt.Println("     setid - Stores vehicle ID for the cmds.  Requires -vid")
  fmt.Println("     getid - Display vehicle ID for the cmds")
  fmt.Println("     getvehiclelist - Displays a list of vehicles and their IDs owned by the login")
  fmt.Println("     help - Display this help")

}


func main() {

  dirPtr := flag.String("rundir", "./", "Directory to exec from")
  databasePtr := flag.String("dbname", "tesla.db", "Name of database")
  cmdPtr := flag.String("cmd", "help", "Command to run")
  vidPtr := flag.String("vid", "notset", "VehicleId")
  clientidPtr := flag.String("clientid", "notset", "API Client ID")
  clientsecPtr := flag.String("clientsec", "notset", "API Client Secret")

  flag.Parse()


  // if help, the user did not set it
  if(*cmdPtr == "help"){
    help()
    os.Exit(1)
  }

  logName := fmt.Sprintf("%s/tesla_admin.log", *dirPtr)
  dbName := fmt.Sprintf("%s/%s",*dirPtr, *databasePtr)

  logmsg.SetLogFile(logName);

  logmsg.Print(logmsg.Info, "dirPtr = ", *dirPtr)
  logmsg.Print(logmsg.Info, "databasePtr = ", *databasePtr)
  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "vidPtr = ", *vidPtr)
  logmsg.Print(logmsg.Info, "clientidPtr = ", *clientidPtr)
  logmsg.Print(logmsg.Info, "clientsecPtr = ", *clientsecPtr)
  logmsg.Print(logmsg.Info, "tail = ", flag.Args())

  var EgcOptionCodes []string

  myTL := tesla.New(dbName)

  switch *cmdPtr {

    case "login":
      if(!myTL.Login()){
        fmt.Println("Login failed")
        os.Exit(4)
      }

    case "getowner":
      myTL.DumpOwnerInfo()

    case "updatescrets":
      myTL.UpdateSecrets()

    case "setid":
      if(*vidPtr == "notset"){
        fmt.Println("Missing Vehicle ID.  Use -vid\n");
        os.Exit(1);
      }
      myTL.AddVehicleId(*vidPtr)

    case "getid":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
      }
   

    case "delowner":
      myTL.DelOwner()

    case "refreshtoken":
      myTL.RefreshToken(false)

    case "wake":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      if(myTL.WakeCmd(vid)){
        fmt.Println("Wake worked")
      }


    case "getvehiclelist":
      myTL.GetVehicleListCmd()

      count := myTL.VehicleList.GetValueInt("count")

      fmt.Printf("Number of vehicles[%d]\n",count)


      for j:= 0; j < count; j++ {
        fmt.Println("Vehicle: ", j)
        fmt.Printf("Ids[%s]\n", myTL.VehicleList.GetArrayValueString(j,"id_s"))
        fmt.Printf("Vin[%s]\n", myTL.VehicleList.GetArrayValueString(j,"vin"))
        fmt.Printf("DisplayName[%s]\n", myTL.VehicleList.GetArrayValueString(j,"display_name"))
        fmt.Printf("State[%s]\n", myTL.VehicleList.GetArrayValueString(j,"state"))

        //fmt.Printf("options[%s]\n", myTL.VehicleList.GetArrayValueString(j,"option_codes"))
        EgcOptionCodes = strings.Split(myTL.VehicleList.GetArrayValueString(j,"option_codes"),",")

        for i:=0; i < len(EgcOptionCodes); i++ {
          code := EgcOptionCodes[i]
          desc := myTL.ModelXOption(code)
          if(desc != ""){
            fmt.Println(code, " ", desc)
          }
      
        }


      }


    case "getsecrets":

      fmt.Println("GetSecrets")
      dberr := myTL.GetSecrets()
      if(dberr){
        fmt.Println("ID[",myTL.GetClientID(),"] secret[",myTL.GetClientSecret(),"]")
      }else{
        fmt.Println("Yikes - DB error!.  Have you stored secrets?");
        os.Exit(3);
      }
    

    case "addsecrets":
  
      if(*clientidPtr == "notset"){
        fmt.Println("addsecrets -clientid not set\n\n");
        os.Exit(2);
      }
      if(*clientsecPtr == "notset"){
        fmt.Println("addsecrets -clientsec not set\n\n");
        os.Exit(2);
      }

      myTL.SetClientID(*clientidPtr)
      myTL.SetClientSecret(*clientsecPtr)

      myTL.AddSecrets();
     

    case "help":
      help()


    default:
      help()
      os.Exit(2);

  } // end switch


os.Exit(0);


}

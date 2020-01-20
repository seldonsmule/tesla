package main

import (
	"os"
	"fmt"
        "strings"
//        "math"
        "github.com/seldonsmule/logmsg"
        "tesla"
)

func help(){

  fmt.Println("tesla_admin - Admin command to setup the tesla.db for use with other commands")
  fmt.Println()
  fmt.Println("Addsecrets - Add Tesla API secrets - needed for all APIs")
  fmt.Println("getsecrets - Displays Tesla API secrets - needed for all APIs")
  fmt.Println("updatesecrets - Updates Tesla API secrets - needed for all APIs")
  fmt.Println("login - Sets up login info")
  fmt.Println("refreshtoke - Refreshes token created via login")
  fmt.Println("getowner - Shows owner details")
  fmt.Println("delowner - Deletes owner details")
  fmt.Println("setid - Stores vehicle ID for the cmds")
  fmt.Println("getid - Display vehicle ID for the cmds")
  fmt.Println("getvehiclelist - Displays a list of vehicles and their IDs owned by the login")
  fmt.Println("help - Display this help")

}


func main() {

  logmsg.SetLogFile("tesla_admin.log");
  var EgcOptionCodes []string

  myTL := tesla.New("./tesla.db")

  args := os.Args;

  if(len(args) < 2){
    fmt.Println("To few arguments\n");
    help()
    os.Exit(1);
  }


  switch args[1]{

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
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }
      myTL.AddVehicleId(args[2])

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
      if(len(args) < 4){
        fmt.Println("addsecrets missing values\n\n");
        os.Exit(2);
      }

      myTL.SetClientID(args[2])
      myTL.SetClientSecret(args[3])

      myTL.AddSecrets();
     

    case "help":
      help()


    default:
      help()
      os.Exit(2);

  } // end switch


os.Exit(0);


}

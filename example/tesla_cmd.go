package main

import (
	"os"
	"fmt"
        "github.com/seldonsmule/restapi"
        "strings"
        "github.com/seldonsmule/logmsg"
        "tesla"
)


func main() {

  logmsg.SetLogFile("tesla.log");
  var EgcOptionCodes []string

  myTL := tesla.New("./tesla.db")
  //myTL.init()

  args := os.Args;

  if(len(args) < 2){
    fmt.Println("To few arguments\n");
    myTL.Help();
    os.Exit(1);
  }


//  fmt.Println(len(args));

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
   

    case "delowner":
      myTL.DelOwner()

    case "refreshtoken":
      myTL.RefreshToken(false)

    case "wake":
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }

      if(myTL.WakeCmd(args[2])){
        fmt.Println("Wake worked")
      }

    case "setchargelimit":
      if(len(args) < 4){
        fmt.Println("Missing Vehicle ID or charge limit value\n");
        os.Exit(1);
      }

      myTL.WakeCmd(args[2])

      if(myTL.SetChargeLimitCmd(args[2], args[3])){
        fmt.Println("SetChargeLimiit worked")
      }


    case "service_data":
      fallthrough
    case "charge_state":
      fallthrough
    case "climate_state":
      fallthrough
    case "drive_state":
      fallthrough
    case "gui_settings":
      fallthrough
    case "vehicle_state":
      fallthrough
    case "vehicle_config":
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }

      if(myTL.DataRequest(args[2], args[1])){
        r := myTL.DataRequestMap[args[1]]
        r.Dump()
      }


    case "nearbycharging":
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }

      if(myTL.NearbyChargingCmd(args[2])){
        //myTL.nearbyCharging.Dump()

        tmp1 := myTL.NearbyCharging.GetValue("superchargers")
      
        fmt.Printf("superchargers[%s]\n", tmp1)

        myarray := restapi.CastArray(tmp1)

        fmt.Println("array len:", len(myarray))

        for i:=0; i < len(myarray); i++ {
    
          tmpmap := restapi.CastMap(myarray[i])
          
          for name, value := range tmpmap{
      
            fmt.Println(name, "=", value)
    
          }
  
        }

      }


    case "getvehicle":
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }

      myTL.GetVehicleCmd(args[2])

      fmt.Printf("Ids[%s]\n", myTL.SingleVehicle.GetValueString("id_s"))
      fmt.Printf("Vin[%s]\n", myTL.SingleVehicle.GetValueString("vin"))
      fmt.Printf("DisplayName[%s]\n", myTL.SingleVehicle.GetValueString("display_name"))
      fmt.Printf("State[%s]\n", myTL.SingleVehicle.GetValueString("state"))

      EgcOptionCodes = strings.Split(myTL.SingleVehicle.GetValueString("option_codes"),",")

      for i:=0; i < len(EgcOptionCodes); i++ {
        code := EgcOptionCodes[i]
        desc := myTL.ModelXOption(code)
        if(desc != ""){
          fmt.Println(code, " ", desc)
        }
      
      }

/*
      fmt.Println("ID: ", myTL.singleVehicle.Single.Response.Id)
      fmt.Println("VehicleId: ", myTL.singleVehicle.Single.Response.VehicleId)
      fmt.Println("Vin: ", myTL.singleVehicle.Single.Response.Vin)
      fmt.Println("DisplayName: ", myTL.singleVehicle.Single.Response.DisplayName)
      //fmt.Println("OptionCodes: ", myTL.singleVehicle.Single.Response.OptionCodes)
      fmt.Println("Color: ", myTL.singleVehicle.Single.Response.Color)
      fmt.Println("Tokens: ", myTL.singleVehicle.Single.Response.Tokens)
      fmt.Println("State: ", myTL.singleVehicle.Single.Response.State)
      fmt.Println("InService: ", myTL.singleVehicle.Single.Response.InService)
      fmt.Println("IdString: ", myTL.singleVehicle.Single.Response.IdString)
      fmt.Println("CalendarEnabled: ", myTL.singleVehicle.Single.Response.CalendarEnabled)
      fmt.Println("ApiVersion: ", myTL.singleVehicle.Single.Response.ApiVersion)
      fmt.Println("BackseatToken: ", myTL.singleVehicle.Single.Response.BackseatToken)
      fmt.Println("BackseatTokenUpdatedAt: ", myTL.singleVehicle.Single.Response.BackseatTokenUpdatedAt)

        for i:=0; i < len(myTL.singleVehicle.Single.Response.EgcOptionCodes); i++ {
          code := myTL.singleVehicle.Single.Response.EgcOptionCodes[i]
          desc := myTL.ModelXOption(code)
          if(desc != ""){
            fmt.Println(code, " ", desc)
          }
      
        }
*/


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
     

    default:
      myTL.Help();
      os.Exit(2);

  } // end switch


os.Exit(0);


}

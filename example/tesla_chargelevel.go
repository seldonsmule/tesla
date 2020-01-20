package main

import (
	"os"
	"fmt"
        "strconv"
//        "math"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/tesla"
)

func help(){

  fmt.Println("tesla_chargelevel - Simple command to run from cron to reest the charging levels")
  fmt.Println()
  fmt.Println("wake - Wake up vehicle")
  fmt.Println("setchargelimit - Set charge limit for next charge")
  fmt.Println("getchargelimit - Get charge limit for next charge")
  fmt.Println("batterylevel - Get current battery level")
  fmt.Println("homesetto50 - If vehicle is home, set to 50% charge limit")
  fmt.Println("homesetto80 - If vehicle is home, set to 80% charge limit")
  fmt.Println("homecharge - If vehicle is home, use lowlimitlevel and highlimit level to adjust next charge limit")
  fmt.Println("             lowlimitlevel - Level to set to so car does not charge every night")
  fmt.Println("             highlimitlevel - Level to set to so car for next charge IF car is at low (or below) lowlevellimit")
  fmt.Println("             NOTE: If set to 100 already and car is below - will skip adjustments.  If charged, will set to low limit")

}


func main() {

  logmsg.SetLogFile("tesla_chargelvl.log");

  myTL := tesla.New("./tesla.db")

  args := os.Args;

  if(len(args) < 2){
    fmt.Println("To few arguments\n");
    help()
    os.Exit(1);
  }


  switch args[1]{

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

    case "setchargelimit":
      if(len(args) < 3){
        fmt.Println("Missing charge limit value\n");
        os.Exit(1);
      }

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)


      if(myTL.SetChargeLimitCmd(vid, args[2])){
        fmt.Println("SetChargeLimiit worked")
      }

    case "homesetto50":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, "vehicle_state")){
        r := myTL.DataRequestMap["vehicle_state"]
        //r.Dump()

        homelink := r.Obj.GetValue("homelink_nearby")

        if(homelink != true){

          fmt.Println("Telsa not at home - not adjusting anything")
          os.Exit(0)

        }else{
          fmt.Println("Tesla in Garage")
        }

      }

      if(myTL.SetChargeLimitCmd(vid, "50" )){
        fmt.Println("SetChargeLimiit worked")
      }

      if(myTL.DataRequest(vid, "charge_state")){
        r := myTL.DataRequestMap["charge_state"]
        //r.Dump()

        limit := r.Obj.GetValue("charge_limit_soc")

        fmt.Println("Charge limit: ", limit)

      }

    case "homecharge":
    
      if(len(args) < 4){
        fmt.Println("Missing charge high/low\n");
        fmt.Println("Usage: tesla_chargelevel low high");
        fmt.Println("If at or below low range, set to high charging level")
        os.Exit(1);
      }

      lowlvl, _ := strconv.ParseFloat(args[2], 64)
      highlvl, _ := strconv.ParseFloat(args[3], 64)


      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, "vehicle_state")){
        r := myTL.DataRequestMap["vehicle_state"]

        homelink := r.Obj.GetValue("homelink_nearby")

        if(homelink != true){

          fmt.Println("Telsa not at home - not adjusting anything")
          os.Exit(0)

        }else{
          fmt.Println("Tesla in Garage")
        }

      }


      if(myTL.DataRequest(vid, "charge_state")){
        chargestate := myTL.DataRequestMap["charge_state"]

        battery := chargestate.Obj.GetValue("battery_level").(float64)
        limit   := chargestate.Obj.GetValue("charge_limit_soc").(float64)

        if(limit == 100){

          fmt.Println("Battery is set to 100")
          if(battery >= 98){
            fmt.Printf("Looks like we are fully charged[%f], letting lower logic go forth\n", battery)
          }else{
            fmt.Printf("Set to full charge, but battery is at [%f], exiting to let charge\n", battery)
            break;
          }

        } // if 100


/*
fmt.Printf("battery type: %T\n", battery)
fmt.Printf("limit type: %T\n", limit)
fmt.Printf("lowlvl type: %T\n", lowlvl)
fmt.Printf("highlvl type: %T\n", highlvl)
*/


        fmt.Println("Low charge level: ", lowlvl)
        fmt.Println("High charge level: ", highlvl)

        fmt.Println("Charge limit: ", limit)
        fmt.Println("Battery Level: ", battery)

        var newchargelimit string

        if(battery <= lowlvl){
          fmt.Println("hmm - looks like we need to charge.  Setting charge limit to high value")
          newchargelimit = fmt.Sprintf("%.0f", highlvl)
        }else{
          fmt.Println("Battery is above lowlevel, no need to charge.  Setting charge limit to low value")
          newchargelimit = fmt.Sprintf("%.0f", lowlvl)
        }

        fmt.Println("New charge limit to set: ", newchargelimit)

        if(myTL.SetChargeLimitCmd(vid, newchargelimit)){
          fmt.Println("SetChargeLimiit worked")
        }

        if(myTL.DataRequest(vid, "charge_state")){
          r := myTL.DataRequestMap["charge_state"]
          //r.Dump()

          limit := r.Obj.GetValue("charge_limit_soc")

          fmt.Println("Charge limit: ", limit)

        }

      } // end get charge_state
   

    case "homesetto80":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, "vehicle_state")){
        r := myTL.DataRequestMap["vehicle_state"]
        //r.Dump()

        homelink := r.Obj.GetValue("homelink_nearby")

        if(homelink != true){

          fmt.Println("Telsa not at home - not adjusting anything")
          os.Exit(0)

        }else{
          fmt.Println("Tesla in Garage")
        }

      }

      if(myTL.SetChargeLimitCmd(vid, "80" )){
        fmt.Println("SetChargeLimiit worked")
      }

      if(myTL.DataRequest(vid, "charge_state")){
        r := myTL.DataRequestMap["charge_state"]
        //r.Dump()

        limit := r.Obj.GetValue("charge_limit_soc")

        fmt.Println("Charge limit: ", limit)

      }

    case "getchargelimit":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, "charge_state")){
        r := myTL.DataRequestMap["charge_state"]
        //r.Dump()

        limit := r.Obj.GetValue("charge_limit_soc")

        fmt.Println("Charge limit: ", limit)

      }

   case "batterylevel":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, "charge_state")){
        chargestate := myTL.DataRequestMap["charge_state"]
        //r.Dump()

        battery := chargestate.Obj.GetValue("battery_level")
        limit   := chargestate.Obj.GetValue("charge_limit_soc")

        fmt.Println("Charge limit: ", limit)
        fmt.Println("Battery Level: ", battery)

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

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.DataRequest(vid, args[1])){
        r := myTL.DataRequestMap[args[1]]
        r.Dump()
      }


    default:
      help()
      os.Exit(2);

  } // end switch


os.Exit(0);


}

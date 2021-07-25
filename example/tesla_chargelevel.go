package main

import (
	"os"
	"fmt"
	"flag"
        "strconv"
//        "math"
        "github.com/seldonsmule/logmsg"
        "github.com/seldonsmule/tesla"
)

const MINAMP = 24

func help(){

  fmt.Println("tesla_chargelevel - Simple command to run from cron to reest the charging levels")
  fmt.Println()
  flag.PrintDefaults()
  fmt.Println()
  fmt.Println("cmds:")
  fmt.Println("     wake - Wake up vehicle")
  fmt.Println("     setchargelimit - Set charge limit for next charge. Use -limit")
  fmt.Println("     getchargelimit - Get charge limit for next charge")
  fmt.Println("     batterylevel - Get current battery level")
  fmt.Println("     showegarage - Shows the garage test setting (on/off)")
  fmt.Println("     enablegarage - Home test know if car is in garage")
  fmt.Println("     disablegarage - Home test know if car is in garage")
  fmt.Println("     homesetto50 - If vehicle is home, set to 50% charge limit")
  fmt.Println("     homesetto80 - If vehicle is home, set to 80% charge limit")
  fmt.Println("     homecharge - If vehicle is home, use -lowlimit, -highlimit and -minamp level to adjust next charge limit")
  fmt.Println("             -lowlimit - Level to set to so car does not charge every night")
  fmt.Println("             -highlimit - Level to set to so car for next charge IF car is at low (or below) -lowlimit")
  fmt.Println("             -minamp - Min amount of AMPs to use for homecharge.  Default is 24.  If not enough, then homecharge will not execute.  Set to zero to force use with even a 12amp circuit")
  fmt.Println("             NOTE: If set to 100 already and car is below - will skip adjustments.  If charged, will set to low limit")

}


func main() {

  dirPtr := flag.String("rundir", "./", "Directory to exec from")
  databasePtr := flag.String("dbname", "tesla.db", "Name of database")
  cmdPtr := flag.String("cmd", "help", "Command to run")
  limitPtr := flag.String("limit", "notset", "charge limit - used with setchargelimit")
  lowlimitPtr := flag.String("lowlimit", "notset", "charge low limit - used with homecharge")
  highlimitPtr := flag.String("highlimit", "notset", "charge high limit - used with homecharge")
  minampPtr := flag.Uint("minamp", MINAMP, "Min amps needed to use homecharge cmds")



  flag.Parse()

  logName := fmt.Sprintf("%s/tesla_chargelvl.log", *dirPtr)
  dbName := fmt.Sprintf("%s/%s",*dirPtr, *databasePtr)

  logmsg.SetLogFile(logName);


  logmsg.Print(logmsg.Info, "dirPtr = ", *dirPtr)
  logmsg.Print(logmsg.Info, "databasePtr = ", *databasePtr)
  logmsg.Print(logmsg.Info, "cmdPtr = ", *cmdPtr)
  logmsg.Print(logmsg.Info, "limitPtr = ", *limitPtr)
  logmsg.Print(logmsg.Info, "lowlimitPtr = ", *lowlimitPtr)
  logmsg.Print(logmsg.Info, "highlimitPtr = ", *highlimitPtr)
  logmsg.Print(logmsg.Info, "tail = ", flag.Args())


  if(*cmdPtr == "help"){
    help()
    os.Exit(1);
  }

  myTL := tesla.New(dbName)


  switch *cmdPtr {

    case "wake":

      rtn, vid := myTL.GetVehicleId()

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has it been stored yet?")
        fmt.Printf("use tesla_admin -cmd setid -vid -rundir=%s to setup\n\n", *dirPtr)
        os.Exit(1);
      }

      if(myTL.WakeCmd(vid)){
        fmt.Println("Wake worked")
      }

    case "setchargelimit":
      if(*limitPtr == "notset"){
        fmt.Println("Missing charge limit value. Use -limit\n");
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


      if(myTL.SetChargeLimitCmd(vid, *limitPtr)){
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

        // check to see if the database has disabled the homelink
        // test logic - if so we will skip worrying if the car in 
        // at home and do the battery adjustments anyhow

        if(myTL.IsHomeLink()){

          homelink := r.Obj.GetValue("homelink_nearby")

          if(homelink != true){

            fmt.Println("Tesla not at home - not adjusting anything")
            os.Exit(0)

          }else{
            fmt.Println("Tesla in Garage")
          }

        }else{
          fmt.Println("Homelink override set - don't care if at home or not")
          logmsg.Print(logmsg.Info, "Homelink override set - don't care if at home or not")
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
    
      if(*lowlimitPtr == "notset"){
        fmt.Println("Missing charge low limit value. Use -lowlimit\n");
        os.Exit(1);
      }

      if(*highlimitPtr == "notset"){
        fmt.Println("Missing charge high limit value. Use -highlimit\n");
        os.Exit(1);
      }

      lowlvl, _ := strconv.ParseFloat(*lowlimitPtr, 64)
      highlvl, _ := strconv.ParseFloat(*highlimitPtr, 64)


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

        if(myTL.IsHomeLink()){

          homelink := r.Obj.GetValue("homelink_nearby")

          if(homelink != true){

            fmt.Println("Tesla not at home - not adjusting anything")
            os.Exit(0)

          }else{
            fmt.Println("Tesla in Garage")
          }

        }else{
          fmt.Println("Homelink override set - don't care if at home or not")
          logmsg.Print(logmsg.Info, "Homelink override set - don't care if at home or not")
        }

      }


      if(myTL.DataRequest(vid, "charge_state")){
        chargestate := myTL.DataRequestMap["charge_state"]

        battery := chargestate.Obj.GetValue("battery_level").(float64)
        limit   := chargestate.Obj.GetValue("charge_limit_soc").(float64)

        current_max  := chargestate.Obj.GetValue("charge_current_request_max").(float64)
        charger_power  := chargestate.Obj.GetValue("charger_power").(float64)
        charger_voltage  := chargestate.Obj.GetValue("charger_voltage").(float64)
        charger_actual_current  := chargestate.Obj.GetValue("charger_actual_current").(float64)
        scheduled_charging_pending := chargestate.Obj.GetValue("scheduled_charging_pending")

        fmt.Println("Current_Max: ", current_max)
        fmt.Println("charger_power: ", charger_power)
        fmt.Println("charger_voltage: ", charger_voltage)
        fmt.Println("charger_actual_current: ", charger_actual_current)
        fmt.Println("Charging_pending: ", scheduled_charging_pending)

        // using the on/off scheduled pending flag that can be set
        // in the car itself as a determination if we want to change
        // the charge levels.  
        //
        // Meaning if off, then assume the driver knows we have a small
        // circuit (less than 24 amps) and we need all the time we can
        // to charge - so don't set it down to 50%!
        //

/* 06-30-2021 - remoed this logic - it backfired
        if(scheduled_charging_pending != true){
          fmt.Println("Scheduled charging is not set - skipping changing charge level")
          os.Exit(0);
        }else{ 
*/

          // now an extra safety check = are we on a big enough circuit
          // at least 30amps (which gives us 24 usable amps

          if(current_max < float64(*minampPtr)){
 
            fmt.Printf("Max Current (amps) is [%f].  Not enough for a fast charge - skipping changing charge level logic.  Needs to be at least [%d]\n", current_max, *minampPtr)

            os.Exit(0);

          }

/* 06-30-2021 - remoed this logic - it backfired
        } // end of else scheduled_charging_pending
*/



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

        if(myTL.IsHomeLink()){

          homelink := r.Obj.GetValue("homelink_nearby")

          if(homelink != true){

            fmt.Println("Tesla not at home - not adjusting anything")
            os.Exit(0)

          }else{
            fmt.Println("Tesla in Garage")
          }
        }else{
          fmt.Println("Homelink override set - don't care if at home or not")
          logmsg.Print(logmsg.Info, "Homelink override set - don't care if at home or not")
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

   case "enablegarage":
    fmt.Println("enablegarage - Set Homelink logic to ON")
    myTL.SetHomeLinkOn();

   case "disablegarage":
    fmt.Println("enablegarage - Set Homelink logic to OFF")
    myTL.SetHomeLinkOff();

   case "showgarage":
     homestate := myTL.IsHomeLink()

     if(homestate){
       fmt.Println("HomeLink logic set to ON.  App will only change charging levels when at home.  This means if traveling, car will not be woken up to test");
     }else{
       fmt.Println("HomeLink logic set to OFF.  App will charging regardless of where the car is located.  This means car will be woken up everytime to test");
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

      if(myTL.DataRequest(vid, *cmdPtr)){
        r := myTL.DataRequestMap[*cmdPtr]
        r.Dump()
      }


    default:
      help()
      os.Exit(2);

  } // end switch


os.Exit(0);


}

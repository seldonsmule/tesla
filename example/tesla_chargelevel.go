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
  fmt.Println("     stopcharging - Stop charging car")
  fmt.Println("     showegarage - Shows the garage test setting (on/off)")
  fmt.Println("     enablegarage - Home test know if car is in garage - ON")
  fmt.Println("     disablegarage - Home test know if car is in garage OFF")
  fmt.Println("     homesetto50 - If vehicle is home, set to 50% charge limit")
  fmt.Println("     homesetto80 - If vehicle is home, set to 80% charge limit")
  fmt.Println("     homecharge - If vehicle is home, use -lowlimit, -highlimit and -minamp level to adjust next charge limit")
  fmt.Println("             -lowlimit - Level to set to so car does not charge every night")
  fmt.Println("             -highlimit - Level to set to so car for next charge IF car is at low (or below) -lowlimit")
  fmt.Println("             -minamp - Min amount of AMPs to use for homecharge.  Default is 24.  If not enough, then homecharge will not execute.  Set to zero to force use with even a 12amp circuit")
  fmt.Println("             NOTE: If set to 100 already and car is below - will skip adjustments.  If charged, will set to low limit")

}


type ChargeState struct {

  BatteryLevel float64
  ChargeLimitSoc float64
  HomeLinkNearby bool
  ChargeCurrentRequestMax float64
  ChargerPower float64
  ChargerVoltage float64
  ChargerActualCurrent float64
  ScheduledChargingPending bool
  UsableBatteryLevel float64
  EstBatteryRange float64
  TimeToFullCharge float64
}

type VehicleState struct {

  HomelinkNearby bool
  CarVersion string

}


type VehicleData struct {
  cs ChargeState
  vs VehicleState
}


func GetVehicleData(myTL *tesla.MyTesla) (bool) {

   rtn, vid := GetVehicleId(myTL)

   if(rtn){
     fmt.Println("VehicleId: ", vid)
   }else{
     fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
     os.Exit(1);
   }

   myTL.WakeCmd(vid)

   if(!myTL.GetVehicleData(vid)){
     fmt.Println("Error Retrieving VehicleData, check log")
     os.Exit(1);
   }

   //fmt.Println("egc VehicleData: ", myTL.VehicleData)



   /*
   fmt.Println()
   fmt.Println("egc exiting")
   os.Exit(0)	
   */

   //vd = myTL.TeslaVehicleData 

   /*
   if(myTL.DataRequest(vid, "vehicle_data")){
     r := myTL.DataRequestMap["vehicle_data"]
     //r.Dump()

     var charge_state_map map[string]interface{}

     charge_state_map = r.Obj.GetValue("charge_state").(map[string]interface{})

     cs.ChargeLimitSoc = charge_state_map["charge_limit_soc"].(float64)
     cs.BatteryLevel = charge_state_map["battery_level"].(float64)
     cs.ChargeCurrentRequestMax = charge_state_map["charge_current_request_max"].(float64)
     cs.ChargerPower = charge_state_map["charger_power"].(float64)
     cs.ChargerVoltage = charge_state_map["charger_voltage"].(float64)
     cs.ChargerActualCurrent = charge_state_map["charger_actual_current"].(float64)
     cs.ScheduledChargingPending = charge_state_map["scheduled_charging_pending"].(bool)
     cs.UsableBatteryLevel = charge_state_map["usable_battery_level"].(float64)
     cs.EstBatteryRange = charge_state_map["est_battery_range"].(float64)
     cs.TimeToFullCharge = charge_state_map["time_to_full_charge"].(float64)

     var vehicle_state_map map[string]interface{}

     vehicle_state_map = r.Obj.GetValue("vehicle_state").(map[string]interface{})

     vs.HomelinkNearby = vehicle_state_map["homelink_nearby"].(bool)
     vs.CarVersion = vehicle_state_map["car_version"].(string)


   }
   */

  return true
}


func GetVehicleId(myTL *tesla.MyTesla) (bool, string) {

  rtn1, vin := myTL.GetVehicleVin()

  if(!rtn1){
    logmsg.Print(logmsg.Error, "GetVehicleId failed - VIN not found in database: ", vin)
    return false, "VIN not found in database"
  }

  rtn2, vid := myTL.GetVehicleIdFromVinCmd(vin)

  if(!rtn2){
    logmsg.Print(logmsg.Error, "GetVehicleId failed - VIN not found in list from Telsa: ", vin)
    return false, "VIN not found in api response"
  }


  return true, vid

}



func main() {

  //var vd VehicleData

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

  rtn := GetVehicleData(myTL)

  if(!rtn){
    fmt.Println("Error Retrieving Vehicle Data")
    os.Exit(1);
  }

  vd := myTL.VehicleData.Response
  vs := vd.VehicleState
  cs := vd.ChargeState

/*
  fmt.Println("ChargeLimitSoc (limit): ", cs.ChargeLimitSoc)
  fmt.Println("BatteryLevel: ", cs.BatteryLevel)
  fmt.Println("ChargeCurrentRequestMax: ", cs.ChargeCurrentRequestMax)
  fmt.Println("ChargerPower: ", cs.ChargerPower)
  fmt.Println("ChargerVoltage: ", cs.ChargerVoltage)
  fmt.Println("ChargerActualCurrent: ", cs.ChargerActualCurrent)
  fmt.Println("ScheduledChargingPending: ", cs.ScheduledChargingPending)
  fmt.Println("UsableBatteryLevel: ", cs.UsableBatteryLevel)
  fmt.Println("EstBatteryRange: ", cs.EstBatteryRange)

  fmt.Println("HomelinkNearby: ", vs.HomelinkNearby)
  fmt.Println("CarVersion: ", vs.CarVersion)
  */

  /*
  fmt.Println("exiting")
  os.Exit(1);
  */

  switch *cmdPtr {

    case "wake":

      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
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

      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)


      if(myTL.SetChargeLimitCmd(vid, *limitPtr)){
        fmt.Println("SetChargeLimiit worked")
      }

    case "homesetto50":

      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      // check to see if the database has disabled the homelink
      // test logic - if so we will skip worrying if the car in 
      // at home and do the battery adjustments anyhow


      if(myTL.IsHomeLink()){
          if(vs.HomelinkNearby != true){

            fmt.Println("Tesla not at home - not adjusting anything")
            os.Exit(0)

          }else{
            fmt.Println("Tesla in Garage")
          }

      }else{
        fmt.Println("Homelink override set - don't care if at home or not")
        logmsg.Print(logmsg.Info, "Homelink override set - don't care if at home or not")
      }

      if(myTL.SetChargeLimitCmd(vid, "50" )){
        fmt.Println("SetChargeLimiit worked")
      }

      rtn = GetVehicleData(myTL)

      vd := myTL.VehicleData.Response
      cs := vd.ChargeState

      if(!rtn){
        fmt.Println("Error Retrieving Vehicle Data - second time")
        os.Exit(1);
      }

      fmt.Println("Charge limit: ", cs.ChargeLimitSoc)

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


      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
        os.Exit(1);
      }

      myTL.WakeCmd(vid)

      if(myTL.IsHomeLink()){

        if(vs.HomelinkNearby != true){

            fmt.Println("Tesla not at home - not adjusting anything")
            os.Exit(0)

        }else{
          fmt.Println("Tesla in Garage")
        }
      }else{
        fmt.Println("Homelink override set - don't care if at home or not")
        logmsg.Print(logmsg.Info, "Homelink override set - don't care if at home or not")
      }

      fmt.Println("BatteryLevel: ", cs.BatteryLevel)
      fmt.Println("ChargeLimitSoc (limit): ", cs.ChargeLimitSoc)
      fmt.Println("ChargeCurrentRequestMax: ", cs.ChargeCurrentRequestMax)
      fmt.Println("ChargerPower: ", cs.ChargerPower)
      fmt.Println("ChargerVoltage: ", cs.ChargerVoltage)
      fmt.Println("ChargerActualCurrent: ", cs.ChargerActualCurrent)
      fmt.Println("ScheduledChargingPending: ", cs.ScheduledChargingPending)

      // now an extra safety check = are we on a big enough circuit
      // at least 30amps (which gives us 24 usable amps

      if(cs.ChargeCurrentRequestMax < int(*minampPtr)){
 
        fmt.Printf("Max Current (amps) is [%f].  Not enough for a fast charge - skipping changing charge level logic.  Needs to be at least [%d]\n", cs.ChargeCurrentRequestMax, *minampPtr)

        os.Exit(0);

      }

      if(cs.ChargeLimitSoc == 100){

        fmt.Println("Battery is set to 100")
        if(cs.BatteryLevel >= 98){
          fmt.Printf("Looks like we are fully charged[%f], letting lower logic go forth\n", cs.BatteryLevel)
        }else{
          fmt.Printf("Set to full charge, but battery is at [%f], exiting to let charge\n", cs.BatteryLevel)
          break;
        }

      } // if 100


      fmt.Println("Low charge level: ", lowlvl)
      fmt.Println("High charge level: ", highlvl)

      fmt.Println("Charge limit: ", cs.ChargeLimitSoc)
      fmt.Println("Battery Level: ", cs.BatteryLevel)

      var newchargelimit string

      if(cs.BatteryLevel <= int(lowlvl)){
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

      rtn = GetVehicleData(myTL)

      vd := myTL.VehicleData.Response
      cs := vd.ChargeState

      if(!rtn){
        fmt.Println("Error Retrieving Vehicle Data - second time")
        os.Exit(1);
      }

      fmt.Println("Charge limit: ", cs.ChargeLimitSoc)


    case "homesetto80":

      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
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

      rtn = GetVehicleData(myTL)

      vd := myTL.VehicleData.Response
      cs := vd.ChargeState

      if(!rtn){
        fmt.Println("Error Retrieving Vehicle Data - second time")
        os.Exit(1);
      }

      fmt.Println("Charge limit: ", cs.ChargeLimitSoc)

    case "getchargelimit":

      fmt.Println("Charge limit: ", cs.ChargeLimitSoc)

   case "enablegarage":
    fmt.Println("enablegarage - Set Homelink logic to ON")
    myTL.SetHomeLinkOn();

   case "disablegarage":
    fmt.Println("disablegarage - Set Homelink logic to OFF")
    myTL.SetHomeLinkOff();

   case "showgarage":

     homestate := myTL.IsHomeLink()

     if(homestate){
       fmt.Println("HomeLink logic set to ON.  App will only change charging levels when at home.  This means if traveling, car will not be woken up to test");
     }else{
       fmt.Println("HomeLink logic set to OFF.  App will charging regardless of where the car is located.  This means car will be woken up everytime to test");
     }

   case "stopcharging":
     fmt.Println("stopcharging")

     rtn, vid := GetVehicleId(myTL)

     if(rtn){
       fmt.Println("VehicleId: ", vid)
     }else{
       fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
       os.Exit(1);
     }

    myTL.StopChargingCmd(vid);

   case "batterylevel":

     fmt.Println("Battery Level: ", cs.BatteryLevel)
     fmt.Println()

     fmt.Println("ChargeLimitSoc (limit): ", cs.ChargeLimitSoc)
     fmt.Println("ChargeCurrentRequestMax: ", cs.ChargeCurrentRequestMax)
     fmt.Println("ChargerPower: ", cs.ChargerPower)
     fmt.Println("ChargerVoltage: ", cs.ChargerVoltage)
     fmt.Println("ChargerActualCurrent: ", cs.ChargerActualCurrent)
     fmt.Println("ScheduledChargingPending: ", cs.ScheduledChargingPending)
     fmt.Println("UsableBatteryLevel: ", cs.UsableBatteryLevel)
     fmt.Println("EstBatteryRange: ", cs.EstBatteryRange)
     fmt.Println("TimeToFullCharge: ", cs.TimeToFullCharge)

     fmt.Println("HomelinkNearby: ", vs.HomelinkNearby)
     fmt.Println("CarVersion: ", vs.CarVersion)

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

      rtn, vid := GetVehicleId(myTL)

      if(rtn){
        fmt.Println("VehicleId: ", vid)
      }else{
        fmt.Println("Error Retrieving VehicleID, has the VIN been stored yet?")
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

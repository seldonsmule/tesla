package main

import (
	"os"
	"fmt"
        "time"
        "github.com/seldonsmule/restapi"
//        "bufio"
        //"syscall"
        "strconv"
        "strings"
	"net/http"
	"io/ioutil"
        "encoding/json"
//        "database/sql"
  //      "time"
//        _ "github.com/mattn/go-sqlite3" 
        "github.com/seldonsmule/logmsg"
//        "golang.org/x/crypto/ssh/terminal"
)

const const_sec_oneweek    = 604800
const const_sec_twoweek    = 1209600
const const_sec_threeweek  = 1814400

const TESLA_API_URL = "https://owner-api.teslamotors.com"

type rest_cmds struct {

  cmd string
  args string
  desc string
  obj *restapi.Restapi

}

type MyTesla struct {

  myDB *MyDatabase


  clientSecret string
  clientID     string

  accessToken string
  email string
  refreshToken string
  expiresTime int

  //modelxoptions map[string]interface{}
  modelxoptions map[string]string

  vehicleList *restapi.Restapi
  singleVehicle *restapi.Restapi
  wake *restapi.Restapi
  nearbyCharging *restapi.Restapi

  dataRequestMap map[string]rest_cmds

}

func (et *MyTesla) authenticationURL(id string, sec string, email string, pwd string) string{

  url := fmt.Sprintf("%s/oauth/token?grant_type=password&client_id=%s&client_secret=%s&email=%s&password=%s", TESLA_API_URL,
                    id, 
                    sec,
                    email,
                    pwd)

  return url

}

func (et *MyTesla) refreshtokenURL(id string, sec string, refreshtoken string) string{

  url := fmt.Sprintf("%s/oauth/token?grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s", TESLA_API_URL,
                    id, 
                    sec,
                    refreshtoken)

  return url

}

func (et *MyTesla) vehiclesURL() string{

  url := fmt.Sprintf("%s/api/1/vehicles", TESLA_API_URL)

  return url

}

func (et *MyTesla) singlevehicleURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s", TESLA_API_URL, id)

  return url

}

func (et *MyTesla) wakeURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s/wake_up", TESLA_API_URL, id)

  return url

}


func (et *MyTesla) nearbyURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s/nearby_charging_sites", TESLA_API_URL, id)

  return url

}


func (et *MyTesla) ModelXOption(code string) string{
    return string(et.modelxoptions[code])
}


func (et *MyTesla) init(){

  logmsg.Print(logmsg.Info, "In MyTesla init");

  jsonFile, err := os.Open("modelx.json")

  if(err == nil){
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), &et.modelxoptions)

  }else{
    fmt.Println("unable to open file")
  }

  et.myDB = new(MyDatabase)

  et.myDB.init()

  et.vehicleList = nil
  et.dataRequestMap = make(map[string]rest_cmds)
  et.dataRequestMapAdd("charge_state","vehicle_id", "Gets charge state data")
  et.dataRequestMapAdd("climate_state","vehicle_id", "Gets climate state data")
  et.dataRequestMapAdd("drive_state","vehicle_id", "Gets drive state data")
  et.dataRequestMapAdd("gui_settings","vehicle_id", "Gets gui settings data")
  et.dataRequestMapAdd("vehicle_config","vehicle_id", "Gets vehicle config data")
  et.dataRequestMapAdd("service_data","vehicle_id", "Gets service data")


/*
// testing being
  et.Login()

  et.dataRequestMap = make(map[string]rest_cmds)
  et.dataRequestMap["nerd"] = rest_cmds{"vehicle_state", "vehicle_id", nil}


  r := et.dataRequestMap["nerd"]
  r.obj = restapi.NewGet(r.cmd, et.data_requestURL(r.cmd, "27207623174674431")) 
  r.obj.SetBearerAccessToken(et.accessToken)
  r.obj.HasInnerMap("response")

  r.obj.Send()

  r.obj.Dump()

  et.dataRequestMap["nerd"] = r
  

fmt.Println("r print:", r)
fmt.Println("last print:", et.dataRequestMap)
fmt.Println("len:", len(et.dataRequestMap))
os.Exit(5)
// testing end
*/

}

func (et *MyTesla) dataRequestMapAdd(name string, args string, desc string){


  et.dataRequestMap[name] = rest_cmds{name, args, desc, nil}
}

func (et *MyTesla) AddSecrets(){

   et.myDB.AddApiDetails(et.clientID, et.clientSecret);

}

func (et *MyTesla) GetSecrets() (bool){

  if(!et.myDB.GetApiDetails(&et.clientID, &et.clientSecret)){
    // not found - try auto updating them
    if(!et.UpdateSecrets()){
      return false
    }
  }

  return true

}

func (et *MyTesla) AddOwner(){

   et.myDB.AddOwner(et.email, et.accessToken, et.refreshToken, et.expiresTime);

}

func (et *MyTesla) GetOwner() (bool){

  return et.myDB.GetOwner(&et.email, &et.accessToken,
                                &et.refreshToken, &et.expiresTime)

}

func (et *MyTesla) DelOwner() (bool){

  return et.myDB.DelOwner()

}

func (et *MyTesla) RefreshToken(skipLogin bool) bool{

  if(!skipLogin){ // if skipping we already have the owner info
    et.Login() // the act of logging in will populate this info
  }

  fmt.Println("Starting RefreshToken")

  r := restapi.NewPost("authentication", et.refreshtokenURL(et.clientID,
                                                          et.clientSecret,
                                                          et.refreshToken))
  if(r.Send()){
    //r.Dump()
  }else{
    logmsg.Print(logmsg.Error,"refresh Failed")
    return false
  }
                                            
  et.accessToken = r.GetValueString("access_token")
  et.refreshToken = r.GetValueString("refresh_token")
  created := r.GetValue("created_at")
  expires := r.GetValue("expires_in")
  et.expiresTime = r.CastFloatToInt(created) + r.CastFloatToInt(expires)
  et.AddOwner()

  return true

}

func (et *MyTesla) Wake(id string) bool{

  var stateStr string

  et.Login() // the act of logging in will populate this info

  if(et.wake != nil){
    return true  // i.e., we already made this call
  }

  et.wake = restapi.NewPost("wake", et.wakeURL(id))

  et.wake.SetBearerAccessToken(et.accessToken)
  et.wake.HasInnerMap("response")

  for et.wake.Send() {

    stateStr = et.wake.GetValueString("state")

    if(strings.Compare(stateStr,"online") == 0){
      et.wake.Dump()
      return true
    }

    logmsg.Print(logmsg.Warning,"Wake still waiting. State:", stateStr)

  }

  logmsg.Print(logmsg.Error,"wake failed")
  return false

}

func (et *MyTesla) DataRequest(id string, cmd string) bool{

  et.Login() // the act of logging in will populate this info

  r := et.dataRequestMap[cmd]

  if(r.obj == nil){ // not setup before
    url := fmt.Sprintf("%s/api/1/vehicles/%s/data_request/%s", TESLA_API_URL, id, cmd)
    r.obj = restapi.NewGet(r.cmd, url) 
    r.obj.SetBearerAccessToken(et.accessToken)
    r.obj.HasInnerMap("response")
  }

  if(r.obj.Send()){
    //et.tmpobj.Dump()
  }else{
    logmsg.Print(logmsg.Error,cmd," Failed ", id)
    return false
  }

  et.dataRequestMap[cmd] = r

  return true
}

func (et *MyTesla) NearbyCharging(id string) bool{

  et.Login() // the act of logging in will populate this info

  et.nearbyCharging = restapi.NewGet("nearby_charging_sites", et.nearbyURL(id))

  et.nearbyCharging.SetBearerAccessToken(et.accessToken)
  et.nearbyCharging.HasInnerMap("response")

  if(et.nearbyCharging.Send()){
    //et.nearbyCharging.Dump()
  }else{
    logmsg.Print(logmsg.Error,"NearbyCharging Failed ", id)
    return false
  }

  return true

}

func (et *MyTesla) GetVehicle(id string) bool{


  et.Login() // the act of logging in will populate this info

  if(et.singleVehicle != nil){
    return true  // i.e., we already made this call
  }

  et.singleVehicle = restapi.NewGet("singlevehicles", et.singlevehicleURL(id))

  et.singleVehicle.SetBearerAccessToken(et.accessToken)
  et.singleVehicle.HasInnerMap("response")

  if(et.singleVehicle.Send()){
    et.singleVehicle.Dump()
  }else{
    logmsg.Print(logmsg.Error,"get vehicle failed:", id)
    return false
  }

/*
  //vehicle := new(RestVehicles)

  if(!et.singleVehicle.sendGetSingleVehicle(id, et.accessToken)){
    fmt.Println("GET Vehicle failed");
    logmsg.Print(logmsg.Error,"GET Vehicle failed")
    return false
  }
*/


  return true

}

func (et *MyTesla) GetVehicleList() bool{

  et.Login() // the act of logging in will populate this info

  //vehicles := new(RestVehicles)

  if(et.vehicleList != nil){
    return true  // i.e., we already made this call
  }

  et.vehicleList = restapi.NewGet("vehicles", et.vehiclesURL())

  et.vehicleList.SetBearerAccessToken(et.accessToken)
  et.vehicleList.HasInnerMapArray("response","count")

  if(et.vehicleList.Send()){
    //et.vehicleList.Dump()
  }else{
    logmsg.Print(logmsg.Error,"get vehicles list failed")
    return false
  }

/*
  if(!et.vehicleList.sendGetVehicleList(et.accessToken)){
    fmt.Println("GET Vehicles failed");
    logmsg.Print(logmsg.Error,"GET Vehicles failed")
    return false
  }
*/


  return true

}

func (et *MyTesla) DumpOwnerInfo(){

  et.Login() // the act of logging in will populate this info

  expiresTime := time.Unix(int64(et.expiresTime), 0)

  fmt.Println("Owner table dump")
  fmt.Println("Email: ", et.email)
  fmt.Println("AccessToken: ", et.accessToken)
  fmt.Println("RefreshToken: ", et.refreshToken)
  fmt.Println("ExpiresTime: ", strconv.Itoa(et.expiresTime))
  fmt.Println("ExpiresTime: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               expiresTime.Month(),expiresTime.Day(), expiresTime.Year(),
             expiresTime.Hour(), expiresTime.Minute(), expiresTime.Second() ) )

  

}


func (et *MyTesla) Login() bool{

  var email1 string
  //var email2 string
  var passwd1 string
  //var passwd2 string
  //var reader *bufio.Reader
  var ok bool

  dberr := et.GetSecrets()
  if(!dberr){
    fmt.Println("Yikes - DB error!.  Have you stored secrets?");
    os.Exit(4);
  }


  // see if we already have an access token

  status := et.GetOwner()

fmt.Println("owner status:", status)

  if(status){

    // check to see if we need to refresh - if within 1 week, refresh
    timeNow := time.Now()

    unixTime := timeNow.Unix()

/* for testing purposes only
    unixTime = int64(et.expiresTime -  const_sec_twoweek) // two weeks
    //unixTime = int64(et.expiresTime -  const_sec_threeweek) // three weeks
    //unixTime = int64(et.expiresTime -  const_sec_oneweek) // one weeks
    timeNow = time.Unix(unixTime, 0)
*/ // end for testing purposes only

/*
    expiresTime := time.Unix(int64(et.expiresTime), 0)
    fmt.Printf("UnixTime[%d] vs expiresTime[%d]\n", unixTime, et.expiresTime)

    fmt.Println("TimeNow: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               timeNow.Month(),timeNow.Day(), timeNow.Year(),
             timeNow.Hour(), timeNow.Minute(), timeNow.Second() ) )

    fmt.Println("ExpiresTime: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               expiresTime.Month(),expiresTime.Day(), expiresTime.Year(),
             expiresTime.Hour(), expiresTime.Minute(), expiresTime.Second() ) )
*/

    if(int(unixTime) >= (et.expiresTime - const_sec_twoweek)){
      fmt.Println("Time to refresh")
      status = et.RefreshToken(true)  
    }

    if(status){
      fmt.Println("Logined to [", et.email, "]")
      return(true)
    }

  } // if status == true  Meaning we have data


  logmsg.Print(logmsg.Warning,"AccessToken not found")
  fmt.Println("AccessToken not found, need to authenticate")

  term := new(MyLogin)
  email1, passwd1, ok = term.TerminalLogin("Email Address", "Password")

  if(!ok){
    fmt.Println("Failure to get credentials")
    return(false)
  }


  r := restapi.NewPost("authentication", et.authenticationURL(et.clientID,
                                                          et.clientSecret,
                                                          email1,
                                                          passwd1))
  r.DebugOn()

  if(r.Send()){
  //  r.Dump()
  }else{
    logmsg.Print(logmsg.Error,"authentication failed")
    return false
  }


  et.email = email1
  et.accessToken = r.GetValueString("access_token")
  et.refreshToken = r.GetValueString("refresh_token")
  created := r.GetValue("created_at")
  expires := r.GetValue("expires_in")
  et.expiresTime = r.CastFloatToInt(created) + r.CastFloatToInt(expires)
  et.AddOwner()

  return true
}

//////////////////////////

func (et *MyTesla) help(){

  fmt.Println("tesla login | getvehiclelist | addsecrets | getsecrets | getowner | delowner | refreshtoken | getvehiclelist | getchargestate | getclimatestate | updatesecrets \n");

  fmt.Println()
  fmt.Println("Get State Cmds")
  for name, value := range et.dataRequestMap {

    fmt.Println(name, value.args, ":", value.desc)

  }

/*
  fmt.Println("tesla login | getvehiclelist | addsecrets | getsecrets | getowner | delowner | refreshtoken | getvehiclelist | getchargestate | getclimatestate | updatesecrets \n");
  fmt.Println("login username password  - Connects to the car\n");
  fmt.Println("getvehiclelist  - get a list of capabilities\n");
  fmt.Println("addsecrets client_id client_secret - stores these\n");
  fmt.Println("getsecrets - shows them\n");
  fmt.Println("getowner - shows details\n");
  fmt.Println("delowner - deletes the owner auth details\n");
  fmt.Println("refreshtoken - Proactively refreshes the accestoken\n");
  fmt.Println("getvehiclelist - Return list of vehicles of the owner\n");
  fmt.Println("getvehicle vehicle_id - Return details of a single vechile by id\n");
  fmt.Println("getchargestate vehicle_id - Get Vehicle charge state\n");
  fmt.Println("getclimatestate vehicle_id - Get Vehicle climate state\n");
  fmt.Println("getdrivestate vehicle_id - Get Vehicle drive state\n");
  fmt.Println("getguisettings vehicle_id - Get Vehicle gui settings\n");
  fmt.Println("getvehiclestate vehicle_id - Get Vehicle vehicle state\n");
  fmt.Println("getvehicleconfig vehicle_id - Get Vehicle vehicle config\n");
  fmt.Println("nearbycharging vehicle_id - Get Nearby Charging data\n");
  fmt.Println("updatescrets - Refresh Secrets from PasteBin URL: https://pastebin.com/raw/pS7Z6yyP\n");
  fmt.Println("wake vehicle_id - Wake up Vehicle\n");
*/


}

func (et *MyTesla) charCleanUp(original[] byte) string{

  var buffer[5000] byte
  var char rune
  var rtnstr string

  index := 0

  for j:=0; j < len(original); j++ {

    char = rune(original[j])

    if(strconv.IsPrint(char)){
      buffer[index] = original[j]
//fmt.Printf("%c\n", original[j])
      index++
    }else{
      if(original[j] == '\n'){
 //       fmt.Println("Found new line")
          buffer[index] = original[j]
          index++
      }else{
//        fmt.Printf("element[%d] not printable value[%d] \n", j, char)
      }
    }

  }

  buffer[index] = '\n'
  index++

//  fmt.Printf("Index is now [%d]\n", index)

  rtnstr = string(buffer[:index])

//  fmt.Printf(rtnstr)

  return rtnstr

}

func (et *MyTesla) UpdateSecrets() bool{

  var client_id string
  var client_secret string

  client_id = "blank"
  client_secret = "blank"
  
  url := "https://pastebin.com/raw/pS7Z6yyP"

  req, _ := http.NewRequest("GET", url, nil)

  res, _ := http.DefaultClient.Do(req)

  defer res.Body.Close()
  body, _ := ioutil.ReadAll(res.Body)

//  tmpstr := string(body)
  tmpstr := et.charCleanUp(body)

  //fmt.Println(tmpstr)

  tmparray := strings.Split(tmpstr, "\n")

  for i:= 0; i < len(tmparray); i++ {
//    fmt.Println(tmparray[i])
    item := strings.Split(tmparray[i], "=")
    if(strings.Compare(item[0], "TESLA_CLIENT_ID") == 0){
      client_id = item[1]
//      fmt.Println("client id found: ", client_id)
    }
    if(strings.Compare(item[0], "TESLA_CLIENT_SECRET") == 0){
      client_secret = item[1]
      //fmt.Println("client secret found: ", client_secret)
    }
  }

  if ( (strings.Compare(client_id, "blank") == 0) || (strings.Compare(client_secret, "blank") == 0) ){

    fmt.Println("Error getting client ID and or secret.  URL returned:")
    fmt.Println(tmpstr)
    return false
  }

  et.clientID = client_id
  et.clientSecret = client_secret

  et.AddSecrets();

//  fmt.Println("client id found: [", client_id, "]")
//  fmt.Println("client secret found: [", client_secret, "]")

  return true

}

func main() {

  logmsg.SetLogFile("tesla.log");
  var EgcOptionCodes []string

  myTL := new(MyTesla)
  myTL.init()

  args := os.Args;

  if(len(args) < 2){
    fmt.Println("To few arguments\n");
    myTL.help();
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


      if(myTL.Wake(args[2])){
        fmt.Println("Wake worked")
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
        r := myTL.dataRequestMap[args[1]]
        r.obj.Dump()
      }

    case "nearbycharging":
      if(len(args) < 3){
        fmt.Println("Missing Vehicle ID\n");
        os.Exit(1);
      }

      if(myTL.NearbyCharging(args[2])){
        myTL.nearbyCharging.Dump()

        tmp1 := myTL.nearbyCharging.GetValue("superchargers")
      
        fmt.Printf("superchargers[%s]\n", tmp1)

        myarray := myTL.nearbyCharging.CastArray(tmp1)

        fmt.Println("array len:", len(myarray))

        for i:=0; i < len(myarray); i++ {
    
          tmpmap := myTL.nearbyCharging.CastMap(myarray[i])
          
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

      myTL.GetVehicle(args[2])

      fmt.Printf("Ids[%s]\n", myTL.singleVehicle.GetValueString("id_s"))
      fmt.Printf("Vin[%s]\n", myTL.singleVehicle.GetValueString("vin"))
      fmt.Printf("DisplayName[%s]\n", myTL.singleVehicle.GetValueString("display_name"))
      fmt.Printf("State[%s]\n", myTL.singleVehicle.GetValueString("state"))

      EgcOptionCodes = strings.Split(myTL.singleVehicle.GetValueString("option_codes"),",")

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
      myTL.GetVehicleList()

      count := myTL.vehicleList.GetValueInt("count")

      fmt.Printf("Number of vehicles[%d]\n",count)

//myTL.vehicleList.Dump()
   

      for j:= 0; j < count; j++ {
        fmt.Println("Vehicle: ", j)
        fmt.Printf("Ids[%s]\n", myTL.vehicleList.GetArrayValueString(j,"id_s"))
        fmt.Printf("Vin[%s]\n", myTL.vehicleList.GetArrayValueString(j,"vin"))
        fmt.Printf("DisplayName[%s]\n", myTL.vehicleList.GetArrayValueString(j,"display_name"))
        fmt.Printf("State[%s]\n", myTL.vehicleList.GetArrayValueString(j,"state"))

        //fmt.Printf("options[%s]\n", myTL.vehicleList.GetArrayValueString(j,"option_codes"))
        EgcOptionCodes = strings.Split(myTL.vehicleList.GetArrayValueString(j,"option_codes"),",")

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
        fmt.Println("ID[",myTL.clientID,"] secret[",myTL.clientSecret,"]")
      }else{
        fmt.Println("Yikes - DB error!.  Have you stored secrets?");
        os.Exit(3);
      }
    

    case "addsecrets":
      if(len(args) < 4){
        fmt.Println("addsecrets missing values\n\n");
        os.Exit(2);
      }

      myTL.clientID = args[2]
      myTL.clientSecret = args[3]

      myTL.AddSecrets();
     

    default:
      myTL.help();
      os.Exit(2);

  } // end switch


os.Exit(0);


}

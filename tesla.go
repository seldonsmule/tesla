package tesla

import (
	"os"
	"fmt"
        "time"
        "github.com/seldonsmule/restapi"
        "strconv"
        "strings"
	"net/http"
	"io/ioutil"
        "encoding/json"
        "github.com/seldonsmule/logmsg"
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

func (rc *rest_cmds) Dump(){

  fmt.Printf("Cmd: %s\n", rc.cmd)
  fmt.Printf("Args: %s\n", rc.args)
  fmt.Printf("Desc: %s\n", rc.desc)
  rc.obj.Dump()

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
  Modelxoptions map[string]string

  VehicleList *restapi.Restapi
  SingleVehicle *restapi.Restapi
  Wake *restapi.Restapi
  Setchargelimit *restapi.Restapi
  NearbyCharging *restapi.Restapi

  DataRequestMap map[string]rest_cmds

}

func (et *MyTesla) GetClientID() string{
  return et.clientID
}

func (et *MyTesla) SetClientID(id string){
  et.clientID = id
}

func (et *MyTesla) GetClientSecret() string{
  return et.clientSecret
}

func (et *MyTesla) SetClientSecret(secret string){
  et.clientSecret = secret
}

func (et *MyTesla) authenticationURL(id string, sec string, email string, pwd string) string{

  url := fmt.Sprintf("%s/oauth/token?grant_type=password&client_id=%s&client_secret=%s&email=%s&password=%s", TESLA_API_URL,
                    id, 
                    sec,
                    email,
                    pwd)

  return url

}

func (et *MyTesla) setchargelimitURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s/command/set_charge_limit", TESLA_API_URL, id)

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
    return string(et.Modelxoptions[code])
}


func New() *MyTesla{

  t := new(MyTesla) 

  logmsg.Print(logmsg.Info, "In MyTesla New");

  jsonFile, err := os.Open("modelx.json")

  if(err == nil){
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), &t.Modelxoptions)

  }else{
    fmt.Println("unable to open file")
  }

  t.myDB = new(MyDatabase)

  t.myDB.init()

  t.VehicleList = nil
  t.DataRequestMap = make(map[string]rest_cmds)
  t.dataRequestMapAdd("charge_state","vehicle_id", "Gts charge state data")
  t.dataRequestMapAdd("climate_state","vehicle_id", "Gets climate state data")
  t.dataRequestMapAdd("drive_state","vehicle_id", "Gets drive state data")
  t.dataRequestMapAdd("gui_settings","vehicle_id", "Gets gui settings data")
  t.dataRequestMapAdd("vehicle_config","vehicle_id", "Gets vehicle config data")
  t.dataRequestMapAdd("vehicle_state","vehicle_id", "Gets vehicle state data")
  t.dataRequestMapAdd("nearbycharging","vehicle_id", "Gets vehicle state data")
  t.dataRequestMapAdd("service_data","vehicle_id", "Gets service data")


  return t

}

func (et *MyTesla) dataRequestMapAdd(name string, args string, desc string){


  et.DataRequestMap[name] = rest_cmds{name, args, desc, nil}
}

func (et *MyTesla) AddSecrets(){

   et.myDB.AddApiDetails(et.GetClientID(), et.GetClientSecret());

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
  et.expiresTime = restapi.CastFloatToInt(created) + restapi.CastFloatToInt(expires)
  et.AddOwner()

  return true

}

func (et *MyTesla) WakeCmd(id string) bool{

  var stateStr string

  et.Login() // the act of logging in will populate this info

  if(et.Wake != nil){
    return true  // i.e., we already made this call
  }

  et.Wake = restapi.NewPost("wake", et.wakeURL(id))

  et.Wake.SetBearerAccessToken(et.accessToken)
  et.Wake.HasInnerMap("response")

  for et.Wake.Send() {

    stateStr = et.Wake.GetValueString("state")

    if(strings.Compare(stateStr,"online") == 0){
      et.Wake.Dump()
      return true
    }

    logmsg.Print(logmsg.Warning,"Wake still waiting. State:", stateStr)

  }

  logmsg.Print(logmsg.Error,"wake failed")
  return false

}

func (et *MyTesla) SetChargeLimitCmd(id string, percent_value string) bool{

  //var stateStr string

  et.Login() // the act of logging in will populate this info

  et.Setchargelimit = restapi.NewPost("Setchargelimit", et.setchargelimitURL(id))

  et.Setchargelimit.SetBearerAccessToken(et.accessToken)
  et.Setchargelimit.HasInnerMap("response")

  jsonstr := fmt.Sprintf("{\"percent\":\"%s\"}", percent_value)

  et.Setchargelimit.SetPostJson(jsonstr)

  if(et.Setchargelimit.Send()){
    et.Setchargelimit.Dump()
  }else{
    logmsg.Print(logmsg.Error,"Setchargelimit"," Failed ", id)
    return false
  }

/*
  for et.wake.Send() {

    stateStr = et.wake.GetValueString("state")

    if(strings.Compare(stateStr,"online") == 0){
      et.wake.Dump()
      return true
    }

    logmsg.Print(logmsg.Warning,"Wake still waiting. State:", stateStr)

  }
*/

  return true

}

func (et *MyTesla) DataRequest(id string, cmd string) bool{

  et.Login() // the act of logging in will populate this info

  r := et.DataRequestMap[cmd]

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

  et.DataRequestMap[cmd] = r

  return true
}

func (et *MyTesla) NearbyChargingCmd(id string) bool{

  et.Login() // the act of logging in will populate this info

  et.NearbyCharging = restapi.NewGet("nearby_charging_sites", et.nearbyURL(id))

  et.NearbyCharging.SetBearerAccessToken(et.accessToken)
  et.NearbyCharging.HasInnerMap("response")

  if(et.NearbyCharging.Send()){
    //et.NearbyCharging.Dump()
  }else{
    logmsg.Print(logmsg.Error,"NearbyCharging Failed ", id)
    return false
  }

  return true

}

func (et *MyTesla) GetVehicleCmd(id string) bool{


  et.Login() // the act of logging in will populate this info

  if(et.SingleVehicle != nil){
    return true  // i.e., we already made this call
  }

  et.SingleVehicle = restapi.NewGet("Singlevehicles", et.singlevehicleURL(id))

  et.SingleVehicle.SetBearerAccessToken(et.accessToken)
  et.SingleVehicle.HasInnerMap("response")

  if(et.SingleVehicle.Send()){
    et.SingleVehicle.Dump()
  }else{
    logmsg.Print(logmsg.Error,"get vehicle failed:", id)
    return false
  }

/*
  //vehicle := new(RestVehicles)

  if(!et.SingleVehicle.sendGetSingleVehicle(id, et.accessToken)){
    fmt.Println("GET Vehicle failed");
    logmsg.Print(logmsg.Error,"GET Vehicle failed")
    return false
  }
*/


  return true

}

func (et *MyTesla) GetVehicleListCmd() bool{

  et.Login() // the act of logging in will populate this info

  //vehicles := new(RestVehicles)

  if(et.VehicleList != nil){
    return true  // i.e., we already made this call
  }

  et.VehicleList = restapi.NewGet("vehicles", et.vehiclesURL())

  et.VehicleList.SetBearerAccessToken(et.accessToken)
  et.VehicleList.HasInnerMapArray("response","count")

  if(et.VehicleList.Send()){
    //et.VehicleList.Dump()
  }else{
    logmsg.Print(logmsg.Error,"get vehicles list failed")
    return false
  }

/*
  if(!et.VehicleList.sendGetVehicleList(et.accessToken)){
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
  et.expiresTime = restapi.CastFloatToInt(created) + restapi.CastFloatToInt(expires)
  et.AddOwner()

  return true
}

//////////////////////////

func (et *MyTesla) Help(){

  fmt.Println("tesla login | getvehiclelist | addsecrets | getsecrets | getowner | delowner | refreshtoken | getvehiclelist | getchargestate | getclimatestate | updatesecrets \n");

  fmt.Println()
  fmt.Println("Get State Cmds")
  for name, value := range et.DataRequestMap {

    fmt.Println(name, value.args, ":", value.desc)

  }
  fmt.Println();
  fmt.Println("Set Cmds");
  fmt.Println("wake vehicle_id - Wake up Vehicle");
  fmt.Println("setchargelimit vehicle_id percent - Sets the charge limit\n");

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

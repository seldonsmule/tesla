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
        "io"
        "crypto/rand"
        "crypto/sha256"
        "encoding/base64"
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
  Obj *restapi.Restapi

}

func (rc *rest_cmds) Dump(){

  fmt.Printf("Cmd: %s\n", rc.cmd)
  fmt.Printf("Args: %s\n", rc.args)
  fmt.Printf("Desc: %s\n", rc.desc)
  rc.Obj.Dump()

}

type MyTesla struct {

  myDB *MyDatabase


  clientSecret string
  clientID     string

  // for new tesla sso auth - begin
  
  SSO_state        string
  SSO_challenge    string
  SSO_challengeSum string

  // for new tesla sso auth - end

  accessToken string
  email string
  refreshToken string
  expiresTime int

  debug bool

  //modelxoptions map[string]interface{}
  Modelxoptions map[string]string

  VehicleList *restapi.Restapi
  SingleVehicle *restapi.Restapi
  Wake *restapi.Restapi
  Setchargelimit *restapi.Restapi
  StopCharging *restapi.Restapi
  NearbyCharging *restapi.Restapi
  SSOrefreshtoken *restapi.Restapi

  DataRequestMap map[string]rest_cmds

  VehicleData TeslaVehicleData

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

func (et *MyTesla) SSOauthorizeURL() string{

  url := fmt.Sprintf("https://auth.tesla.com/oauth2/v3/authorized?client_id=%s&code_challenge=%s&code_challenge_method=%s&redirect_uri=%s&response_type=%s&scope=%s&state=%s",
                    "ownerapi", 
                    et.SSO_challengeSum,
                    "S256",
                    "https://auth.tesla.com/void/callback",
                    "code",
                    "openid email offline_access",
                    et.SSO_state)

  return url

}

func (et *MyTesla) setchargelimitURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s/command/set_charge_limit", TESLA_API_URL, id)

  return url

}

func (et *MyTesla) stopChargingURL(id string) string{

  url := fmt.Sprintf("%s/api/1/vehicles/%s/command/charge_stop", TESLA_API_URL, id)

  return url

}

func (et *MyTesla) refreshtokenURL(id string, sec string, refreshtoken string) string{

  url := fmt.Sprintf("%s/oauth/token?grant_type=refresh_token&client_id=%s&client_secret=%s&refresh_token=%s", TESLA_API_URL,
                    id, 
                    sec,
                    refreshtoken)

  return url

}

func (et *MyTesla) ssorefreshtokenURL() string{

  url := "https://auth.tesla.com/oauth2/v3/token"
  //url := "http://localhost:3000/refreshtoken"


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


func New(dbName string) *MyTesla{

  t := new(MyTesla) 


  logmsg.Print(logmsg.Info, "In MyTesla New");

  jsonFile, err := os.Open("modelx.json")

  if(err == nil){
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    json.Unmarshal([]byte(byteValue), &t.Modelxoptions)

  }else{
    logmsg.Print(logmsg.Error, "unable to open file")
  }

  t.myDB = new(MyDatabase)

  t.myDB.init(dbName)


  logmsg.Print(logmsg.Info,"=== DEBUG TURNED OFF, use debugOn to enable stdout of dumps")
  t.debug = false

  t.VehicleList = nil
  t.DataRequestMap = make(map[string]rest_cmds)
  t.dataRequestMapAdd("vehicle_data","vehicle_id", "Gets all the data")
  t.dataRequestMapAdd("charge_state","vehicle_id", "Gts charge state data")
  t.dataRequestMapAdd("climate_state","vehicle_id", "Gets climate state data")
  t.dataRequestMapAdd("drive_state","vehicle_id", "Gets drive state data")
  t.dataRequestMapAdd("gui_settings","vehicle_id", "Gets gui settings data")
  t.dataRequestMapAdd("vehicle_config","vehicle_id", "Gets vehicle config data")
  t.dataRequestMapAdd("vehicle_state","vehicle_id", "Gets vehicle state data")
  t.dataRequestMapAdd("nearbycharging","vehicle_id", "Gets vehicle state data")
  t.dataRequestMapAdd("service_data","vehicle_id", "Gets service data")

  // SSO setup stuff - begin
  // note - i copied this directly from :
  // github.com/uhthomas/tesla_exporter/cmd/login

     // this doesn't have to be 9 bytes, or base64. Just preference.
     var b [9]byte
     if _, err := io.ReadFull(rand.Reader, b[:]); err != nil {
             fmt.Errorf("rand state: %w", err)
             return t // egc - probably not the right thing :)
     }

     t.SSO_state = base64.RawURLEncoding.EncodeToString(b[:])

     var p [86]byte
     if _, err := io.ReadFull(rand.Reader, p[:]); err != nil {
             fmt.Errorf("rand challenge: %w", err)
             return t // egc - probably not the right thing :)
     }   

     t.SSO_challenge = base64.RawURLEncoding.EncodeToString(p[:])
     sum := sha256.Sum256([]byte(t.SSO_challenge))
     t.SSO_challengeSum = base64.RawURLEncoding.EncodeToString(sum[:])

/*
fmt.Printf("sso_state: [%s]\n", t.SSO_state);
fmt.Printf("sso_challenge: [%s]\n", t.SSO_challenge);
fmt.Printf("sso_challengeSum: [%s]\n", t.SSO_challengeSum);
*/

  // SSO setup stuff - end




  return t

}

func (et *MyTesla) debugOn(){
  et.debug = true
}

func (et *MyTesla) debugOff(){
  et.debug = false
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

func (et *MyTesla) IsHomeLink() bool{
  return et.myDB.IsHomeLink();
}

func (et *MyTesla) SetHomeLinkOn() bool{
  return et.myDB.SetHomeLinkOn();
}

func (et *MyTesla) SetHomeLinkOff() bool{
  return et.myDB.SetHomeLinkOff();
}

func (et *MyTesla) AddOwner(){

   et.myDB.AddOwner(et.email, et.accessToken, et.refreshToken, et.expiresTime);

}

func (et *MyTesla) AddVehicleId(id string){

fmt.Println("Get/Add VehicleId is being deprecated - use GetVehicleIdFromVinCmd")

   et.myDB.AddVehicleId(id);

}

func (et *MyTesla) AddVehicleVin(vin string){

   et.myDB.AddVehicleVin(vin);

}

func (et *MyTesla) GetOwner() (bool){

  return et.myDB.GetOwner(&et.email, &et.accessToken,
                                &et.refreshToken, &et.expiresTime)

}

func (et *MyTesla) GetVehicleId() (bool, string){

  var id string

fmt.Println("Get/Add VehicleId is being deprecated - use GetVehicleIdFromVinCmd")


  return et.myDB.GetVehicleId(&id), id
  

}

func (et *MyTesla) GetVehicleVin() (bool, string){

  var id string

  return et.myDB.GetVehicleVin(&id), id

}

func (et *MyTesla) DelOwner() (bool){

  return et.myDB.DelOwner()

}

func (et *MyTesla) DelVehicleId() (bool){

  return et.myDB.DelVehicleId()

}

func (et *MyTesla) RefreshToken(skipLogin bool) bool{

  if(!skipLogin){ // if skipping we already have the owner info
    et.Login() // the act of logging in will populate this info
  }

  logmsg.Print(logmsg.Info, "Starting RefreshToken")

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
      if(et.debug){ et.Wake.Dump() }
      return true
    }

    logmsg.Print(logmsg.Warning,"Wake still waiting. State:", stateStr)

  }

  logmsg.Print(logmsg.Error,"wake failed")
  return false

}

func (et *MyTesla) SSORefreshToken() bool{

  // fmt.Println("wow - called SSORefreshToken")

  dberr := et.GetSecrets()
  if(!dberr){
    logmsg.Print(logmsg.Error, "Yikes - DB error!.  Have you stored secrets?");
    os.Exit(4);
  }

  // see if we already have an access token

  status := et.GetOwner()

  if(!status){
    fmt.Println("Initial owner info include access tokens are missing - use importtoken command to fix")
    os.Exit(4)
  }

  et.SSOrefreshtoken = restapi.NewPost("SSOrefreshtoken", et.ssorefreshtokenURL())

  //et.SSOrefreshtoken.DebugOn()

  et.SSOrefreshtoken.SetBearerAccessToken(et.accessToken)
  //et.SSOrefreshtoken.HasInnerMap("response")

  jsonstr := fmt.Sprintf("{\"grant_type\": \"refresh_token\", \"client_id\": \"ownerapi\", \"refresh_token\": \"%s\", \"scope\": \"openid email offline_access\" }", et.refreshToken)

//fmt.Println(jsonstr)

  et.SSOrefreshtoken.SetPostJson(jsonstr)

  //et.SSOrefreshtoken.Dump()

  if(et.SSOrefreshtoken.Send()){
    if(et.debug){et.SSOrefreshtoken.Dump()}
  }else{
    logmsg.Print(logmsg.Error,"SSOrefreshtoken"," Failed ")
    return false
  }

  //et.SSOrefreshtoken.Dump()

  et.accessToken = et.SSOrefreshtoken.GetValueString("access_token")
  et.refreshToken = et.SSOrefreshtoken.GetValueString("refresh_token")
  expires := et.SSOrefreshtoken.GetValue("expires_in")
  //expires := 2592000; // 30 days - tesla is sending back 300

/*
  fmt.Println("returned expires time:", expires)
  fmt.Printf("expires is [%T]\n", expires)

  fmt.Println("AccessToken:", et.accessToken)
  fmt.Println("RefreshToken:", et.refreshToken)

  fmt.Println("before ExpiresTime: ", et.expiresTime)
  fmt.Println("before ExpiresTime: ", et.ExpiresTimeStr())
*/

  timeNow := time.Now()
  created := timeNow.Unix()

  et.expiresTime = int(created) + restapi.CastFloatToInt(expires)
  //et.expiresTime = int(created) + int(expires)
  

/*
  fmt.Println("after ExpiresTime: ", et.expiresTime)
  fmt.Println("after ExpiresTime: ", et.ExpiresTimeStr())
*/


  et.AddOwner()

  return true

}

func (et *MyTesla) StopChargingCmd(id string) bool{

  //var stateStr string

  et.Login() // the act of logging in will populate this info

  et.StopCharging = restapi.NewPost("StopCharging", et.stopChargingURL(id))

  et.StopCharging.SetBearerAccessToken(et.accessToken)
  et.StopCharging.HasInnerMap("response")

  if(et.StopCharging.Send()){
    if(et.debug){et.StopCharging.Dump()}
  }else{
    logmsg.Print(logmsg.Error,"StopCharging"," Failed ", id)
    return false
  }

  return true

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
    if(et.debug){et.Setchargelimit.Dump()}
  }else{
    logmsg.Print(logmsg.Error,"Setchargelimit"," Failed ", id)
    return false
  }

/*
  for et.wake.Send() {

    stateStr = et.wake.GetValueString("state")

    if(strings.Compare(stateStr,"online") == 0){
      if(et.debug){et.wake.Dump()}
      return true
    }

    logmsg.Print(logmsg.Warning,"Wake still waiting. State:", stateStr)

  }
*/

  return true

}

func (et *MyTesla) GetVehicleData(id string) bool{

  et.Login() // the act of logging in will populate this info

  cmd := "vehicle_data"

  r := et.DataRequestMap[cmd]

  if(r.Obj == nil){ // not setup before
    url := fmt.Sprintf("%s/api/1/vehicles/%s/%s", TESLA_API_URL, id, cmd)
    r.Obj = restapi.NewGet(r.cmd, url) 
    r.Obj.SetBearerAccessToken(et.accessToken)
    r.Obj.HasInnerMap("response")
  }

  r.Obj.JsonOnly()

  if(r.Obj.Send()){
    //et.tmpobj.Dump()
  }else{
    logmsg.Print(logmsg.Error,cmd," Failed ", id)
    return false
  }

  json.Unmarshal(r.Obj.BodyBytes, &et.VehicleData)
  /*
  fmt.Println("in tesla class - ID: ", et.VehicleData.Response.ID)
  // why is display name not in the response?
  fmt.Println("in tesla class - ID: ", et.VehicleData.Response.DisplayName)
  */

  return true

}

func (et *MyTesla) DataRequest(id string, cmd string) bool{

  et.Login() // the act of logging in will populate this info

  r := et.DataRequestMap[cmd]

  if(r.Obj == nil){ // not setup before
    //url := fmt.Sprintf("%s/api/1/vehicles/%s/data_request/%s", TESLA_API_URL, id, cmd)
    url := fmt.Sprintf("%s/api/1/vehicles/%s/%s", TESLA_API_URL, id, cmd)
    r.Obj = restapi.NewGet(r.cmd, url) 
    r.Obj.SetBearerAccessToken(et.accessToken)
    r.Obj.HasInnerMap("response")
  }

  if(r.Obj.Send()){
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
    if(et.debug){et.SingleVehicle.Dump()}
  }else{
    logmsg.Print(logmsg.Error,"get vehicle failed:", id)
    return false
  }

/*
  //vehicle := new(RestVehicles)

  if(!et.SingleVehicle.sendGetSingleVehicle(id, et.accessToken)){
    logmsg.Print(logmsg.Error,"GET Vehicle failed")
    return false
  }
*/


  return true

}

func (et *MyTesla) GetVehicleIdFromVinCmd(vin string) (bool, string) {

  if(!et.GetVehicleListCmd() ){
    logmsg.Print(logmsg.Error,"get vehicles list failed")
    return false, "actionfailed"
  }  

  // ok - we have the list, lets see if we can find the vin

fmt.Println("Looking for vin: ", vin)

  count := et.VehicleList.GetValueInt("count")


  for j:= 0; j < count; j++ {

    if( et.VehicleList.GetArrayValueString(j,"vin") == vin){
      return true, et.VehicleList.GetArrayValueString(j,"id_s")
    }
  
  }

  return false, "not found in list"
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
    logmsg.Print(logmsg.Error,"GET Vehicles failed")
    return false
  }
*/


  return true

}

func (et *MyTesla) ExpiresTimeStr() string {

  var expireStr string = "nerd"
  
  expiresTime := time.Unix(int64(et.expiresTime), 0)

  expireStr = fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               expiresTime.Month(),expiresTime.Day(), expiresTime.Year(),
             expiresTime.Hour(), expiresTime.Minute(), expiresTime.Second() )

  return expireStr
}

func (et *MyTesla) DumpOwnerInfo(){

  et.Login() // the act of logging in will populate this info

  //expiresTime := time.Unix(int64(et.expiresTime), 0)

  fmt.Println("Owner table dump")
  fmt.Println("Email: ", et.email)
  fmt.Println("AccessToken: ", et.accessToken)
  fmt.Println("RefreshToken: ", et.refreshToken)
  fmt.Println("ExpiresTime: ", strconv.Itoa(et.expiresTime))
  fmt.Println("ExpiresTime: ", et.ExpiresTimeStr())
/*
  fmt.Println("ExpiresTime: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               expiresTime.Month(),expiresTime.Day(), expiresTime.Year(),
             expiresTime.Hour(), expiresTime.Minute(), expiresTime.Second() ) )
*/

  

}

func (et *MyTesla) ImportTokens() bool{

  var filename string
  var email string
  var ok bool

  fmt.Println("Import tokens from another authentication program/script")

  timeNow := time.Now()

  unixTime := timeNow.Unix()
  // cheesy way of having non magic numbers - but are magic :)
  one_hour := int(60 * 60)
  one_day  := int(24 * one_hour)
  num_days := int(40)

  et.expiresTime = int(unixTime) + (one_day * num_days)


  expiresTime := time.Unix(int64(et.expiresTime), 0)
  fmt.Printf("UnixTime[%d] vs expiresTime[%d]\n", unixTime, et.expiresTime)

  fmt.Println("TimeNow: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
               timeNow.Month(),timeNow.Day(), timeNow.Year(),
             timeNow.Hour(), timeNow.Minute(), timeNow.Second() ) )

  fmt.Println("ExpiresTime: ", fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d", 
             expiresTime.Month(),expiresTime.Day(), expiresTime.Year(),
             expiresTime.Hour(), expiresTime.Minute(), expiresTime.Second() ) )


  term := new(MyLogin) 

  status := et.GetOwner()

  logmsg.Print(logmsg.Info, "owner status: ", status)

  if(!status){ // need email
    email, ok = term.Prompt("Account Email", true);

    if(!ok){
      fmt.Println("Error getting email address");
      return false
    }

    et.email = email

  }


  filename, ok = term.Prompt("Token response json filename:", false);

  if(!ok){
    fmt.Println("Error inputing token");
    return false
  }

  fmt.Printf("File[%s]\n", filename)

  // Open our jsonFile
  jsonFile, err := os.Open(filename)

  // if we os.Open returns an error then handle it
  if err != nil {
      fmt.Println(err)
      return false
  }

  fmt.Println("Successfully Opened: ", filename)

  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  var result map[string]interface{}
  json.Unmarshal([]byte(byteValue), &result)


  myMap := restapi.CastMap(result)

  //fmt.Println(myMap)

  tokenMap := restapi.CastMap(myMap["tokens"])

  //fmt.Println(tokenMap)

 // accessToken := restapi.CastString(tokenMap["owner_access_token"])
  //ssoRefreshToken := restapi.CastString(tokenMap["sso_refresh_token"])

  et.accessToken = restapi.CastString(tokenMap["owner_access_token"])
  et.refreshToken = restapi.CastString(tokenMap["sso_refresh_token"])

  
  fmt.Println("accessToken: ", et.accessToken)
  fmt.Println()
  fmt.Println("refreshToken: ", et.refreshToken)

/*
  created := r.GetValue("created_at")
  expires := r.GetValue("expires_in")
  et.expiresTime = restapi.CastFloatToInt(created) + restapi.CastFloatToInt(expires)
*/


   et.myDB.AddOwner(et.email, et.accessToken, et.refreshToken, et.expiresTime);

  return true
}

func (et *MyTesla) SSOLogin() bool{

/*
  var email1 string
  var passwd1 string

  fmt.Println("In new SSO Login logic")

  email1 = "nerd";
  passwd1 = "geek";
*/

  dberr := et.GetSecrets()
  if(!dberr){
    logmsg.Print(logmsg.Error, "Yikes - DB error!.  Have you stored secrets?");
    os.Exit(4);
  }

  //fmt.Println(email1);
  //fmt.Println(passwd1);

  r := restapi.NewGet("authentication", et.SSOauthorizeURL())
  r.DebugOn()

  //r.Header.Add("User-Agent", "tesla_admin_cntl")

  r.Dump()

  if(r.Send()){
  //  r.Dump()
  }else{
    logmsg.Print(logmsg.Error,"authentication failed")
    return false
  }



  return true

}

func (et *MyTesla) Login() bool{

  //var email1 string
  //var email2 string
  //var passwd1 string
  //var passwd2 string
  //var reader *bufio.Reader
  //var ok bool

  dberr := et.GetSecrets()
  if(!dberr){
    logmsg.Print(logmsg.Error, "Yikes - DB error!.  Have you stored secrets?");
    os.Exit(4);
  }


  // see if we already have an access token

  status := et.GetOwner()

  logmsg.Print(logmsg.Info, "owner status: ", status)

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

// egc 3/6/2021 - removed 2 week pad, tesla updated the exipre time to
// 300 seconds (5minutes) so this forced a new token every since time
// if(int(unixTime) >= (et.expiresTime - const_sec_twoweek)){
 if(int(unixTime) >= et.expiresTime){
      logmsg.Print(logmsg.Info, "Time to refresh token")
      //status = et.RefreshToken(true)  
      status = et.SSORefreshToken()  
    }

    if(status){
      //logmsg.Print(logmsg.Info, fmt.Sprintf("Logined to [", et.email, "]"))
      logmsg.Print(logmsg.Info, "Logined to [", et.email, "]")
      return(true)
    }

  } // if status == true  Meaning we have data


  logmsg.Print(logmsg.Warning,"AccessToken not found")
  fmt.Println("AccessToken not found - run cheesy python script and import new token")

  return false

/*
  logmsg.Print(logmsg.Warning,"AccessToken not found")

  term := new(MyLogin)
  email1, passwd1, ok = term.TerminalLogin("Email Address", "Password")

  if(!ok){
    logmsg.Print(logmsg.Error, "Failure to get credentials")
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
*/
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

    logmsg.Print(logmsg.Error, "Error getting client ID and or secret.  URL returned:")
    logmsg.Print(logmsg.Error, fmt.Sprintf(tmpstr))
    return false
  }

  et.clientID = client_id
  et.clientSecret = client_secret

  et.AddSecrets();

//  fmt.Println("client id found: [", client_id, "]")
//  fmt.Println("client secret found: [", client_secret, "]")

  return true

}

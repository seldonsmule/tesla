package tesla

import "fmt"
import "time"
import "github.com/seldonsmule/logmsg"
import "github.com/denisbrodbeck/machineid"
import "database/sql"
//import "database/sql/driver"
import  _ "github.com/mattn/go-sqlite3"

//const constDbName = "./tesla.db"

type MyDatabase struct {

  nerd int
  handle *sql.DB  // i figured out the type by using the %T to print out the variable

  dbName string

}

func (edb *MyDatabase) checkErr(err error) {
  if(err != nil){
    panic(err.Error())
  }
  
}

func (edb *MyDatabase) createTable(tableName string, tableSql string) bool{

  if(edb.handle == nil){
    logmsg.Print(logmsg.Error, "database handle not initiated");
    return false
  }

  row, err := edb.handle.Query("select 1 from $1 limit 1", tableName);

  _ = row // found a way to undo the variable since we don't really need it
          // and we get an error otherwise.  I am sure some geek will comment
          // on how bad this code is :)

  if(err == nil){
    logmsg.Print(logmsg.Info, 
                 fmt.Sprintf("table already exist [%s]", tableName));
    return true
  }

  logmsg.Print(logmsg.Info, fmt.Sprintf("Create db table [%s]", tableName));

  stmt, err := edb.handle.Prepare(tableSql);
  logmsg.Print(logmsg.Debug03, stmt)

  if(err != nil){
    logmsg.Print(logmsg.Error, fmt.Sprintf("Invalid SQL: %s", err.Error()))
    return false
  }
  
  stmt.Exec()


  return true

}

func (edb *MyDatabase) GetOwner(pEmail *string, 
                                      pAccessToken *string, 
                                      pRefreshToken *string,
                                      pExpiresTime *int ) (bool) {

  var email string
  var accessToken string
  var refreshToken string
  var expiresTime int
  var gotRow bool
  var expires time.Time


  //sqlmsg := fmt.Sprintf("SELECT * FROM owner;");

  rows, err := edb.handle.Query("SELECT * FROM owner;");

  if(err != nil){
    logmsg.Print(logmsg.Error, "Db error: ", err)
    return false
  }

  defer rows.Close() // close the resource later 

  gotRow = false;

  for rows.Next(){
//fmt.Println("in rows.Next loop")

    gotRow = true;

    //err = rows.Scan(&email, &accessToken, &refreshToken, &expiresTime)
    err = rows.Scan(&email, &accessToken, &refreshToken, &expires)

//fmt.Println("Email[",email,"]")

    expiresTime = int(expires.Unix())

    if(err != nil){
      logmsg.Print(logmsg.Error, "Db error: ", err)
      return false
    }
  }

  if(!gotRow){
    logmsg.Print(logmsg.Error, "AccessToken not set");
    return false
    
  }

  *pEmail = email
  *pAccessToken = accessToken
  *pRefreshToken = refreshToken
  *pExpiresTime = expiresTime

  return true
}

func (edb *MyDatabase) GetVehicleId(pId *string) (bool) {

  var id string
  var gotRow bool

  //sqlmsg := fmt.Sprintf("SELECT * FROM vehicle_id;");

  rows, err := edb.handle.Query("SELECT * FROM vehicle_id;");

  if(err != nil){
    logmsg.Print(logmsg.Error, "Error getting VehicleID Db error: ", err)
    return false
  }

  defer rows.Close() // close the resource later 

  gotRow = false;

  for rows.Next(){
//fmt.Println("in rows.Next loop")

    gotRow = true;

    err = rows.Scan(&id)

    logmsg.Print(logmsg.Info, "VehicleId[",id,"]")

    if(err != nil){
      logmsg.Print(logmsg.Error, "Db error: ", err)
      logmsg.Print(logmsg.Error, "Error getting VehicleID (not stored) Db error: ", err)
      return false
    }
  }

  if(!gotRow){
    logmsg.Print(logmsg.Error, "VehicleId not set");
    return false
    
  }

  *pId = id

  return true
}

func (edb *MyDatabase) GetVehicleVin(pVin *string) (bool) {

  var id string
  var gotRow bool

  rows, err := edb.handle.Query("SELECT * FROM vin;");

  if(err != nil){
    logmsg.Print(logmsg.Error, "Error getting Vin Db error: ", err)
    return false
  }

  defer rows.Close() // close the resource later 

  gotRow = false;

  for rows.Next(){
//fmt.Println("in rows.Next loop")

    gotRow = true;

    err = rows.Scan(&id)

    logmsg.Print(logmsg.Info, "Vin[",id,"]")

    if(err != nil){
      logmsg.Print(logmsg.Error, "Db error: ", err)
      logmsg.Print(logmsg.Error, "Error getting Vin (not stored) Db error: ", err)
      return false
    }
  }

  if(!gotRow){
    logmsg.Print(logmsg.Error, "Vin not set");
    return false
    
  }

  *pVin = id

  return true
}

func (edb *MyDatabase) DelVehicleId() bool {

  //sql := fmt.Sprintf("DELETE FROM vehicle_id;")

  edb.handle.Exec("DELETE FROM vehicle_id;")

  return true
}

func (edb *MyDatabase) DelVehicleVin() bool {

  edb.handle.Exec("DELETE FROM vin;")

  return true
}

func (edb *MyDatabase) DelOwner() bool {

  //sql := fmt.Sprintf("DELETE FROM owner;")

  edb.handle.Exec("DELETE FROM owner;")

  return true
}

func (edb *MyDatabase) AddOwner(email string, 
                                      accessToken string, 
                                      refreshToken string,
                                      expiresTime int ) (bool) {

  logmsg.Print(logmsg.Debug03, 
               fmt.Sprintf("Email[%s] AccessToken[%s] RefreshToken[%s] Expires[%d]",
                           email, accessToken, refreshToken, expiresTime));

  // first delete any existing entry

  //sql := fmt.Sprintf("DELETE FROM owner;")

  edb.handle.Exec("DELETE FROM owner;")

//  sql = fmt.Sprintf("INSERT INTO owner (email, access_token, refresh_token, expires_in) VALUES ('%s', '%s', '%s', '%d');",
//                     email, accessToken, refreshToken, expiresTime)

  edb.handle.Exec("INSERT INTO owner (email, access_token, refresh_token, expires_in) VALUES ($1, $2, $3, $4);",
                     email, accessToken, refreshToken, expiresTime)

  return true
}

func (edb *MyDatabase) AddVehicleId(id string) (bool) {

  // first delete any existing entry

//  sql := fmt.Sprintf("DELETE FROM vehicle_id;")

  edb.handle.Exec("DELETE FROM vehicle_id;")

//  sql = fmt.Sprintf("INSERT INTO vehicle_id (id) VALUES ('%s');",
//                     id)

  edb.handle.Exec("INSERT INTO vehicle_id (id) VALUES ($1);", id)

  return true
}

func (edb *MyDatabase) AddVehicleVin(id string) (bool) {

  // first delete any existing entry

  edb.handle.Exec("DELETE FROM vin;")

  edb.handle.Exec("INSERT INTO vin (id) VALUES ($1);", id)

  return true
}


func (edb *MyDatabase) GetApiDetails(pId *string, pSecret *string) (bool){

  var clientID string
  var clientSecret string
//  var id int
  var gotRow bool

  //sqlmsg := fmt.Sprintf("SELECT * FROM api_details LIMIT 1;");

  rows, err := edb.handle.Query("SELECT * FROM api_details LIMIT 1;");

  if(err != nil){
    logmsg.Print(logmsg.Error, "Db error: ", err)
    return false
  }

  defer rows.Close() // close the resource later 

  gotRow = false;

  for rows.Next(){

    gotRow = true;

    err = rows.Scan(&clientID, &clientSecret)

    if(err != nil){
      logmsg.Print(logmsg.Error, "Db error: ", err)
      return false
    }
  }

  if(!gotRow){
    logmsg.Print(logmsg.Error, "ClientID and ClientSecret not set");
    return false
    
  }

  *pSecret = clientSecret
  *pId = clientID

  return true
}

func (edb *MyDatabase) AddApiDetails(clientId string, clientSecret string) bool{

  //var sql2 string

  c := clientId
  s := clientSecret

  logmsg.Print(logmsg.Info, "AddApiDetails clientId: ", c); 
  logmsg.Print(logmsg.Info, "AddApiDetails clientSecret: ", s); 


  // first delete any existing entry

//  sql := fmt.Sprintf("DELETE FROM api_details;")

  edb.handle.Exec("DELETE FROM api_details;")

  //sql2 = fmt.Sprintf("INSERT INTO api_details (client_id, client_secret) VALUES ('%s', '%s');", c, s)

  edb.handle.Exec("INSERT INTO api_details (client_id, client_secret) VALUES ($1, $2);", c, s)

  return true
}

  

func (edb *MyDatabase) init(dbName string){

  logmsg.Print(logmsg.Debug03, "in init");

  edb.dbName = dbName
  logmsg.Print(logmsg.Debug03, edb.dbName)

  db, err := sql.Open("sqlite3", edb.dbName)

  edb.checkErr(err)

  edb.nerd = 2
  edb.handle = db


  if(!edb.createTable("api_details", "CREATE TABLE `api_details` (`client_id` VARCHAR(256) NULL, `client_secret` VARCHAR(256) NULL)") ){
    logmsg.Print(logmsg.Warning,"Unable to create table api_details")
  }

  if(!edb.createTable("owner", "CREATE TABLE `owner` (`email` VARCHAR(64) NULL, `access_token` VARCHAR(256) NULL, `refresh_token` VARCHAR(256) NULL, `expires_in` DATE NULL) ") ){
    logmsg.Print(logmsg.Warning,"Unable to create table owner")
  }

  if(!edb.createTable("tamper", "CREATE TABLE `tamper` (`machineid` VARCHAR(256) NULL) ") ){
    logmsg.Print(logmsg.Warning,"Unable to create table tamper")
  }

  if(!edb.createTable("vehicle_id", "CREATE TABLE `vehicle_id` (`id` VARCHAR(256) NULL) ") ){
    logmsg.Print(logmsg.Warning,"Unable to create table vehicle_id")
  }

  if(!edb.createTable("vin", "CREATE TABLE `vin` (`id` VARCHAR(256) NULL) ") ){
    logmsg.Print(logmsg.Warning,"Unable to create table vin")
  }

  if(!edb.createTable("homelink", "CREATE TABLE `homelink` (`setup` VARCHAR(256) NULL, `homelogic` INTEGER) ") ){
    logmsg.Print(logmsg.Warning,"Unable to create table homelink")
  }


  edb.initHomeLink();

  id , _ := machineid.ProtectedID(edb.dbName);

  logmsg.Print(logmsg.Info, "machineid: ", id);


  //sqlmsg := "select * from tamper limit 1";

  row, err := edb.handle.Query("select * from tamper limit 1;")

  //_ = row // found a way to undo the variable since we don't really need it
          // and we get an error otherwise.  I am sure some geek will comment
          // on how bad this code is :)

  if(err != nil){
    logmsg.Print(logmsg.Error, "Error reading table tamper");
    return 
  }

  defer row.Close()

  var mid string

  row.Next()

  err = row.Scan(&mid)
  
  logmsg.Print(logmsg.Info, "mid = ", mid)

  if(mid == ""){
    logmsg.Print(logmsg.Error, "machine id not set")

//    sqlinsert := fmt.Sprintf("INSERT INTO tamper (machineid) VALUES('%s');",
//                             id)
//    logmsg.Print(logmsg.Debug03, sqlinsert)
    edb.handle.Exec("INSERT INTO tamper (machineid) VALUES($1);",
                             id)
  }else{
    logmsg.Print(logmsg.Info, "we have the mid")
  }


}

func (edb *MyDatabase) SetHomeLink(state bool) bool{

  logmsg.Print(logmsg.Info, "Setting Homelink state to:", state);

  edb.handle.Exec("DELETE FROM homelink;")


  edb.handle.Exec("INSERT INTO homelink (setup, homelogic) VALUES ($1, $2);",
                     "yes",
                     state)

  return true
}

func (edb *MyDatabase) SetHomeLinkOn() bool{

  return (edb.SetHomeLink(true))
}

func (edb *MyDatabase) SetHomeLinkOff() bool{

  return (edb.SetHomeLink(false))
}

func (edb *MyDatabase) intToBool(value int) bool{

  if(value == 1){
    return true
  }

  return false
}

func (edb *MyDatabase) IsHomeLink() bool{

  row, err := edb.handle.Query("select * from homelink limit 1;")

  //_ = row // found a way to undo the variable since we don't really need it
          // and we get an error otherwise.  I am sure some geek will comment
          // on how bad this code is :)

  if(err != nil){
    logmsg.Print(logmsg.Error, "Error reading table homelink");
    return false 
  }

  defer row.Close()

  var setup string
  var homelogic int

  row.Next()

  err = row.Scan(&setup, &homelogic)


  logmsg.Print(logmsg.Info, "setup = ", setup)

  if(setup == ""){
    logmsg.Print(logmsg.Error, "setup not set")
    edb.SetHomeLinkOn()
    homelogic = 1
  }

  logmsg.Print(logmsg.Info, "Homelinklogic is:", edb.intToBool(homelogic));

  return edb.intToBool(homelogic)
}

func (edb *MyDatabase) initHomeLink() bool{

  row, err := edb.handle.Query("select * from homelink limit 1;")

  //_ = row // found a way to undo the variable since we don't really need it
          // and we get an error otherwise.  I am sure some geek will comment
          // on how bad this code is :)

  if(err != nil){
    logmsg.Print(logmsg.Error, "Error reading table homelink");
    return false 
  }

  defer row.Close()

  var setup string
  var homelogic int

  row.Next()

  err = row.Scan(&setup, &homelogic)


  logmsg.Print(logmsg.Info, "setup = ", setup)

  if(setup == ""){
    logmsg.Print(logmsg.Error, "setup not set")

    edb.SetHomeLinkOn()
  }else{
    logmsg.Print(logmsg.Info, "homelink logic already set")
    edb.IsHomeLink() // forces a print to the log of current state
  }


  return true
}


func (db *MyDatabase) hello(){

  logmsg.Print(logmsg.Debug01, "hi there")
}

func (db *MyDatabase) Hello(){
  db.hello()
}

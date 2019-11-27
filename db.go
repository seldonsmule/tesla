package main

import "fmt"
import "time"
import "github.com/seldonsmule/logmsg"
import "github.com/denisbrodbeck/machineid"
import "database/sql"
//import "database/sql/driver"
import  _ "github.com/mattn/go-sqlite3"

const constDbName = "./tesla.db"

type MyDatabase struct {

  nerd int
  handle *sql.DB  // i figured out the type by using the %T to print out the variable

}

func (edb *MyDatabase) checkErr(err error) {
  if(err != nil){
    panic(err.Error())
  }
  
}

func (edb *MyDatabase) createTable(tableName string, tableSql string) bool{

  var sqlmsg string

  if(edb.handle == nil){
    logmsg.Print(logmsg.Error, "database handle not initiated");
    return false
  }

  sqlmsg = fmt.Sprintf("select 1 from %s limit 1", tableName);

  row, err := edb.handle.Query(sqlmsg);
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


  sqlmsg := fmt.Sprintf("SELECT * FROM owner;");

  rows, err := edb.handle.Query(sqlmsg);

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

func (edb *MyDatabase) DelOwner() bool {

  sql := fmt.Sprintf("DELETE FROM owner;")

  edb.handle.Exec(sql)

  return true
}

func (edb *MyDatabase) AddOwner(email string, 
                                      accessToken string, 
                                      refreshToken string,
                                      expiresTime int ) (bool) {

  logmsg.Print(logmsg.Info, 
               fmt.Sprintf("Email[%s] AccessToken[%s] RefreshToken[%s] Expires[%d]",
                           email, accessToken, refreshToken, expiresTime));

  // first delete any existing entry

  sql := fmt.Sprintf("DELETE FROM owner;")

  edb.handle.Exec(sql)

  sql = fmt.Sprintf("INSERT INTO owner (email, access_token, refresh_token, expires_in) VALUES ('%s', '%s', '%s', '%d');",
                     email, accessToken, refreshToken, expiresTime)

  edb.handle.Exec(sql)

  return true
}


func (edb *MyDatabase) GetApiDetails(pId *string, pSecret *string) (bool){

  var clientID string
  var clientSecret string
//  var id int
  var gotRow bool

  sqlmsg := fmt.Sprintf("SELECT * FROM api_details LIMIT 1;");

  rows, err := edb.handle.Query(sqlmsg);

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

  var sql2 string

  c := clientId
  s := clientSecret

  logmsg.Print(logmsg.Info, "AddApiDetails clientId: ", c); 
  logmsg.Print(logmsg.Info, "AddApiDetails clientSecret: ", s); 


  // first delete any existing entry

  sql := fmt.Sprintf("DELETE FROM api_details;")

  edb.handle.Exec(sql)

  sql2 = fmt.Sprintf("INSERT INTO api_details (client_id, client_secret) VALUES ('%s', '%s');", c, s)

  edb.handle.Exec(sql2)

  return true
}

  

func (edb *MyDatabase) init(){

  logmsg.Print(logmsg.Debug03, "in init");
  logmsg.Print(logmsg.Debug03, constDbName)

  db, err := sql.Open("sqlite3", constDbName)

  edb.checkErr(err)

  edb.nerd = 2
  edb.handle = db


  if(!edb.createTable("api_details", "CREATE TABLE `api_details` (`client_id` VARCHAR(256) NULL, `client_secret` VARCHAR(256) NULL)") ){
    logmsg.Print(logmsg.Error,"Unable to create table api_details")
  }

  if(!edb.createTable("owner", "CREATE TABLE `owner` (`email` VARCHAR(64) NULL, `access_token` VARCHAR(256) NULL, `refresh_token` VARCHAR(256) NULL, `expires_in` DATE NULL) ") ){
    logmsg.Print(logmsg.Error,"Unable to create table owner")
  }

  if(!edb.createTable("tamper", "CREATE TABLE `tamper` (`machineid` VARCHAR(256) NULL) ") ){
    logmsg.Print(logmsg.Error,"Unable to create table tamper")
  }

  id , _ := machineid.ProtectedID(constDbName);
fmt.Println(id)


  sqlmsg := "select * from tamper limit 1";

  row, err := edb.handle.Query(sqlmsg);
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
  
  fmt.Printf("mid = %s\n", mid)

  if(mid == ""){
    fmt.Println("machine id not set")

    sqlinsert := fmt.Sprintf("INSERT INTO tamper (machineid) VALUES('%s');",
                             id)
    fmt.Println(sqlinsert)
    edb.handle.Exec(sqlinsert)
  }else{
    fmt.Println("we have the id")
  }


}

func (db *MyDatabase) hello(){

  fmt.Println("hi there")
}

func (db *MyDatabase) Hello(){
  db.hello()
}

3/6/2021 - updated for new tesla authentication method

will require a python app for the initial token

# Telsa class
This class is an example of how to use the unpublished tesla APIs to get information about a tesla and to change settings.
# Needed documentation
The class uses the unpublished Tesla API to access a Telsa.  Below are two different references.  Both write-ups were needed to figure out how to communicate.
### https://tesla-api.timdorr.com
### https://www.teslaapi.io
### https://pastebin.com/pS7Z6yyP
Secret values for the Authentication

# Capabilities

## Setup (house keeping)
These methods are used setup the authentication information into a database for future use with the Tesla API. The class has an sqlite3 database to store this information. **It is recommend that you look at how to do a better job of securing this information.  This is just an example.**

### Secrets
These are referenced in the authentication piece of the API.  They are setup to be modified in the database in case they change. They never have changed since this project was started.  The latest versions are here: https://pastebin.com/pS7Z6yyP
#### add secrets
Adds client id and client secret to database
#### get secrets
Gets the current values
#### updates secrets
Changes the current values

## Usage recommendation
Use tesla_admin (see below) for all management of your database of credentials and vehicles.  Then have your application using the tesla library do the rest of the work.  You can do it differenntly - but the initial login to get the accesstoken has been changed by tesla.  Instead, use the python script (in the python director) to get a json version of your initial accesstoken and ssorefreshtoken.  Store that in a file and then use the "tesla_admin -cmd importtoken" command to import into your database.


## Authentication
Provides authentication.  Uses the user id and password of the tesla mobile app.  Will receive a token that is used for auth going forward. The class has an sqlite3 database to store this information. **It is recommend that you look at how to do a better job of securing this information.  This is just an example.**

#### login
Does not work anymore - Tesla changed their api.  See: https://tesla-api.timdorr.com/api-basics/authentication#post-https-auth-tesla-com-oauth-2-v-3-token-1

for details.  You will have to use the script in the python directory to create your initial accesstoken.

#### refreshtoken
Refreshes the token so user id/password prompting is not required.  **login** above will auto call **refreshtoken** based on the time to live value in the token
#### wake
Wakes up the Telsa's computer for conversations with the rest of the APIs.  From lots of testing, the **wake** command is important to have responded as successful.  You can complete a **login** succesfully, but the car will not respond.  ***Even if you have just completed talking to the car, it goes to sleep quickly :)***
### get owner
Provides details about the owner of the Tesla as stored in the database.  ***Does not show password, but does show email and token information***
### delete owner
Removes all information from the database

# Vehicle
All of the get/set commands require the unique vehicle ID.  These methods provide two ways to get that information.  See https://tesla-api.timdorr.com/api-basics/vehicles for details
### Get Vehicle list
List of all vehicles registered to a single Telsa account - **email address**
### Get Vehicle
Provides information about a specific vehicle based on the vehicle ID

## Get Information
Provides a json of different types of information about the vehicle.  See the API documentation links provided above for the details
### Service data, Charge state, Drive state, GUI settings, Vehicle state, Vehicle config, Near By Charging

## Set Information
Allows for the setting various items on the vehicle.  At this point, the class only implemeted a call to **setchargelimit**.  Others can be added as needed
### Set Charge Limit
Sets the charging limit from 50-100% for the next charging cycle



# Example code
## tesla_cmd - Demo of how to use the tesla.go class
```
tesla login | getvehiclelist | addsecrets | getsecrets | getowner | delowner | refreshtoken | getvehiclelist | getchargestate | getclimatestate | updatesecrets 


Get State Cmds
nearbycharging vehicle_id : Gets vehicle state data
service_data vehicle_id : Gets service data
charge_state vehicle_id : Gts charge state data
climate_state vehicle_id : Gets climate state data
drive_state vehicle_id : Gets drive state data
gui_settings vehicle_id : Gets gui settings data
vehicle_config vehicle_id : Gets vehicle config data
vehicle_state vehicle_id : Gets vehicle state data

Set Cmds
wake vehicle_id - Wake up Vehicle
setchargelimit vehicle_id percent - Sets the charge limit
```
# Utilities
## tesla_admin - Admin command to setup the tesla.db for use with other utilities
```
  -clientid string
    	API Client ID (default "notset")
  -clientsec string
    	API Client Secret (default "notset")
  -cmd string
    	Command to run (default "help")
  -dbname string
    	Name of database (default "tesla.db")
  -rundir string
    	Directory to exec from (default "./")
  -vid string
    	VehicleId (default "notset")

cmds:
     Addsecrets - Add Tesla API secrets - needed for all APIs
     getsecrets - Displays Tesla API secrets - needed for all APIs
     updatesecrets - Updates Tesla API secrets - needed for all APIs
     login - Sets up login info
     refreshtoke - Refreshes token created via login
     getowner - Shows owner details
     delowner - Deletes owner details
     setid - Stores vehicle ID for the cmds.  Requires -vid
     getid - Display vehicle ID for the cmds
     getvehiclelist - Displays a list of vehicles and their IDs owned by the login
     help - Display this help
```
## tesla_chargelevel - Simple command to run from cron to reest the charging levels
```
  -cmd string
    	Command to run (default "help")
  -dbname string
    	Name of database (default "tesla.db")
  -highlimit string
    	charge high limit - used with homecharge (default "notset")
  -limit string
    	charge limit - used with setchargelimit (default "notset")
  -lowlimit string
    	charge low limit - used with homecharge (default "notset")
  -rundir string
    	Directory to exec from (default "./")

cmds:
     wake - Wake up vehicle
     setchargelimit - Set charge limit for next charge. Use -limit
     getchargelimit - Get charge limit for next charge
     batterylevel - Get current battery level
     homesetto50 - If vehicle is home, set to 50% charge limit
     homesetto80 - If vehicle is home, set to 80% charge limit
     homecharge - If vehicle is home, use -lowlimit and -highlimit level to adjust next charge limit
             -lowlimit - Level to set to so car does not charge every night
             -highlimit - Level to set to so car for next charge IF car is at low (or below) -lowlimit
             NOTE: If set to 100 already and car is below - will skip adjustments.  If charged, will set to low limit
```

# Dependancies

## Add these packages
go get github.com/mattn/go-sqlite3

go get github.com/seldonsmule/logmsg

go get github.com/seldonsmule/restapi

go get golang.org/x/crypto/ssh/terminal

go get github.com/denisbrodbeck/machineid


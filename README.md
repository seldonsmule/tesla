# Telsa class
This class is an example of how to use the unpublished tesla APIs to get information about a tesla and to change settings.

# Capabilities
The class uses the unpublished Tesla API (https://tesla-api.timdorr.com) to access a Tesla.

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


## Authentication
Provides authentication.  Uses the user id and password of the tesla mobile app.  Will receive a token that is used for auth going forward. The class has an sqlite3 database to store this information. **It is recommend that you look at how to do a better job of securing this information.  This is just an example.**

#### login
If a valid OAuth token is not available, will prompt for user id and password
#### refreshtoken
Refreshes the token so user id/password prompting is not required.  **login** above will auto call **refreshtoken** based on the time to live value in the token
#### wake
Wakes up the Telsa's computer for conversations with the rest of the APIs.  From lots of testing, the **wake** command is important to have responded as successful.  You can complete a **login** succesfully, but the car will not respond.  ***Even if you have just completed talking to the car, it goes to sleep quickly :)***
### get owner
Provides details about the owner of the Tesla as stored in the database.  ***Does not show password, but does show email and token information***
### delete owner
Removes all information from the database

## updatesecrets
## get secrets
## add secrets
## refreshtoken
## getowner
## delowner
## wake
## setchargelimit
## service_data
## charge_state
## drive_state
## gui_settings
## vehicle_state
## vehicle_config
## nearbycharging
## getvehicle
## getvehicle list



# Example/Utilties

# Dependancies

Inital code to access a tesla

More doc coming later


Add these packages
go get github.com/mattn/go-sqlite3

go get github.com/seldonsmule/logmsg

go get github.com/seldonsmule/restapi

go get golang.org/x/crypto/ssh/terminal

go get github.com/denisbrodbeck/machineid

You need to have access to a tesla and a login for this to work


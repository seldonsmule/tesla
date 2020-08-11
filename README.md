# Telsa class
This class is an example of how to use the unpublished tesla APIs to get information about a tesla and to change settings.

# Capabilities
The class uses the unpublished Tesla API (https://tesla-api.timdorr.com) to access a Tesla.

## Setup (house keeping)
These methods are used setup the authentication information into a database for future use with the Tesla API. The class has an sqlite3 database to store this information. **It is recommend that you look at how to do a better job of securing this information.  This is just an example.**

### Secrets
These are referenced in the authentication piece of the API.  They are setup to be modified in the database in case they change. They never have changed since this project was started.  The latest versions are here: https://pastebin.com/pS7Z6yyP
#### add secrets - adds client id and client secret to database
#### get secrets - gets the current values
#### updates secrets changes the current values


## Authentication
Provides authentication.  Uses the user id and password of the tesla mobile app.  Will receive a token that is used for auth going forward.  
#### login
#### refreshtoken
#### wake

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


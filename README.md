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



# Example/Utilties

# Dependancies

## Add these packages
go get github.com/mattn/go-sqlite3

go get github.com/seldonsmule/logmsg

go get github.com/seldonsmule/restapi

go get golang.org/x/crypto/ssh/terminal

go get github.com/denisbrodbeck/machineid


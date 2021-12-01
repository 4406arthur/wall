# wall

This service used to be a developer manager Auth token service

## Architecture

![alt text](CleanArchitecture.jpg "Clean Arch")

* Independent of Frameworks: The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
* Testable: The business rules can be tested without the UI, Database, Web Server, or any other external element.
* Independent of UI: The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
* Independent of Database: You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
* Independent of any external agency: In fact your business rules simply don’t know anything at all about the outside world


## API SPEC
getToken

### Description

this api could get a jwt token

### Endpoint

     HTTP POST /developer/getToken

### Request Header filed description

* **"Authorization: Basic {BasicToekn}"**: only accept Basic auth token issusd by wall service 
* **"Content-Type: application/json"**: request payload only accpet json format

### Request Body filed description

* **"userID"** *(O,string)*: jwt token consumer

### Response field description

* **"token"** *(M,string)*

### JWT payload description

To make collaboration 3rd party service do more authorization policy this JWT token will provide some extra info


* **"aud"** *(M,string)*: means developerID which developer generate this token
* **"sub"** *(M,string)*: means userID, the token conusmer
* **"scopes"** *(M,listOfString)*: extra info for 3rd party service do something like RBAC

## Dev

please ensure you have install golang version same as go.mod determin

```sh
go mod vendor
docker-compose up --build
```

## Test

go test


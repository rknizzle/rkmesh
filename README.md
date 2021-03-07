# RKMESH
A platform for managing 3D models for manufacturing and 3D-printing

# Starting Development in this Repository
## Dependencies:  
TODO: Instructions to install dependencies  
- docker & docker-compose
- Make

## Run locally
`make start-local`  
The above command will use docker-compose to spin up a postgres database and run the app locally.
The app server will reload on any file changes.  

# Database Migrations
go-migrate [docs](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md)

### Create a new migration
`migrate create -ext sql -dir migrations -seq create_examples_table`

TODO: migrate CLI should be a dependency so that the up AND down migrations can be tested.


# Backend Architecture
Based on [go-clean-arch](https://github.com/bxcodec/go-clean-arch)
## 4 layers
#### domain  
Contains any objects structs and methods. Used in all layers.  
#### controller  
Accepts inputs and responds. This layer will mainly be a REST API for this app   .  
#### service  
Contains business logic.  
#### repository
Handles interactions with any external data sources. (databases, microsercives, file storage)  

# Testing
### Controller
Starts an actual HTTP server and sends an HTTP request to it. Mocks the service layer to return mock
data

### Service
Mocks the repository layer to return mock data
  
### Repository
TBD: 
... Maybe by simulating a database connection when running queries but it might make more sense to
just have integration tests that touch a database to cover the repository layer


# Mocking
Generate mocks of all the interfaces in domain
```
cd domain
mockery --all
```

TODO: generate mocks of a new package

# Connect to local database with psql
`psql postgresql://postgres:postgres@localhost:5432/postgres`

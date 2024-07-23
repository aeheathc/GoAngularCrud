# GoAngularCrud

A simple demo app showing CRUD functionality using:
- Go backend
  - "Squirrel" SQL builder
  - "Echo" HTTP framework
- Angular frontend


## Run

Use either of these methods to run, and once the builds are done visit http://localhost:4200

### with Docker

With the Docker daemon running, run `docker-compose up -d` in the project root

### separately

- Have Go and npm installed
- Have a postgres server running.
  - Set its connection string in the env variable DATABASE_URL if it needs to be different than the app's default of postgresql://admin:password@localhost/GoAngularCrud
  - Run the create table statement found in sql.sh
- In the backend folder run `go run .`
- In the frontend folder run `npm start`


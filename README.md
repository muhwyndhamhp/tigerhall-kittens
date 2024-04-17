# Tigerhall Kittens

## Introduction
This project is to fulfill my application as a Senior Backend Engineer at Tigerhall. 

## Demo
You can access the demo via [https://tigerhall-kittens.fly.dev](https://tigerhall-kittens.fly.dev)
It has 2 built-in options for GraphQL Editor / Playground:
- [Altair (Recommended)](https://tigerhall-kittens.fly.dev/altair)
- [GraphiQL](https://tigerhall-kittens.fly.dev/graphiql)

## Infrastructure Overview
![Infra Overview](schema.png)
The infrastructure follows the "Progressive Enhancement" concept, which means the system will be able to self-sustain, but will become increasingly more capable as more feature available to itself. 

One of such example is the core of the system is the GraphQL Server, which is responsible for handling the business logic and data access. The GraphQL Server is connected to Turso Database, wich uses LibSQL (SQLite fork with enhancement in distributed replica) as the core driver, but to run the app it does not requires Turso DB connection as it able to run locally, whilst the Turso DB is used for the production environment.

The GraphQL Server also has a Message Queue, which is used to send email notifications. The Message Queue is connected to the Email Service via SendGrid. The Email Service is responsible for sending the email notifications whenever new sightings are created. 


## Local Setup
***(If You're from Tigerhall, Please Use `.env` file provided in the email I sent!!)***

### Prerequisites
- Go 1.16
- Make
- [Air](https://github.com/cosmtrek/air) (Optional, for hot reload)
- [Mockery](https://github.com/vektra/mockery) (Optional, for mocking interfaces)
- [GqlGen](https://gqlgen.com/) (Optional, for generating GraphQL Schema)

### Steps
To run the project locally, you can follow these steps:
- Clone the repository
- Create the `.env` file based on `.env.example`. Graph below will explain the environment variables below.
- Run the following command:
```bash
make auto-migrate
```
- Run the following command to start the server:
```bash
make run
```
### Environment Variables
| Variable Name | Description | Default Value | Required |
| ------------- | ----------- | ------------- | -------- |
| `PORT` | Port for the server to run | `8080` | Yes |
| `BASE_URL` | Base URL for the server | `http://localhost:8080` | Yes |
| `LIBSQL_URL` | URL for LibSQL | - | No (Will Default to Local DB File) |
| `LIBSQL_TOKEN` | Token for LibSQL | - | No |
| `CF_ACCOUNT_ID` | Cloudflare Account ID | - | Yes |
| `CF_R2_ACCESS_KEY_ID` | Cloudflare R2 Access Key | - | Yes |
| `CF_R2_SECRET_ACCESS_KEY` | Cloudflare R2 Secret Access Key | - | Yes |
| `JWT_SECRET ` | Secret for JWT | `MuhWyndham-TigerHall-Kittens-Test` | Yes |
| `SENDGRID_API_KEY` | SendGrid API Key | - | Yes |
| `SENDGRID_SENDER_EMAIL` | SendGrid Email Origin | - | Yes |

## All Available Commands
- `make auto-migrate` : Run the automigrate to create the tables
- `make dry-run-migrate` : Run the automigrate without executing the migration. Use this to see the SQL that will be executed when running the migration.
- `make run` : Run the server
- `make test` : Run the unit tests with coverage report (will automatically open browser window)
- `make gen` : Generate the GraphQL Schema
- `make gen-mocks` : Generate the mocks for the interfaces
- `make print-dsn` : Print the DSN for the current database connection

## Todo List
- [x] Graph the Infrastructure Design
- [x] Create new GraphQL Server with GqlGen
- [x] Create Schema for `User`, `Tiger`, and `Sighting`
- [x] Connect Turso DB
- ~~[ ] Create Migration mechanism (Probably using Automigrate Dry Run & gomigrate)~~
- [x] Create Automigrate
- [x] Create Auth Implementation using JWT and Middleware
- [x] Create CRUD for `User`, `Tiger`, and `Sighting` w/ Pagination
- [x] Implement Sighting Rules (Only Beyond 5 km from prev. Sightings)
- [x] Implement Image Upload for Sightings
- [x] Create Message Queue using Go Channel and Send Email Notification on Consumer Side
- [ ] Add transaction for Create operations
- [x] Create Unit Test for Each Function
  - [x] Create Unit Test for Sighting
  - [x] Create Unit Test for User
  - [x] Create Unit Test for Tiger
- [ ] Create Integration Test for Each Endpoint
- [x] Create Documentation
- [x] Create Dockerfile
- [x] Create Fly io Deployment
- [ ] Fix Migration to use proper Migration Mechanism
- [ ] Better Error Presentation (Maybe add extensions by appending error codes?)

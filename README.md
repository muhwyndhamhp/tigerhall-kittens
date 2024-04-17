# Tigerhall Kittens

## Introduction
This project is to fulfill my application as a Senior Backend Engineer at Tigerhall. 

## Demo
You can access the demo at [https://tigerhall-kittens.fly.dev/altair.html](https://tigerhall-kittens.fly.dev/altair.html)

## Infrastructure Overview
![Infra Overview](schema.png)

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

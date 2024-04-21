# Database
This section will contains the documentation of the database implementation.
## What is Turso?
Turso was built on top of LibSQL. LibSQL itself is a fork of SQLite that was modified to be used in edge-computing environments, and has core strenght in it's ability to maintain consistency between multiple replicas (including on-device replicas!). 

It also maintains 100% compatibility with SQLite Query, Syntax, and Planner so you can use it as a drop-in replacement for SQLite. In fact, I use SQLite as the Database for Testing in this very project but use Turso in the running environment.

## Why Turso?
I use Turso here for several strenghts it has over other databases:
- It's very easy to spin new database.
- It has reasonable performance, in the tens of thousands of queries per second.
- It is fully managed and we can manage the databases both using the CLI and the Web App (https://turso.tech).
- It has incredibly generous free tier.
- Because it is a fork of SQLite, it is very easy to use and has a lot of documentation available.
- Very lightweight and support in-memory database. Very useful for testing purposes and quick prototyping.

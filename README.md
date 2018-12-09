# user-profile-api

## Overview ##
This repository contains the necessary files to host restful API's using Protocol Buffers (a.k.a protobuf) under golang to run a database. Information on protocol buffers
can be found on [protobufs Google Developers site](https://developers.google.com/protocol-buffers/docs/proto3).
All of the endpoints are hosted using [Heroku](https://www.heroku.com). The database was implemented using [MongoDB](https://mongodb.com)
with the help of the public MongoDB driver [mgo](https://github.com/globalsign/mgo) and is being hosted using [mLab](https://mlab.com).
This repository handles requests from multiple pages of our website dealing with user profile information.

## Program Execution ##
Make sure [mgo](https://github.com/globalsign/mgo), [glog](https://github.com/golang/glog), [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), 
[cors](https://github.com/rs/cors), and [grpc](https://godoc.org/google.golang.org/grpc) are installed in your golang environment. To execute the program 
run the server.go file as follows,

	go run server.go

This will execute the server file.

## Endpoints ##
Each endpoint expects to receive specific fields to process a request. The following are the expectations for each endpoint and the response. The API link is https://tea-user-profile-api.herokuapp.com/. Note that each endpoint only accepts POST requests, so the link will only display a "Not Found" message since GET requests aren't supported.

| Endpoint | Request | Response |
|:--------:|---------|----------|
| getuserprofile | string email = 1; | string profileimage = 1;<br>string profiledescription = 2;<br>repeated string endorsements = 3;<br>repeated string currentprojects = 4;<br>repeated string previousprojects = 5;|
| updateuserprofile | string email = 1;<br>string profileimage = 2;<br>string profiledescription = 3;<br>repeated string currentprojects = 4;<br>repeated string previousprojects = 5;| bool success = 1; |

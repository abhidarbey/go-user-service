# GO User Service

## This projects is based upon Golang v1.14

## Description
This project has  4 Domain layer :
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Controller/Delivery Layer (uses gin-gonic framework)

#### The diagram:

<img src="/docs/img/clean-arch.png">

This project uses Go Module

#### Run the Applications

#move to directory
$ cd workspace

# Clone this repo
$ git clone https://github.com/abhidarbey/user-service.git

#move to project
$ cd user-service

# Build the docker image first

# Run the application

# check if the containers are running

# Execute the call
$ curl localhost:8080/ping

# Stop
```


### Tools Used:
In this project, I use some tools listed below. But you can use any simmilar library that have the same purposes. But, well, different library will have different implementation type. Just be creative and use anything that you really need. 

- All libraries listed in [`go.mod`](https://github.com/abhidarbey/user-service/go.mod)
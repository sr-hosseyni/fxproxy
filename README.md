### About the task ###

fxporxy code challenge is designed to demonstrate [sidecar pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar) implementation. fxproxy is sidecar service which is responsible to govern incoming traffic to the downstream application services based on allowed path list `allowedList` in `main.go`.


### Task objective ###

* implement `ValidatePath` and pass failing tests in `main_test.go`, however we don't want to limit you so feel free to change the function name or even project structure/packages in way you believe can suit a production grade Go project. Code and tests are added to give you an idea about the business requirement of the task.

* It won't be a sidecar service if it doesn't handle incoming Http request and send it to downstream application service, so please extend project functionality to handle http proxy responsibilities well, you can use any 3rd party packages or stick with standard libraries.

* Provide docker-compose to run the sidecar service along side with a dummy application service 

#### Good to have ####

* Implementing the task in TDD way
* Effective use of source control
* Follow Idiomatic Go
* Design a loosely coupled and highly maintainable code
* Since this service is expected to handle all traffic to pod/application service, code is expected to be fast and resilient

## RUN
```shell script
git clone git@github.com:sr-hosseyni/fxproxy.git
cd fxproxy
docker-compose up -d
```

http://localhost:8888/company  
http://127.0.0.1:8000/health-check

## TODO
* Request handling duration (time takes to get response from downstream) and some other metrics.
* Current fxproxy container has prepared for development env by hot-reload, For deployment we need to compile app and copy
 native binary inside a production docker image.
* fxproxy should be able to be more configurable!
* It must be able to handle http 2 in communication with other sidecar proxies, and even RPC
* Needs more test!
* Functionality for capability to connect to service discovery.
* Handling outgoing traffic from the backend micro-service to other services
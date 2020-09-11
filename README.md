### About the task ###

fxporxy code challenge is designed to demonstrate [sidecar pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar) implementation. fxproxy is sidecar service which is responsible to govern incoming traffic to the downstream application services based on allowed path list `allowedList` in `main.go`.


### Task objective ###

* implement `ValidatePath` and pass failing tests in `main_test.go`, however we don't want to limit you so feel free to change the name or project structure/packages in way you believe can suit a production grade Go project. Code and tests are added to give you an idea about the business requirement of the project. 

* It won't be a sidecar service if it doesn't handle incoming Http request and send it to downstream application service, so please extend project functionality to handle http proxy responsibilities well, you can use any 3rd party packages or stick with standard libraries.

* Provide docker-compose to run the sidecar service along side with a dummy application service 

#### Good to have ####

* Implementing the task in TDD way
* Effective use of source control
* Follow Idiomatic Go
* Design a loosely coupled and highly maintainable code
* Since this service is expected to handle all traffic to pod/application service, code is expected to be fast and resilient
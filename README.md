### About the task ###

fxporxy is designed to demonstrate [sidecar pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar) implementation. fxproxy is sidecar service which is responsible to govern incoming traffic to the downstream application services (REVERSE PROXY) based on allowed path list in configuration.


## RUN
```shell script
git clone git@github.com:sr-hosseyni/fxproxy.git
cd fxproxy
docker-compose up -d
```

http://localhost:8888/company  
http://127.0.0.1:8000/health-check

## TODO
* General structure of project needs improvments!
* Request handling duration (time takes to get response from downstream) and some other metrics.
* Current fxproxy container has prepared for development env by hot-reload, For deployment we need to compile app and copy native binary inside a production docker image.
* fxproxy should be able to be more configurable!
* It must be able to handle http 2 in communication with other sidecar proxies, and even RPC
* Needs more test!
* Functionality for capability to connect to service discovery.
* Handling outgoing traffic from the backend micro-service to other services (DIRECT PROXY)
* Using logrus for structured logging
* Using viper for configuration management
* Using cobra for extending CLI functions
* A sidecar can be under very heavy traffic and easily can turn into single point of failure. to avoid it and have better resilience we can use in memory caching mechanism, utilize goroutine to process path in concurrent way rather than loop. Also having benchmark test which is part of Go testing
* Makefile to handle test, build, deploy commands could help streamlining the development/deployment process
* Docker file could be more optimized by using multi-stage building approach
# go-service-template

Draft template for Go microservices. See [main.go](main.go), [app.go](app/app.go) and [service.go](service/service.go) for the meaty bits.

This template covers:
- Go microservice
- [ConnectRPC](https://connectrpc.com/) ("gRPC but better", and backwards compatible)
- Simple but clean Dependency Injection (see [app.go](app/app.go))
- Graceful termination handling (see [app.go](app/app.go))
- Easy E2E testing, reusing the live app config (see [service_test.go](service/service_test.go))
- Example middleware ([auth](service/interceptors/auth.go))

Intended to be converted into a [cookiecutter](https://www.cookiecutter.io/) / [cruft](https://cruft.github.io/cruft/) template for easier deployment and maintenance.

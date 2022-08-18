# IBAN Demo Project

## Introduction
This project implements a simple IBAN format validator.

## Development Progres Documentation
Link: [Here](./docs/Development.md)

## Development Environment
For details, please refer to the [Makefile](./Makefile)

### Important Commands:

```bash
make image # Builds the runtime image

make run # Builds and runs the runtime image (webserver on port 18888)
HTTP_PORT=22334 make run # Same as above, but with a different port

make test # Runs all tests
make coverage # Creates a test coverage report and automatically opens it in the browser

make vendor # Updates dependencies

make serve_docs # Serves the Swagger documentation
```

## Project Requirements

* Ref: [IBAN](https://en.wikipedia.org/wiki/International_Bank_Account_Number)
* You are not allowed to to use someone else's IBAN library. Please write the solution
  yourself and cover it with tests. Note: Libraries are allowed for web-server itself.
* Your solution should be documented well enough so your future colleagues will be
  able to easily run it and use it in a micro-service environment.
* Your solution should be cross-platform and work on both MacOS and Ubuntu.
    * Please include a Dockerfile and verify that you can run your app in Docker
* Do not over-engineer it, please write just enough code and tests.
* We expect you to use Git, and we are interested in looking at your commits.
* This isn’t meant to take too much of your time, so if you find yourself stuck on
  anything, please reach out.

## Implementation Steps
(in no particular order)

- [ ] HTTP API with one endpoint: `/validate/{iban}`
    - [ ] JSON response with:
        * `error:string` is the error message if the IBAN is invalid and null otherwise
        * `isValid:bool` is true if the IBAN is valid and false otherwise
        * `parsedComponents:struct` contains componentized IBAN if isValid is true and null otherwise
            * `country:string` is the country code
            * `checkDigits:uint` is the check digits
            * `bankCode:string` is the bank code (looks like this can be a string?)
            * `accountNumber: uint` is the account number
    - [ ] [Swagger](https://swagger.io/) specifications
        - [ ] Include Swagger-Docs in `make generate` command
    - [ ] Unit Tests
    - [ ] Functional Tests
- [ ] Service Layer
    - [ ] Investigate IBAN structure
    - [ ] Regex-Hashmap vs. abstraction layer
    - [ ] Exported generic validation method
        - e.g. `func (svc *Service) Parse(iban string) (IBAN, error)`
    - [ ] Unit tests
- [ ] Dev-Environment
    - [x] Makefile
    - [x] Compile in Build Container
    - [x] Create Runtime Image
    - [x] `make build` to compile binary
    - [x] `make image` to create runtime image
    - [x] `make run` to run/stop compiled binary on dev machine
    - [ ] `make test` to run unit & functional tests
    - [ ] `make mod` to update and vendor dependencies
- [ ] CI
    - [ ] `make lint` to run linters
    - [ ] Include Mock-Generation in `make generate` command
    - [ ] Verify Generated Files are up-to-date
    - [ ] GitHub Actions for all of the above
- [ ] CD
    - [ ] Push Swagger Docs to Swagger Hub (and tag new version + latest)
    - [ ] K8S Manifests
        - [ ] Namespace
        - [ ] Deployment
        - [ ] Service
        - [ ] Ingress
    - [ ] GitHub Action for tagging new version by label
    - [ ] GitHub Action for applying k8s manifests to EKS cluster

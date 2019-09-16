# tradeshift
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

## How to run
1. Build with `$ CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-s' -o appbinary`
2. `docker build -t mvp .`
3. `docker-compose up` (Known issue: May need to ru multiple times if `api` container starts before the `db-msq` container)


## Sample use-case:
- Create node :
`curl --request POST http://localhost:8080/node/create -d '{"id":"root", "pid":""}' --header "Content-Type: application/json" --ipv4`
- Get childrens:
`curl http://localhost:8080/children/root`
- Update parent:
`curl --request PUT http://localhost:8080/node/b/make_parent/d  --header "Content-Type: application/json" --ipv4`


### TODO:
- Input validation
- Test coverage (Unit+Functional)
- Rate limit
- Code smell (Some rough edges)
- Make sure `api` starts after `db-msq`. Now this behaviour is randome.

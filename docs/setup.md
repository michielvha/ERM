# Setup commands

## init

````shell
go clean -modcache
go mod tidy
go build -o main .
````

## refresh

````shell
rm go.mod
rm  go.sum
go mod init github.com/MKTHEPLUGG/ERM
go get github.com/golang-migrate/migrate/v4@v4.18.1
go get github.com/golang-migrate/migrate/v4/source/iofs
go get github.com/lib/pq
go get github.com/rs/zerolog/log
go get github.com/golang-jwt/jwt/v5
go get github.com/gin-gonic/gin
go get github.com/rs/zerolog
````



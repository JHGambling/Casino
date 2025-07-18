module jhgambling/backend

go 1.23.2

require jhgambling/protocol v0.0.0

replace jhgambling/protocol => ../protocol

require (
	github.com/fatih/color v1.18.0
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/gorilla/websocket v1.5.3
	github.com/matoous/go-nanoid/v2 v2.1.0
	golang.org/x/crypto v0.39.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.28 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)

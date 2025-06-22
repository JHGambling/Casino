# Example Game Plugin

## Build

To build a plugin, run this command:
```sh
go build -buildmode=plugin -o example.so main.go
```

Replace `example` with a better suiting name for you plugin and then move this file into /casino-backend/games/ so it can be loaded

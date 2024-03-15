@echo off
go generate
go install github.com/tc-hib/go-winres@latest
go-winres simply --icon web/static/favico.ico
go build -ldflags "-H windowsgui" -tags=release -o bin/v2ray-panel-windows.exe
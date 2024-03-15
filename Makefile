.PHONY: server client install_service

client:
	@go build -o bin/v2ray-panel-cilent cmd/client/main.go

server:
	@go build -o bin/v2ray-panel-server cmd/server/main.go


install_service:
	@./scripts/install_service.sh


service_start:
	@systemctl start v2ray-panel.service

service_stop:
	@systemctl stop v2ray-panel.service

service_restart:
	@systemctl restart v2ray-panel.service

log:
	@grep "v2ray-panel" /var/log/syslog
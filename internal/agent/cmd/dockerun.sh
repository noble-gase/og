#!/bin/bash

docker rm -f {{.AppName}}
docker rmi -f {{.AppName}}
{{ if eq .DockerF "Dockerfile"}}
docker build -t {{.AppName}} .
{{- else}}
docker build -f {{.DockerF}} -t {{.AppName}} .
{{- end}}
docker image prune -f

docker run -d --name={{.AppName}} --restart=always --privileged -p 10086:8000 -v /data/{{.AppName}}:/data {{.AppName}}

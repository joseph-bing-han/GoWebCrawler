#!/usr/bin/env bash
echo "Build start."
cd `dirname $0`
go build -o bin/spider src/spider/main/spider.go
go build -o bin/crawler src/crawler/main/crawler.go
go build -o bin/updater src/updater/main/updater.go
go build -o bin/CacheClear src/tools/CacheClear.go
go build -o bin/MQClear src/tools/MQClear.go
go build -o bin/ProxyFresh src/tools/ProxyFresh.go
echo "Build complete"



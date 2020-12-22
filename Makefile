cli:
	go build -mod vendor -o bin/server cmd/server/main.go

debug:
	go run -mod vendor cmd/server/main.go -enable-www -enable-properties -spatial-database-uri 'sqlite:///?dsn=$(DSN)' -properties-reader-uri 'sqlite:///?dsn=$(DSN)' -geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=$(DSN)' -nextzen-apikey $(APIKEY) -mode directory:// 

debug-fs:
	go run -mod vendor cmd/server/main.go -enable-www -enable-properties -spatial-database-uri 'sqlite:///?dsn=$(DSN)' -properties-reader-uri 'sqlite:///?dsn=$(DSN)' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-woeplanet:
	go run -mod vendor cmd/server/main.go -enable-www -geojson-path-resolver-uri wofid:// -enable-properties -spatial-database-uri 'sqlite:///?dsn=$(DSN)' -properties-reader-uri 'sqlite:///?dsn=$(DSN)' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-geojson:
	go run -mod vendor cmd/server/main.go -enable-www -geojson-path-resolver-uri wofid:// -enable-properties -spatial-database-uri 'sqlite:///?dsn=$(DSN)&index-geojson=true' -properties-reader-uri 'sqlite:///?dsn=$(DSN)' -geojson-reader-uri 'sql://sqlite3/geojson/id/body?dsn=$(DSN)' -nextzen-apikey $(APIKEY) -mode featurecollection:// $(GEOJSON)


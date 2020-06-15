#!/bin/bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get github.com/twinj/uuid
go get go.mongodb.org/mongo-driver/bson
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options
go get github.com/lithammer/shortuuid
go get github.com/twinj/uuid
(cd app && rm -f app && go build && echo "Build success" && cd ..)


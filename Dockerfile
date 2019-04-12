FROM golang:1.11.2-alpine3.8

WORKDIR /go/src/github.com/kine-dmd/api/

COPY . .
RUN rm -rf vendor/ && rm -f **/mock*.go

EXPOSE 80

RUN apk add --no-cache git
RUN apk add --no-cache gcc
RUN apk add --no-cache libc-dev

RUN wget -O dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64
RUN echo '287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep' | sha256sum -c -
RUN chmod +x dep
RUN mv dep /usr/bin/


RUN dep ensure
RUN GOBIN=$PWD/vendor/bin/ go install ./vendor/github.com/golang/mock/mockgen/

RUN vendor/bin/mockgen -destination=kinesisqueue/mock_kinesis_queue.go -package=kinesisqueue github.com/kine-dmd/api/kinesisqueue KinesisQueueInterface
RUN vendor/bin/mockgen -destination=dynamoDB/mock_dynamo_db.go -package=dynamoDB github.com/kine-dmd/api/dynamoDB DynamoDBInterface
RUN vendor/bin/mockgen -destination=api_time/mock_time.go -package=api_time github.com/kine-dmd/api/api_time ApiTime
RUN vendor/bin/mockgen -destination=watch_position_db/mock_watch_pos_db.go  -package=watch_position_db github.com/kine-dmd/api/watch_position_db WatchPositionDatabase

RUN go test -v ./...

RUN rm -f **/*_test.go && rm -rf vendor/bin && rm -f **/mock*.go

RUN go build -o ~/go/bin/main .

ENTRYPOINT ~/go/bin/main

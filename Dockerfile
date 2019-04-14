FROM golang:1.11.2-alpine3.8

WORKDIR /go/src/github.com/kine-dmd/api/

COPY . .

# Remove old dependencies and mocks (should also be in dockerignore)
RUN rm -rf vendor/ && rm -f **/mock*.go

# Expose the web port
EXPOSE 80

# Add C dependensies required by Go runtime
RUN apk add --no-cache git
RUN apk add --no-cache gcc
RUN apk add --no-cache libc-dev

# Intstall Go Dep package manager
RUN wget -O dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64
RUN echo '287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep' | sha256sum -c -
RUN chmod +x dep
RUN mv dep /usr/bin/

# Install dependency libraries and mockgen for testing
RUN dep ensure
RUN GOBIN=$PWD/vendor/bin/ go install ./vendor/github.com/golang/mock/mockgen/

# Generate mocks
RUN vendor/bin/mockgen -destination=kinesisqueue/mock_kinesis_queue.go -package=kinesisqueue github.com/kine-dmd/api/kinesisqueue KinesisQueueInterface
RUN vendor/bin/mockgen -destination=dynamoDB/mock_dynamo_db.go -package=dynamoDB github.com/kine-dmd/api/dynamoDB DynamoDBInterface
RUN vendor/bin/mockgen -destination=api_time/mock_time.go -package=api_time github.com/kine-dmd/api/api_time ApiTime
RUN vendor/bin/mockgen -destination=watch_position_db/mock_watch_pos_db.go  -package=watch_position_db github.com/kine-dmd/api/watch_position_db WatchPositionDatabase
RUN vendor/bin/mockgen -destination=apple_watch_3/mock_data_writer.go  -package=apple_watch_3 github.com/kine-dmd/api/apple_watch_3 Aw3DataWriter
RUN vendor/bin/mockgen -destination=binary_file_appender/mock_binary_file_appender.go  -package=binary_file_appender github.com/kine-dmd/api/binary_file_appender BinaryFileAppender

# Run all tests
RUN go test -v ./...

# Remove mocks and test files
RUN rm -f **/*_test.go && rm -rf vendor/bin && rm -f **/mock*.go

# Build the executable and remove all source files to save space
RUN go build -o ~/go/bin/main .
RUN rm -rf *

ENTRYPOINT ~/go/bin/main

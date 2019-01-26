FROM golang:1.11.2-alpine3.8

WORKDIR /go/src/github.com/kine-dmd/api/

COPY . .

EXPOSE 80

RUN apk add --no-cache git
RUN apk add --no-cache gcc
RUN apk add --no-cache libc-dev

RUN wget -O dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64
RUN echo '287b08291e14f1fae8ba44374b26a2b12eb941af3497ed0ca649253e21ba2f83  dep' | sha256sum -c -
RUN chmod +x dep
RUN mv dep /usr/bin/


RUN dep ensure

RUN go test -v ./...

RUN rm -f **/*_test.go

RUN go build -o ~/go/bin/main .

ENTRYPOINT ~/go/bin/main

FROM golang:1.21

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# build binary and place into /usr/local/bin/
COPY . .
RUN go build -v -o /usr/local/bin/q3rcon-proxy ./cmd/q3rcon-proxy/

# Command to run when starting the container
ENTRYPOINT [ "q3rcon-proxy" ]
FROM golang:alpine

WORKDIR /dist

COPY . .

# build binary and place into /usr/local/bin
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/q3rcon-proxy ./cmd/q3rcon-proxy

# Command to run when starting the container
ENTRYPOINT [ "q3rcon-proxy" ]
FROM golang:1.21 AS build_image

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# build binary, place into ./bin/
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/q3rcon-proxy ./cmd/q3rcon-proxy/

FROM scratch AS final_image

WORKDIR /bin/

COPY --from=build_image /usr/src/app/bin/q3rcon-proxy .

# Command to run when starting the container
ENTRYPOINT [ "./q3rcon-proxy" ]
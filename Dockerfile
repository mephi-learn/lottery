FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app ./app/server

# Final image
FROM gcr.io/distroless/static-debian11:nonroot

WORKDIR /app

COPY --from=build /go/bin/app ./app

USER nonroot:nonroot

ENTRYPOINT ["/app/app", "-help"]

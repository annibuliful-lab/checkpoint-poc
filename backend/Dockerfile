FROM golang:1.21.6-alpine3.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
COPY ../generated.graphql ./
COPY ../.env ./.env
RUN go mod download

# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend

EXPOSE 8080

CMD ["/backend"]
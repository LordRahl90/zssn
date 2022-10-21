FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o zssn ./cmd/


FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/zssn zssn

EXPOSE 8080

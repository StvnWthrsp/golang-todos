FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /todos-server

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /todos-server /todos-server

EXPOSE 8080

USER nonroot:nonroot

CMD [ "/todos-server" ]
FROM golang:1.19.3-buster AS build

WORKDIR /src/github.com/gerardrbentley/places/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . /src/github.com/gerardrbentley/places/
RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux go build -o /bin/app && \
    chmod 111 /bin/app

FROM gcr.io/distroless/base as final
COPY --from=build /bin/app /bin/app
CMD ["/bin/app"]

FROM golang:1.13 AS builder

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# copy files
WORKDIR /go/src/github.com/DostonAkhmedov/lucy
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

# run test
# RUN go test -v ./...
# RUN go test -bench=. -v ./...
# RUN rm -rf *_test.go

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lucy .

FROM alpine:latest

# copy app
WORKDIR /app
COPY --from=builder /go/src/github.com/tusupov/exmoarbitrage/lucy .

# create appuser
RUN adduser -S -D -H -h /app appuser
USER appuser

CMD ["./lucy"]
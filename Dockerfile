FROM golang:1.19 AS build
COPY . /work
WORKDIR /work
RUN go build -trimpath -ldflags '-s -w' -o main
RUN ldd main 2>&1 | grep -q "not a dynamic executable"

FROM scratch
COPY --from=build /work/main /
ENTRYPOINT ["/main"]

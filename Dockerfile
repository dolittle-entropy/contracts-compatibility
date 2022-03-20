FROM --platform=$BUILDPLATFORM golang:1.18.0 as build
ARG TARGETARCH

WORKDIR /go/src/contracts-compatibility
COPY . /go/src/contracts-compatibility
ENV CGO_ENABLED=0
ENV GOARCH $TARGETARCH
RUN go build -o /go/bin/contracts-compatibility .

FROM gcr.io/distroless/static-debian11
COPY --from=build /go/bin/contracts-compatibility /bin/contracts-compatibility
ENTRYPOINT ["/bin/contracts-compatibility"]
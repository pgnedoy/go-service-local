FROM golang:1.14.2 as build

WORKDIR /src/app

COPY . .

ARG GITHUB_DEPLOY_KEY

RUN bin/build

FROM gcr.io/distroless/base
COPY --from=build /bin/app /

CMD ["/app", "server"]

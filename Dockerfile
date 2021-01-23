FROM golang:1.15.6-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN go build -o /server .

FROM alpine:3.12.3
COPY --from=build /src/swagger.yaml /app/
COPY --from=build /src/html/verifyTokenEmail.html /app/html/
COPY --from=build /src/html/itemOpenGraph.html /app/html/
COPY --from=build /server /app/server

CMD ["/app/server"]

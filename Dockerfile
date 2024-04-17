FROM golang:1.22-alpine3.19 as builder

RUN apk update && apk add git make bash

WORKDIR /app

COPY . ./

# Do dep installs outside, due to private git modules
# RUN make dep

RUN go build -v -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/
COPY --from=builder /app/public /app/public
COPY --from=builder /app/utils/email /app/utils/email

EXPOSE 4001

CMD [ "/app/main" ]


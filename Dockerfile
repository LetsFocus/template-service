FROM alpine:edge
COPY ./configs /configs
ADD main /main
RUN apk add --no-cache postgresql-client
EXPOSE 9001
CMD ["/main"]

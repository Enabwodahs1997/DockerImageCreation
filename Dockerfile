FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o genx .

FROM scratch
COPY --from=build /src/genx /genx
LABEL maintainer="Enabwodahs (Summer)" \
    description="Go data generator with linear, cosine, log, exp, and web preview modes" \
    cohort="22" \
    animal="Wolf"
EXPOSE 8080/tcp
ENTRYPOINT ["/genx"]
CMD ["--help"]
# I wasn't in class when the animal was chosen. So I picked one at random.
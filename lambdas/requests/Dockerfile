FROM golang:1.18.1-buster
# install compiler
RUN apt update
RUN apt install -y golang tar
RUN curl -fsSL https://get.pulumi.com | sh
ENV PATH "${HOME}/.pulumi/bin:${PATH}"
RUN go env -w GOPROXY=https://proxy.golang.org,direct

WORKDIR /app
COPY . .

EXPOSE 8080

RUN go mod init requests-lambda
RUN go mod tidy
RUN go build requests.go

ENTRYPOINT [ "/app/requests" ]

# copy artifacts to a clean image
# FROM public.ecr.aws/lambda/provided:al2
# RUN yum install tar gzip -y
# RUN curl -fsSL https://get.pulumi.com | sh
# ENV PATH "/root/.pulumi/bin:${PATH}"
# WORKDIR /app
# COPY --from=build /app/requests .
# ENTRYPOINT [ "/app/requests" ]
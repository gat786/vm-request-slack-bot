FROM public.ecr.aws/lambda/provided:al2 as build
# install compiler
RUN yum install -y golang tar
RUN curl -fsSL https://get.pulumi.com | sh
RUN go env -w GOPROXY=https://proxy.golang.org,direct

WORKDIR /app
COPY . .

RUN go mod init requests-lambda
RUN go mod tidy
RUN go build requests.go

# copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
RUN yum install tar gzip -y
RUN curl -fsSL https://get.pulumi.com | sh
ENV PATH="/root/.pulumi/bin:${PATH}"
WORKDIR /app
COPY --from=build /app/requests .
ENTRYPOINT [ "/app/requests" ]
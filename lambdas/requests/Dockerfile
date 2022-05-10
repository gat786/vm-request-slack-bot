FROM pulumi/pulumi-go:latest as build

RUN go env -w GOPROXY=https://proxy.golang.org,direct

WORKDIR /app
COPY . .

RUN go mod init requests-lambda
RUN go mod tidy
RUN go build requests.go

# copy artifacts to a clean image
FROM pulumi/pulumi-go:latest
WORKDIR /
COPY --from=build /app/requests .

ENV APP_ENV ""
ENV PULUMI_ACCESS_TOKEN ""
ENV LINODE_TOKEN ""

# (Optional) Add Lambda Runtime Interface Emulator and use a script in the ENTRYPOINT for simpler local runs
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod 755 /usr/bin/aws-lambda-rie
COPY entry.sh /
RUN chmod 755 /entry.sh
ENTRYPOINT [ "/entry.sh" ]
CMD [ "/requests" ]
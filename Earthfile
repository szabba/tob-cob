VERSION 0.6

linux-x64:
    FROM ubuntu:bionic-20230308

    WORKDIR /work

    ENV BUILD_FLAGS "-trimpath"

    RUN apt-get update -qq
    RUN apt-get install -qy gcc pkg-config xorg-dev libglfw3-dev curl
    RUN mkdir -p ./out/dl/go
    RUN curl -L -o ./out/dl/go/go1.20.3.linux-amd64.tar.gz https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
    RUN tar -C /usr/local -xzf ./out/dl/go/go1.20.3.linux-amd64.tar.gz
    ENV PATH /usr/local/go/bin:$PATH

    COPY go.mod go.sum .

    RUN go mod download

    COPY . .

    RUN go build $BUILD_FLAGS ./...
    RUN go test $BUILD_FLAGS -cover ./...

    RUN mkdir -p ./out/linux-x64
    RUN go build $BUILD_FLAGS -o ./out/linux-x64/tob-cob
    RUN mkdir -p ./out/linux-x64/assets
    RUN cp -r ./assets/* ./out/linux-x64/assets

    SAVE ARTIFACT ./out/linux-x64/ /linux-x64

butler:
    FROM ubuntu:lunar-20230314

    WORKDIR /work

    RUN apt-get update -qq
    RUN apt-get install -qy curl unzip

    RUN mkdir -p ./out/dl
    RUN curl -L -o ./out/dl/butler.zip https://broth.itch.ovh/butler/linux-amd64/15.21.0/archive/default

    RUN mkdir -p ./out/butler
    RUN unzip ./out/dl/butler.zip -d ./out/butler
    RUN chmod +x ./out/butler/butler

    ENV PATH /work/out/butler:$PATH

deploy:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +linux-x64/linux-x64 ./out/linux-x64

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out/linux-x64 szabba/tears-of-butterflies-colors-of-blood:linux-x64

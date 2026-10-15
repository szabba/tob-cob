VERSION 0.7

ARG REVISION

assets:
    FROM scratch
    COPY ./assets .
    SAVE ARTIFACT ./* /

linux-x64:
    RUN GOOS=linux GOARCH=amd64 go build -o tob-cob
    SAVE ARTIFACT ./tob-cob /linux-x64/tob-cob

linux-arm64:
    RUN GOOS=linux GOARCH=arm64 go build -o tob-cob
    SAVE ARTIFACT ./tob-cob /linux-arm64/tob-cob

osx-arm64:
    RUN GOOS=darwin GOARCH=arm64 go build -o tob-cob
    SAVE ARTIFACT ./tob-cob /osx-arm64/tob-cob

windows-x64:
    RUN GOOS=darwin GOARCH=arm64 go build -o tob-cob
    SAVE ARTIFACT ./tob-cob /windows-x64/tob-cob

windows-arm64:
    RUN GOOS=darwin GOARCH=arm64 go build -o tob-cob
    SAVE ARTIFACT ./tob-cob /windows-arm64/tob-cob

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

deploy-linux-x64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +linux-x64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:linux-x64

deploy-linux-arm64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +linux-arm64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:linux-arm64

deploy-osx-arm64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +osx-arm64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:osx-arm64


deploy-windows-x64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +windows-x64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:windows-x64

deploy-windows-arm64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +windows-arm64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:windows-arm64

deploy:
    WAIT
        BUILD +deploy-linux-x64
        BUILD +deploy-linux-arm64
        BUILD +deploy-osx-arm64
        BUILD +deploy-windows-x64
        BUILD +deploy-windows-arm64
    END

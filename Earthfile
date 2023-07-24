VERSION 0.7

IMPORT github.com/prelift/earthly-udcs/go:test-revisions

assets:
    FROM scratch
    COPY ./assets .
    SAVE ARTIFACT ./* /

module:
    DO go+MODULE --BUILDER_TAG=1.19 \
        --TIDY=false # Go 1.19 did not support the -x flag in go tidy

    DO go+MODULE --BUILDER_TAG=1.20

    SAVE IMAGE module:current

linux-x64:
    DO go+BINARY --BUILDER=module --BUILDER_TAG=current \
        --REVISION=$EARTHLY_TARGET_TAG \
        --PACKAGE=github.com/szabba/tob-cob --OUTPUT=tob-cob \
        --CGO_ENABLED=1 \
        --GOOS=linux --GOARCH=amd64

    SAVE ARTIFACT ./tob-cob /linux-x64/tob-cob

osx-arm64:
    DO go+BINARY --BUILDER=module --BUILDER_TAG=current \
        --REVISION=$EARTHLY_TARGET_TAG \
        --PACKAGE=github.com/szabba/tob-cob --OUTPUT=tob-cob \
        --CGO_ENABLED=1 \
        --GOOS=darwin --GOARCH=arm64
    
    SAVE ARTIFACT ./tob-cob /osx-arm64/tob-cob

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

deploy-osx-arm64:
    FROM +butler
    WORKDIR /work

    RUN mkdir -p ~/.config/itch

    COPY +assets/ ./out/assets
    COPY +osx-arm64/tob-cob ./out

    RUN butler -V
    RUN --push --secret BUTLER_API_KEY \
        butler push ./out szabba/tears-of-butterflies-colors-of-blood:osx-arm64

deploy:
    WAIT
        BUILD +deploy-linux-x64
        BUILD +deploy-osx-arm64
    END

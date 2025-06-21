
FROM docker.io/golang:1.24.4-bookworm AS build
WORKDIR /go/src/app
COPY ../../go.mod ../../go.sum ../../Makefile ./
RUN make deps
COPY ../.. ./
RUN make bin/main

FROM scratch AS fonts
ADD https://github.com/notofonts/notofonts.github.io/raw/refs/heads/main/fonts/NotoSans/full/ttf/NotoSans-Bold.ttf /usr/share/fonts/truetype/noto/

FROM scratch
WORKDIR /app
COPY --from=fonts /usr/share/fonts/truetype/noto/NotoSans-Bold.ttf /usr/share/fonts/truetype/noto/
# zone db is probably better mounted from the host than bundled with the image...
# COPY --from=build /usr/share/zoneinfo/Europe/London /usr/share/zoneinfo/Europe/London
COPY --from=build /go/src/app/bin/main ./
CMD ["/app/main"]

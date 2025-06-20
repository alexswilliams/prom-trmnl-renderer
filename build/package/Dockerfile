
FROM docker.io/golang:1.24.2-bookworm AS build
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
COPY --from=build /usr/share/zoneinfo/Europe/London /usr/share/zoneinfo/Europe/London
COPY --from=build /go/src/app/bin/main ./
CMD ["/app/main"]

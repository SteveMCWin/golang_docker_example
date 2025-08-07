FROM golang:1.24

RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev \
    build-essential \
    file \
    libpcre3-dev \
    libpcre3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN wget -O spellfix.c https://raw.githubusercontent.com/sqlite/sqlite/master/ext/misc/spellfix.c
RUN mkdir -p extensions
RUN gcc -fPIC -shared -o extensions/spellfix.so spellfix.c -lsqlite3 -lpcre
RUN rm spellfix.c

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /myapp

EXPOSE 8080

CMD ["/myapp"]

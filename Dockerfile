FROM golang:1.24

# get dependencies
RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-dev \
    build-essential \
    file \
    libpcre3-dev \
    libpcre3 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# get the sqlite extension
RUN wget -O spellfix.c https://raw.githubusercontent.com/sqlite/sqlite/master/ext/misc/spellfix.c
RUN mkdir -p extensions
RUN gcc -fPIC -shared -o extensions/spellfix.so spellfix.c -lsqlite3 -lpcre
RUN rm spellfix.c

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# set up the database
RUN touch data/person.db
RUN sqlite3 data/person.db < data/people.sql
RUN sqlite3 data/person.db < data/spellfix_people.sql

RUN CGO_ENABLED=1 GOOS=linux go build -o myapp

EXPOSE 8080

CMD ["./myapp"]

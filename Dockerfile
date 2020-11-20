#dockerイメージ
FROM golang:latest

#コンテナ内に作業ディレクトリを
RUN mkdir /app

#コンテナログイン時のディレクトリ指定
WORKDIR /app

#RUNを書いてからbuildしないとダメっぽい（？）
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/jinzhu/gorm
RUN go get github.com/joho/godotenv
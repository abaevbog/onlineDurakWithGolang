FROM golang:1.12


WORKDIR /golang/src/github/abaevbog/onlineDurakWithGolang

ENV GOPATH "/Users/bogdanabaev/RandomProgramming/golang/"

COPY . .

RUN go get -d -v 
RUN go install -v 

CMD ["./main"]

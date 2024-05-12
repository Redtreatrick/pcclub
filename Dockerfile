FROM scratch

COPY . .

#ADD main test_file.txt /

#RUN go build -o main .

CMD ["./main", "test_file.txt"]
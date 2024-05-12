FROM scratch

COPY main test_file.txt /

CMD ["./main", "test_file.txt"]
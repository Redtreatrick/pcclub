Test Assignment for YaDRO

to launch (only tested on linux) you follow these steps:

1) update modules using go mod tidy (use go mod init before if needed)
2) build an application using go build -o {binary_name}
3) put {test_file_txt_name} in directory
4) change main and test_file.txt to {binary_name} and {test_file_txt_name}
5) build docker container using docker build -t {docker_container_name} .
6) run docker container using docker run -it {docker_container_name}

FROM alpine:3.7
COPY go_all_demo_back /usr/local/

RUN ["chmod", "+x", "usr/local/go_all_demo_back"]
CMD usr/local/go_all_demo_back


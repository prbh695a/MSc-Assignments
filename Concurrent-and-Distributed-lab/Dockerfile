FROM tudinfse/cds_server
#FROM golang:latest 

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y build-essential vim golang-go

ADD cds_server.json /etc/cds_server.json
ADD mopp-2018-t0-harmonic-progression-sum/ mopp-2018-t0-harmonic-progression-sum/
ADD mopp-2018-t1-average mopp-2018-t1-average/
ADD mopp-2018-t2-levenshtein mopp-2018-t2-levenshtein/

WORKDIR /mopp-2018-t0-harmonic-progression-sum
RUN make

WORKDIR /mopp-2018-t1-average
RUN make

WORKDIR /mopp-2018-t2-levenshtein
RUN make

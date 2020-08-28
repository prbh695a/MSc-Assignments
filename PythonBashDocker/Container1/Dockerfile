FROM ubuntu:16.04
      
MAINTAINER Your name <your.name@mailbox.tu-dresden.de>

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install -y python make autoconf build-essential openssh-server

RUN mkdir /var/run/sshd
RUN echo 'root:screencast' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
CMD ["/usr/sbin/sshd", "-D"]

RUN adduser --disabled-password --gecos "" ubuntu
WORKDIR "/home/ubuntu"
COPY id_rsa.pub /home/ubuntu/.ssh/authorized_keys


WORKDIR "/home/ubuntu"
RUN wget http://monkey.org/~provos/libevent-1.3e.tar.gz
RUN tar zxvf libevent-1.3e.tar.gz
WORKDIR "libevent-1.3e"
RUN ./configure
RUN make
RUN make install
RUN ldconfig -v

WORKDIR "/home/ubuntu"
RUN wget http://www.memcached.org/files/memcached-1.4.33.tar.gz
RUN tar zxvf memcached-1.4.33.tar.gz
WORKDIR "memcached-1.4.33"
RUN ./configure
RUN make
RUN make install

EXPOSE 11211 22

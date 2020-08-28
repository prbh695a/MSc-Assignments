FROM ubuntu:16.04
      
MAINTAINER Your name <your.name@mailbox.tu-dresden.de>

RUN apt-get update
RUN apt-get upgrade -y

RUN apt-get install -y openssh-server python mercurial git make autoconf build-essential gcc g++ r-base-core

RUN mkdir /var/run/sshd
RUN echo 'root:screencast' | chpasswd
RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed -i 's/StrictModes yes/StrictModes no/' /etc/ssh/sshd_config
# SSH login fix. Otherwise user is kicked off after login
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
ENV NOTVISIBLE "in users profile"
RUN echo "export VISIBLE=now" >> /etc/profile
RUN adduser --disabled-password --gecos "" ubuntu

#COPY id_rsa.pub /root/.ssh/authorized_keys
COPY id_rsa.pub /home/ubuntu/.ssh/authorized_keys
COPY id_rsa /home/ubuntu/

WORKDIR "/home/ubuntu"
RUN hg clone https://bitbucket.org/db7/dude
WORKDIR "dude"
RUN python setup.py install

WORKDIR "/home/ubuntu"
RUN git clone git://github.com/twitter/twemperf.git
WORKDIR "twemperf"
RUN autoreconf -fvi
RUN ./configure
RUN make
RUN make install

WORKDIR "/home/ubuntu"
COPY lattice_0.20-38.tar.gz .
COPY Benchmark.py .
COPY Dudefile .
COPY graphs.R .
COPY run.sh .
RUN chmod 0777 *

RUN gunzip -c lattice_0.20-38.tar.gz | tar xvf -
RUN rm -f lattice_0.20-38.tar.gz
WORKDIR "/home/ubuntu"


ENV PATH="${PATH}:$HOME/usr/local/bin"

CMD ["/usr/sbin/sshd", "-D"]
EXPOSE 22


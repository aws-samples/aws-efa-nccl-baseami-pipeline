FROM nvidia/cuda:10.2-devel

ENV USER root

# -------------------------------------------------------------------------------------
# install needed software -
# openssh
# mpi
# awscli
# supervisor
# -------------------------------------------------------------------------------------

RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt install -y unzip iproute2 nano cmake openssh-server openssh-client python3 python3-pip build-essential wget curl pigz && apt clean
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip && ./aws/install

RUN mkdir -p /var/run/sshd
ENV DEBIAN_FRONTEND noninteractive

ENV NOTVISIBLE "in users profile"

#####################################################
## SSH SETUP

RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd
RUN echo "export VISIBLE=now" >> /etc/profile

RUN echo "${USER} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
ENV SSHDIR /root/.ssh
RUN mkdir -p ${SSHDIR} && \
    touch ${SSHDIR}/sshd_config && \
    ssh-keygen -t rsa -f ${SSHDIR}/ssh_host_rsa_key -N '' && \
    cp ${SSHDIR}/ssh_host_rsa_key.pub ${SSHDIR}/authorized_keys && \
    cp ${SSHDIR}/ssh_host_rsa_key ${SSHDIR}/id_rsa && \
    echo "    IdentityFile ${SSHDIR}/id_rsa" >> /etc/ssh/ssh_config && \
    echo "Host *" >> /etc/ssh/ssh_config && echo "    StrictHostKeyChecking no" >> /etc/ssh/ssh_config && \
    chmod -R 600 ${SSHDIR}/* && \
    chown -R ${USER}:${USER} ${SSHDIR}/
# check if ssh agent is running or not, if not, run
RUN eval `ssh-agent -s` && ssh-add ${SSHDIR}/id_rsa

##################################################
## CUDA MPI

RUN wget -O /tmp/openmpi.tar.gz https://download.open-mpi.org/release/open-mpi/v4.0/openmpi-4.0.3.tar.gz && \
    tar -xvf /tmp/openmpi.tar.gz -C /tmp && cd /tmp/openmpi* && ./configure --prefix=/opt/openmpi --with-cuda --enable-mpirun-prefix-by-default && \
    make -j $(nproc) && make install
RUN echo "export PATH=$PATH:/opt/openmpi/bin:/usr/local/cuda/bin" >> /etc/profile
RUN echo "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/openmpi/lib:/usr/local/cuda/include:/usr/local/cuda/lib64" >> /etc/profile

###################################################
## APPLICATION LAYER
RUN /usr/local/cuda/bin/nvcc --version

###################################################
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D", "-p", "22"]

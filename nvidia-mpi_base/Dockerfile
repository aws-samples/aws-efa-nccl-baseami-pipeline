FROM nvidia/cuda:11.0-devel-ubuntu18.04

ARG EFA_INSTALLER_VERSION=latest
ARG NCCL_VERSION=master
ARG AWS_OFI_NCCL_VERSION=aws
ARG NCCL_TESTS_VERSION=v2.0.0

RUN apt-get update -y
RUN apt-get install -y --allow-unauthenticated \
    git \
    gcc \
    yum-utils \
    vim \
    kmod \
    openssh-client \
    openssh-server \
    build-essential \
    curl \
    autoconf \
    libtool \
    gdb \
    automake


ENV HOME /tmp

ENV LD_LIBRARY_PATH=/usr/local/cuda/extras/CUPTI/lib64:/opt/amazon/openmpi/lib:$HOME/nccl/build/lib:/opt/amazon/efa/lib:/$HOME/aws-ofi-nccl/install/lib:$LD_LIBRARY_PATH
ENV PATH=/opt/amazon/openmpi/bin/:/opt/amazon/efa/bin:$PATH


#################################################
## Install EFA installer
RUN cd $HOME \
    && curl -O https://efa-installer.amazonaws.com/aws-efa-installer-${EFA_INSTALLER_VERSION}.tar.gz \
    && tar -xf $HOME/aws-efa-installer-${EFA_INSTALLER_VERSION}.tar.gz \
    && cd aws-efa-installer \
    && ./efa_installer.sh -y -d --skip-kmod --skip-limit-conf --no-verify | tee $HOME/install.log

## Test EFA installation
RUN /opt/amazon/efa/bin/fi_info --version

###################################################
## Install NCCL
RUN git clone https://github.com/NVIDIA/nccl.git $HOME/nccl \
    && cd $HOME/nccl \
    && make -j src.build CUDA_HOME=/usr/local/cuda \
    NVCC_GENCODE="-gencode=arch=compute_80,code=sm_80 -gencode=arch=compute_75,code=sm_75 -gencode=arch=compute_70,code=sm_70"

###################################################
## Install AWS-OFI-NCCL plugin
RUN git clone https://github.com/aws/aws-ofi-nccl.git $HOME/aws-ofi-nccl \
    && cd $HOME/aws-ofi-nccl \
    && git checkout ${AWS_OFI_NCCL_VERSION} \
    && ./autogen.sh \
    && ./configure --prefix=$HOME/aws-ofi-nccl/install \
       --with-libfabric=/opt/amazon/efa/ \
       --with-cuda=/usr/local/cuda \
       --with-nccl=$HOME/nccl/build \
       --with-mpi=/opt/amazon/openmpi/ \
    && make && make install

###################################################
## Install NCCL-tests
RUN git clone https://github.com/NVIDIA/nccl-tests.git $HOME/nccl-tests \
    && cd $HOME/nccl-tests \
    && git checkout ${NCCL_TESTS_VERSION} \
    && make MPI=1 \
       MPI_HOME=/opt/amazon/openmpi/ \
       CUDA_HOME=/usr/local/cuda \
       NCCL_HOME=$HOME/nccl/build \
       NVCC_GENCODE="-gencode=arch=compute_80,code=sm_80 -gencode=arch=compute_75,code=sm_75 -gencode=arch=compute_70,code=sm_70"

CMD echo "EFA Setup Complete! "

FROM phusion/baseimage:focal-1.0.0

ENV DEBIAN_FRONTEND=noninteractive
RUN locale-gen en_US.UTF-8

ENV LANGUAGE=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8
ENV LC_CTYPE=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV TERM xterm


RUN set -eux \
    # Add the "PHP 8" ppa
    && apt-get install -y software-properties-common \
    && add-apt-repository -y ppa:ondrej/php \
    #
    #--------------------------------------------------------------------------
    # Software's Installation
    #--------------------------------------------------------------------------
    #
    && echo 'DPkg::options { "--force-confdef"; };' >> /etc/apt/apt.conf \
    # Install "PHP Extentions", "libraries", "Software's"
    && apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --allow-downgrades --allow-remove-essential \
        --allow-change-held-packages \
        php8.1-cli \
        php8.1-common \
        php8.1-curl \
        php8.1-intl \
        php8.1-xml \
        php8.1-mbstring \
        php8.1-mysql \
        php8.1-pgsql \
        php8.1-sqlite \
        php8.1-sqlite3 \
        php8.1-zip \
        php8.1-bcmath \
        php8.1-memcached \
        php8.1-gd \
        php8.1-dev \
        pkg-config \
        libcurl4-openssl-dev \
        libedit-dev \
        libssl-dev \
        libxml2-dev \
        xz-utils \
        libsqlite3-dev \
        rsync \
        sqlite3 \
        git \
        curl \
        vim \
        nano \
        tree \
    && apt-get clean \
    #####################################
    # Composer:
    #####################################
    # Install composer and add its bin to the PATH.
    && curl -s http://getcomposer.org/installer | php \
    && echo "export PATH=${PATH}:/var/www/vendor/bin" >> ~/.bashrc \
    && mv composer.phar /usr/local/bin/composer

USER root

ARG PUID=1000
ARG PGID=1000


# always run apt update when start and after add new source list, then clean up at end.
RUN set -xe; \
    apt-get update -yqq \
    && pecl channel-update pecl.php.net \
    && groupadd -g ${PGID} laradock \
    && useradd -l -u ${PUID} -g laradock -m laradock -G docker_env \
    && usermod -p "*" laradock -s /bin/bash \
    && apt-get install -yqq apt-utils libzip-dev zip unzip php8.1-zip nasm \
    && php -m | grep -q 'zip' \
    && ln -snf /usr/share/zoneinfo/UTC /etc/localtime && echo UTC > /etc/timezone \
    && composer self-update --2 \
    && apt-get install -yqq \
        php8.1-redis \
        php8.1-igbinary \
        libxrender1 \
        libfontconfig1 \
        libx11-dev \
        libjpeg62 \
        libxtst6 \
        fontconfig \
        libjpeg-turbo8-dev \
        xfonts-base \
        xfonts-75dpi \
        wget



USER root

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    rm /var/log/lastlog /var/log/faillog

# Set default work directory
WORKDIR /var/www

FROM php:8.3-fpm

LABEL MAINTAINER="vahagn99ghukasyan@gmail.com"

ARG WWW_DATA_USER
ARG WWW_DATA_GROUP
ENV TZ=Asia/Erevan

# Set working directory
WORKDIR /var/www/

# Clean up and install required packages
RUN apt-get clean && apt-get update && apt-get install -y --no-install-recommends \
    git \
    zip \
    unzip \
    libpng-dev \
    libjpeg-dev \
    libfreetype6-dev \
    libonig-dev \
    libxml2-dev \
    curl \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) gd mbstring pdo pdo_mysql xml bcmath \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Install composer
RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" \
    && php composer-setup.php \
    && php -r "unlink('composer-setup.php');" \
    && mv composer.phar /usr/local/bin/composer

# Add user for app
RUN groupadd -g ${WWW_DATA_GROUP:-1000} www \
    && useradd -u ${WWW_DATA_USER:-1000} -ms /bin/bash -g www www

# Copy existing application directory contents
COPY . /var/www

# Copy existing application directory permissions
COPY --chown=www:www . /var/www

# Change current user to www
USER www

# Expose port 9000
EXPOSE 9000

# Start PHP-FPM
CMD ["php-fpm"]

FROM php:8.2-apache

COPY --from=composer/composer /usr/bin/composer /usr/bin/composer

RUN curl -1sLf 'https://dl.cloudsmith.io/public/symfony/stable/setup.deb.sh' | bash
RUN apt install unzip -y
RUN pecl install xdebug
RUN docker-php-ext-install -j$(nproc) pdo pdo_mysql
RUN docker-php-ext-enable xdebug

RUN usermod -u 1000 www-data
RUN chown -R www-data:www-data /var/www
RUN a2enmod rewrite

FROM php:8.2-apache

COPY --from=composer/composer /usr/bin/composer /usr/bin/composer

RUN pecl install xdebug
RUN docker-php-ext-install -j$(nproc) pdo pdo_mysql
RUN docker-php-ext-enable xdebug

RUN usermod -u 1000 www-data
RUN chown -R www-data:www-data /var/www
RUN a2enmod rewrite

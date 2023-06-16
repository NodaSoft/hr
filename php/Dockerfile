
FROM php:8.2-cli-alpine

ARG UID=600
ARG GID=600

RUN apk --update add --no-cache \
    $PHPIZE_DEPS \
    shadow \
    && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
    && docker-php-ext-install -j $("nproc") pdo_mysql \
    && apk del $PHPIZE_DEPS \
    && rm -rf /tmp/* /var/tmp/* /usr/share/doc/* /var/cache/apk/* /usr/share/php8 \
    && usermod -u $UID www-data -s /bin/bash && groupmod -g $GID www-data \
    && chown -R www-data:www-data /var/www

USER www-data
WORKDIR /var/www/ns

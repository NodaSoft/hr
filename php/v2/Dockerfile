FROM php:7.1-cli

RUN apt-get update && apt-get install -y git unzip

COPY --from=composer:2.2 /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html

EXPOSE 8080
CMD ["php", "-S", "0.0.0.0:8080"]
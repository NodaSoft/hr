#!/bin/sh

vendor/bin/phpcbf -p --standard=PSR12 ./src/ ./tests/ --parallel=$("nproc")
vendor/bin/php-cs-fixer fix src
vendor/bin/php-cs-fixer fix tests
vendor/bin/phpcs -p --standard=PSR12 ./src/ ./tests/ --parallel=$("nproc")
vendor/bin/phpmd src/ text cleancode,codesize,design
vendor/bin/phpmd tests/ text cleancode,codesize,design
vendor/bin/phpstan analyse --level=7 --no-progress -vvv --memory-limit=512M
vendor/bin/psalm --no-cache

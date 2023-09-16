fix-code-style:
	@vendor/bin/php-cs-fixer fix --allow-risky=yes --verbose --using-cache=no

analysis-code:
	@php -d memory_limit=-1 vendor/bin/phpstan analyse -c phpstan.neon --memory-limit=-1
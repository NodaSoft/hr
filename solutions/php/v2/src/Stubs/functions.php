<?php

function __(string $name, ...$args): array
{
	return [
		'from' => "FROM: $name",
		'to'   => "TO: $name",
	];
}
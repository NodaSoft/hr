<?php

declare(strict_types=1);

namespace NW\WebService\Request;

class Request
{
    public static function input(string $key): mixed
    {
        return $_POST[$key] ?? null;
    }

    public static function get(string $key): mixed
    {
        return $_GET[$key] ?? null;
    }

    private static function clear(mixed $val): mixed
    {
        return $val;
    }

}


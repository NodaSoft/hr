<?php

namespace App\Connection\Tools;


class ConverterTools
{
    /**
     * @param string $data
     * @return array
     */
    public static function jsonDecode(string $data): array
    {
        if (!$data) {
            return [];
        }

        try {
            return json_decode($data, true, 512, JSON_THROW_ON_ERROR);
        } catch (\JsonException $e) {
            throw new \RuntimeException($e->getMessage(), $e->getCode(), $e);
        }
    }
}
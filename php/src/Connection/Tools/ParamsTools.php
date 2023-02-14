<?php

namespace App\Connection\Tools;


class ParamsTools
{
    /**
     * @param array $array
     * @return string
     */
    public static function whereArrToStr(array $array): string
    {
        if (!$array) {
            return '';
        }

        try {
            return str_repeat('?,', count($array) - 1) . '?';
        } catch (\Exception $e) {
            throw new \RuntimeException($e->getMessage(), $e->getCode(), $e);
        }
    }
}
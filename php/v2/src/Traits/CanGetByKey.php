<?php

namespace Nodasoft\Testapp\Traits;

use Exception;

trait CanGetByKey
{
    public function getBy(
        array  $arr,
        mixed  $value,
        string $key = 'id'): array|false
    {
        $item = array_filter($arr, function ($item) use ($key, $value) {
            return $item[$key] === $value;
        });

        return current($item);
    }

    /**
     * @throws Exception
     */
    public function getByKeyOrThrow(
        array  $arr,
        mixed  $value,
        string $key = 'id',
        string $message = null,
        int    $code = 400): array
    {
        $item = $this->getBy($arr, $value, $key);
        if (!$item) throw new Exception($message ?? 'record not found!', $code);
        return $item;
    }
}
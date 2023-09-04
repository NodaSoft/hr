<?php

declare(strict_types=1);

namespace pseudovendor;

/**
 * Имплементация реквеста от условного самого лучшего на свете фреймворка
 */
class Request
{
    /**
     * Вернет тело запроса, которое фреймворк предварительно конвертировал в зависимости от content-type
     *
     * @return mixed
     */
    public function getBody(): mixed
    {
        return 'stub';
    }
}

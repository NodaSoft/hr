<?php

namespace NW\WebService\References\Operations\Notification;

abstract class ReferencesOperation
{
    abstract public function doOperation(): array;

    /**
     * Работать с глобальными переменными таким образом опасно как минимум из-за SQL инъекций, переполнения буфера.
     * $_REQUEST содержит в себе параметры из $_GET, $_POST, $_COOKIE. Поэтому может возникнуть непредсказуемое поведение.
     *
     * Нужно фильтровать данные, экранировать.
     *
     * @param $pName
     * @return mixed
     */
    public function getRequest($pName)
    {
        return $_REQUEST[$pName];
    }
}
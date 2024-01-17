<?php

namespace App\v2\Facades;

abstract class Event
{
    public abstract function dispatch(...$args): void;
}
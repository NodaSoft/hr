<?php

namespace NW\WebService\References\Operations\Notification\Clients;

/**
 * Interface ClientInterface
 */
interface ClientInterface
{
    public function send(...$args): bool;
}
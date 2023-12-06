<?php

namespace NW\WebService\References\Operations\Notification\Clients;

interface ClientInterface
{
    public static function getById(int $clientId): ?static;

    public function getFullName(): string;
}
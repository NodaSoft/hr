<?php

namespace NW\WebService\References\Operations\Notification\Clients;

/**
 * Абстрактная сущность "Клиент"
 *
 * @property int $id
 * @property int $type
 * @property string $name
 */
abstract class Client implements ClientInterface
{
    /** @var int $id */
    public $id;

    /** @var int $type */
    public $type;

    /** @var string $name */
    public $name;

    public static function getById(int $clientId): ?static
    {
        return new static($clientId); // fakes the getById method
    }

    public function getFullName(): string
    {
        return $this->name;
    }
}
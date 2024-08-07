<?php

namespace NW\WebService\References\Operations\Notification\Serializers;

use JsonSerializable;
use NW\WebService\References\Operations\Notification\DTOs\IndexDTO;

class IndexSerializer implements JsonSerializable
{
    public function __construct(
        protected IndexDTO $indexDTO,
    ) {
    }

    public function jsonSerialize(): array
    {
        $dto = $this->indexDTO;
        return [
            'notificationEmployeeByEmail' => $dto->notificationEmployeeByEmail,
            'notificationClientByEmail' => $dto->notificationClientByEmail,
            'notificationClientBySms' => new SmsSerializer($dto->notificationClientBySms),
        ];
    }
}

<?php

namespace NW\WebService\References\Operations\Notification\Serializers;

use JsonSerializable;
use NW\WebService\References\Operations\Notification\DTOs\SmsDTO;

class SmsSerializer implements JsonSerializable
{
    public function __construct(
        protected SmsDTO $indexDTO,
    ) {
    }

    public function jsonSerialize(): array
    {
        $dto = $this->indexDTO;
        return [
            'isSent'  => $dto,
            'message' => $dto,
        ];
    }
}

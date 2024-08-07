<?php

namespace NW\WebService\References\Operations\Notification\models;

use NW\WebService\References\Operations\Notification\DTOs\StatusEnum;
use NW\WebService\References\Operations\Notification\Exceptions\NotFoundException;

class  Status
{
    const ERROR_STATUS_NOT_FOUND = 'Status not found';

    public function __construct(
        public int $id,
        public string $name
    ) {
    }

    /**
     * @ throws NotFoundException
     */
    public static function findById(int $id): ?self
    {
        $statusMap = StatusEnum::STATUS_MAP;
        if (!array_key_exists($id, $statusMap)) {
            return null;
            // or // throw new NotFoundException(self::ERROR_STATUS_NOT_FOUND);
        }
        return new self(
            id: $id,
            name: $statusMap[$id]
        );
    }
}
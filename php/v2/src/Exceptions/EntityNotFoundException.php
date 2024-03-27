<?php

namespace App\Exceptions;

class EntityNotFoundException extends \HttpException
{
    public function __construct(string $entityType, int $id)
    {
        parent::__construct(sprintf('Entity [%s] with id: [%s] not found', $entityType, $id), 404);
    }
}

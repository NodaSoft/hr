<?php

namespace NodaSoft\DataMapper\Mapper;

use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\EntityInterface\Entity;

class NotificationMapper implements Mapper
{
    /**
     * @param int $id
     * @return null|Notification
     */
    public function getById(int $id): ?Entity
    {
        // TODO: Implement getById() method.
    }

    /**
     * @param string $string
     * @return Notification|null
     */
    public function getByName(string $string): ?Entity
    {
        // TODO: Implement getByName() method.
    }
}

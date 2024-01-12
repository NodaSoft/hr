<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

/**
 * @property Seller $Seller
 */
class Contractor extends AbstractContractor
{
    public const TYPE_CUSTOMER = 0;
    public int $id;
    public int $type;
    public string $name;
    public string $email;
    public string $mobile;
}

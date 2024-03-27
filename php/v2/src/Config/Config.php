<?php

namespace App\Config;

class Config
{
    public const RESELLER_FROM_EMAIL = 'contractor@example.com';
    public const EVENT_CHANGE_RETURN_STATUS = 'changeReturnStatus';
    public const EVENT_NEW_RETURN_STATUS = 'newReturnStatus';

    /**
     * @param int $resellerId
     * @param string $eventName
     * @return string[]
     */
    public static function getEmailsByEventForReseller(int $resellerId, string $eventName): array
    {
        //TODO: This is fake method, do not forget to change it to something usefull in future.
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }
}

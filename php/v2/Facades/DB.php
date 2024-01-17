<?php

namespace App\v2\Facades;

use App\v2\Models\Contractors\Contractor;

class DB
{

    private const string RESELLER_EMAIL_FROM = 'contractor@example.com';

    /**
     * @return string
     */
    static function getResellerEmailFrom(): string
    {
        try {
            return self::RESELLER_EMAIL_FROM;
        } catch (\Exception $ex) {
            return '';
        }
    }

    /**
     * @param $resellerId
     * @return string[]
     */
    static function getEmailsByPermit($resellerId): array
    {
        try {
            return ['someemeil@example.com', 'someemeil2@example.com']; //fake search in db by id, and get emails
        } catch (\Exception $exception) {
            return [];
        }
    }

    public static function getClientById(int $clientId): Contractor
    {
        return new Contractor();
    }
}
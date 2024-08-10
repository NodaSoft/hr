<?php

namespace app\Domain\Notification\DTO;

use Spatie\DataTransferObject\DataTransferObject;

class NotificationData extends DataTransferObject
{
    /** @var int */
    public $complaint_id;

    /** @var string */
    public $complaint_number;

    /** @var int */
    public $creator_id;

    /** @var int */
    public $reseller_id;

    /** @var int */
    public $expert_id;

    /** @var int */
    public $client_id;

    /** @var int */
    public $consumption_id;

    /** @var string */
    public $consumption_number;

    /** @var string */
    public $agreement_number;

    /** @var string */
    public $date;

    /** @var int */
    public $notification_type;

    /** @var array */
    public $differences;

    public static function fromNotificationRequest(array $request): NotificationData
    {

        $differences = [];

        if($request['differences']){
            $differences = [
                'FROM' => (int)$request['differences']['from'] ?? null,
                'TO' => (int)$request['differences']['to'] ?? null
            ];
        }

        $data = [
            'complaint_id' => (int)$request['complaintId'],
            'complaint_number' => (string)$request['complaintNumber'],
            'reseller_id' => (int)$request['resellerId'] ?? null,
            'creator_id' => (int)$request['creatorId'],
            'expert_id' => (int)$request['expertId'],
            'client_id' => (int)$request['clientId'],
            'consumption_id' => (int)$request['consumptionId'],
            'consumption_number' => (string)$request['consumptionNumber'],
            'agreement_number' => (string)$request['agreementNumber'],
            'date' => (string)$request['date'],
            'notification_type' => (int)$request['notificationType'] ?? null,
            'differences' => $differences,
        ];

        return new self($data);
    }
}
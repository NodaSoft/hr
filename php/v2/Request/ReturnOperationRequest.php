<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Request\RequestAbstract;

readonly class ReturnOperationRequest extends RequestAbstract
{


    public function __construct(
        public int $complaintId,
        public int $resellerId,
        public int $complaintNumber,
        public int $creatorId,
        public int $expertId,
        public int $clientId,
        public int $consumptionId,
        public string $consumptionNumber,
        public string $agreementNumber,
        public string $date,
        public int $notificationType,
        public ?array $differences,
    ) {
    }

    protected static function rules(): array
    {
        return [
            'complaintId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'resellerId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'complaintNumber' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'creatorId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'expertId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'clientId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'consumptionId' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'consumptionNumber' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'agreementNumber' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'date' => ['required' => true, 'params' => ['filter' => FILTER_VALIDATE_INT]],
            'notificationType' => [
                'required' => true,
                'params' => ['filter' => FILTER_CALLBACK, 'options' => new EnumRule(NotificationTypeEnum::class)],
            ],
            'differences' => [
                'required' => true,
                'flags' => FILTER_REQUIRE_ARRAY,
                'params' => ['filter' => FILTER_CALLBACK, 'options' => new EnumRule(StatusEnum::class)],
            ]
        ];
    }
}
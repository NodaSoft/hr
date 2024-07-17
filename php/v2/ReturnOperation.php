<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Enums\NotificationTypeEnum;
use NW\WebService\References\Operations\Notification\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\Validation\ValidationPipeline;

/**
 * Class ReturnOperation
 *
 * This class processes return operations and sends notifications to employees and clients.
 */
readonly class ReturnOperation
{
    public function __construct(
        private ValidationPipeline  $validator,
        private NotificationManager $notificationManager,
    )
    {
    }

    /**
     * Processes the return operation and sends notifications to employees and clients.
     * @throws ValidationException
     */
    public function process(NotificationData $data): array
    {
        $this->validator->process($data);

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $this->notificationManager->sendEmployeeNotifications($data,  $result);

        //  Шлём клиентское уведомление, только если произошла смена статуса
        if ($data->notificationType === NotificationTypeEnum::CHANGE && !empty($data->differences['to'])) {
            $this->notificationManager->sendClientNotifications($data,  $result);
        }

        return $result;
    }
}

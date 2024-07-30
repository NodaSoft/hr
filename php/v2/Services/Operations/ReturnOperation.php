<?php

namespace NW\WebService\References\Operations\Notification\Services\Operations;

use Exception;
use NW\WebService\References\Operations\Notification\Models\Contractor;
use NW\WebService\References\Operations\Notification\Models\Employee;
use NW\WebService\References\Operations\Notification\Models\MessageTypes;
use NW\WebService\References\Operations\Notification\Models\NotificationEvents;
use NW\WebService\References\Operations\Notification\Models\Status;
use NW\WebService\References\Operations\Notification\Services\Notifications\MessagesClient;
use NW\WebService\References\Operations\Notification\Services\Notifications\NotificationManager;
use NW\WebService\References\Operations\Notification\Services\OperationValidator;

/**
 * Class ReturnOperation
 * @package NW\WebService\References\Operations\Notification\Services\Operations;
 *
 * Код предназначен для осуществления уведомлений о возврате товара.
 *
 * Что было измененено:
 * 1. Добавлена типизация для свойств и методов
 * 2. Классы распределены по отдельным файлам и сгруппированы по папкам по смыслу
 * 3. Документация PHPDoc
 * 4. Добавлены константы и некоторые фейковые методы
 * 5. Добавлена валидация входящих параметров
 * 6. Переведены комментарии
 *
 */
class ReturnOperation extends ReferencesOperation
{
    public const REQUEST_DATA = 'data';

    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @return array
     */
    private function initiateResult(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }

    /**
     * @return array
     * @throws Exception
     */
    public function doOperation(): array
    {
        $result = $this->initiateResult();
        $data = $this->getRequest(self::REQUEST_DATA);

        $validator = OperationValidator::make($data);
        $error = $validator->validate();
        $customNotification = $validator->getNotificationMessage();

        if (!empty($customNotification)) {
            $result['notificationClientBySms']['message'] = $customNotification;
            return $result;
        }

        if (!empty($error)) {
            throw new Exception($error, 400);
        }

        $resellerId = $data['resellerId'];
        $notificationType = (int)$data['notificationType'];
        $client = Contractor::getById((int)$data['clientId']);

        $templateData = $this->prepareTemplate($data, $client, $notificationType, $resellerId);

        $this->notifyStaff($resellerId, $templateData, $result);

        // Send Client's notifications on status changed
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $this->notifyClient($client, $resellerId, $templateData, $result, $data);
        }

        return $result;
    }

    /**
     * @param int $resellerId
     * @param array $templateData
     * @param array $result
     * @return void
     */
    private function notifyStaff(int $resellerId, array $templateData, array &$result): void
    {
        // Get staff emails from settings
        $emails = $this->getEmailsByPermit($resellerId, NotificationEvents::EVENT_GOODS_RETURNS);
        $emailFrom = $this->getResellerEmailFrom($resellerId);
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    MessageTypes::EMAIL => [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;
            }
        }
    }

    /**
     * @param Contractor $client
     * @param int $resellerId
     * @param array $templateData
     * @param array $result
     * @param array $data
     * @return void
     */
    private function notifyClient(Contractor $client, int $resellerId, array $templateData, array &$result, array $data): void
    {
        $emailFrom = $this->getResellerEmailFrom($resellerId);

        if (!empty($emailFrom) && !empty($client->email)) {
            MessagesClient::sendMessage([
                MessageTypes::EMAIL => [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->email,
                    'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
            $result['notificationClientByEmail'] = true;
        }

        if (!empty($client->mobile)) {
            $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($error)) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }

    /**
     * @param int $notificationType
     * @param int $resellerId
     * @return string
     */
    private function getDiff(int $notificationType, int $resellerId): string
    {
        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        return $differences;
    }

    /**
     * @param array $data
     * @param Contractor $client
     * @param int $notificationType
     * @param int $resellerId
     * @return array
     * @throws Exception
     */
    private function prepareTemplate(array $data, Contractor $client, int $notificationType, int $resellerId): array
    {
        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->name;
        }

        $cr = Employee::getById((int)$data['creatorId']);
        $et = Employee::getById((int)$data['expertId']);

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => (int)$data['creatorId'],
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $cFullName,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $this->getDiff($notificationType, $resellerId),
        ];

        // Do not notify if any of template datum is empty
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ($key) is empty!", 500);
            }
        }

        return $templateData;
    }

    /**
     * @return string
     */
    public function getResellerEmailFrom(): string
    {
        return 'contractor@example.com';
    }

    /**
     * @param int $resellerId
     * @param string $event
     * @return string[]
     */
    public function getEmailsByPermit(int $resellerId, string $event): array
    {
        // fakes the method
        return ['someemeil@example.com', 'someemeil2@example.com'];
    }
}
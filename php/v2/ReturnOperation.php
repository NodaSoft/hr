<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use Exception;
use NW\WebService\References\Operations\Notification\Exceptions\DifferencesNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\TemplateKeyException;

class ReturnOperation extends ReferencesOperation
{
    public const int TYPE_NEW = 1;
    public const int TYPE_CHANGE = 2;

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $data = (array) $this->getRequest('data');

        $returnOperationDTO = new ReturnOperationDTO($data);

        $notificationType = $returnOperationDTO->notificationType;
        $client = $returnOperationDTO->client;
        $creator = $returnOperationDTO->creator;
        $expert = $returnOperationDTO->expert;
        $reseller = $returnOperationDTO->reseller;

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $reseller->id);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) $data['differences']['from']),
                'TO' => Status::getName((int) $data['differences']['to']),
            ], $reseller->id);
        }

        if ($differences === '') {
            throw new DifferencesNotFoundException('Differences are empty');
        }

        $templateData = [
            'COMPLAINT_ID' => (int) $data['complaintId'],
            'COMPLAINT_NUMBER' => (string) $data['complaintNumber'],
            'CREATOR_ID' => $creator->id,
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $expert->id,
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_ID' => $client->id,
            'CLIENT_NAME' => $client->getFullName(),
            'CONSUMPTION_ID' => (int) $data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string) $data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string) $data['agreementNumber'],
            'DATE' => (string) $data['date'],
            'DIFFERENCES' => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new TemplateKeyException("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom();
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($reseller->id, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $this->sendEmailNotification($emailFrom, $email, $templateData, $reseller);
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                $this->sendEmailNotification($emailFrom, $client->email, $templateData, $reseller, $client->id,
                    (int) $data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = '';
                $res = $this->sendClientMobileNotification($reseller, $client, $templateData,
                    (int) $data['differences']['to'], $error);
                if ($res) {
                    $result['notificationClientBySms']['isSent'] = true;
                } else {
                    $result['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $result;
    }

    /**
     * @param  string  $emailFrom
     * @param  string  $email
     * @param  array  $templateData
     * @param  Seller  $reseller
     * @param  int|null  $clientId
     * @param  int|null  $status
     * @return void
     */
    private function sendEmailNotification(
        string $emailFrom,
        string $email,
        array $templateData,
        Seller $reseller,
        int $clientId = null,
        int $status = null
    ): void {
        MessagesClient::sendMessage(
            parameters: [
                0 => [
                    'emailFrom' => $emailFrom,
                    'emailTo' => $email,
                    'subject' => __('complaintClientEmailSubject', $templateData, $reseller->id),
                    'message' => __('complaintClientEmailBody', $templateData, $reseller->id),
                ],
            ],
            resellerId: $reseller->id,
            notificationEvent: NotificationEvents::CHANGE_RETURN_STATUS,
            clientId: $clientId,
            status: $status
        );
    }

    private function sendClientMobileNotification(
        Seller $reseller,
        Contractor $client,
        array $templateData,
        int $status,
        string $error
    ) {
        return NotificationManager::send(
            resellerId: $reseller->id,
            id: $client->id,
            notificationEvent: NotificationEvents::CHANGE_RETURN_STATUS,
            status: $status,
            templateData: $templateData,
            error: $error
        );
    }
}

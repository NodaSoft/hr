<?php

namespace app\Domain\Notification\Actions;

require_once 'others.php';

use app\Domain\Notification\DTO\DifferencesDTO;
use app\Domain\Notification\Exceptions\ActionException;
use app\Domain\Notification\Gateways\EmployeeGateway;
use app\Domain\Notification\Models\NotificationEvents;
use app\Domain\Notification\Models\NotificationType;
use app\Domain\Notification\Models\Status;
use app\Domain\Notification\DTO\NotificationData;
use app\Services\EmailSender\DTO\EmailDTO;
use app\Services\EmailSender\EmailSender;
use app\Services\SmsSender\DTO\SmsDTO;
use app\Services\SmsSender\SmsSender;
use function NW\WebService\References\Operations\Notification\getEmailsByPermit;
use function NW\WebService\References\Operations\Notification\getResellerEmailFrom;

class NotificationAction
{
    /**
     * @param NotificationData $data
     * @return array
     * @throws ActionException
     */
    public function notify(NotificationData $data): array
    {
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        if (empty($data->reseller_id)) {
            throw new DataTransferObjectError('Empty resellerId', 500);
        }

        if (empty($data->notification_type)) {
            throw new DataTransferObjectError('Empty notificationType', 500);
        }

        if (empty($data->differences['from']) || empty($data->differences['to'])) {
            throw new DataTransferObjectError('Empty differences', 500);
        }

        $reseller = (new EmployeeGateway)->getEmployeeById($data->reseller_id);
        if (!$reseller) {
            throw new ActionException('Reseller not found!', 400);
        }

        $creator = (new EmployeeGateway)->getEmployeeById($data->creator_id);
        if (!$creator) {
            throw new ActionException('Creator not found!', 400);
        }

        $et = (new EmployeeGateway)->getEmployeeById($data->expert_id);
        if (!$et) {
            throw new ActionException('Expert not found!', 400);
        }

        $client = (new EmployeeGateway())->getEmployeeByIdAndByType($data->client_id);
        if (!$client) {
            throw new ActionException('Client not found!', 400);
        }

        $templateData = [
            'complaint_id' => $data->complaint_id,
            'complaint_number' => $data->complaint_number,
            'creator_id' => $data->creator_id,
            'creator_name' => $creator->name,
            'expert_id' => $data->expert_id,
            'expert_name' => $et->name,
            'client_id' => $data->client_id,
            'client_name' => $client->name,
            'consumption_id' => $data->consumption_id,
            'consumption_number' => $data->consumption_number,
            'agreement_number' => $data->agreement_number,
            'date' => $data->date,
            'differences' => $this->getDifferences(
                new DifferencesDTO([
                    'notification_type' => $data->notification_type,
                    'user_id' => $reseller->id,
                    'differences' => $data->differences
                ])),
        ];

        $emailFrom = getResellerEmailFrom();
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit();

        if ($emailFrom && count($emails) > 0) {

            foreach ($emails as $email) {
                $isSent = (new EmailSender())->sendEmail(new EmailDTO([
                    'email_from' => $emailFrom,
                    'email_to' => $email,
                    'data' => $templateData,
                    'status' => NotificationEvents::CHANGE_RETURN_STATUS,
                    'user_id' => $reseller->id
                ]));
                $result['notificationEmployeeByEmail'] = $isSent;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($data->notification_type === NotificationType::TYPE_CHANGE && !empty($data->differences['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                $isSentEmailClient = (new EmailSender())->sendEmail(new EmailDTO([
                    'email_from' => $emailFrom,
                    'email_to' => $client->email,
                    'data' => $templateData,
                    'status' => NotificationEvents::CHANGE_RETURN_STATUS,
                    'user_id' => $reseller->id,
                    'client_id' => $client->id,
                    'to' => $data->differences['to'],
                ]));
                $result['notificationClientByEmail'] = $isSentEmailClient;
            }

            if (!empty($client->mobile)) {

                $isSentSms = (new SmsSender())->send(new SmsDTO([
                    'data' => $templateData,
                    'status' => NotificationEvents::CHANGE_RETURN_STATUS,
                    'user_id' => $reseller->id,
                    'client_id' => $client->id,
                    'to' => $data->differences['to'],
                ]));

                if ($isSentSms) {
                    $result['notificationClientBySms']['isSent'] = $isSentSms['isSent'];
                }
                if ($isSentSms['error']) {
                    $result['notificationClientBySms']['message'] = $isSentSms['error'];
                }
            }
        }

        return $result;
    }

    private function getDifferences(DifferencesDTO $data): string
    {
        $result = '';
        switch ($data->notification_type) {
            case NotificationType::TYPE_NEW:
                $result = __('NewPositionAdded', null, $data->user_id);
                break;
            case NotificationType::TYPE_CHANGE:
                $result = __('PositionStatusHasChanged', [
                    'FROM' => Status::getStatus($data->differences['from']),
                    'TO' => Status::getStatus($data->differences['to']),
                ], $data->user_id);
                break;
        }
        return $result;
    }
}
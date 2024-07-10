<?php

namespace NW\WebService\References\Operations\Notification\Services;

use NW\WebService\References\Operations\Notification\Entities\Contractor;
use NW\WebService\References\Operations\Notification\Helpers;
use NW\WebService\References\Operations\Notification\Events\NotificationEvent;

class NotificationService
{
    public function getResellerEmailFrom(int|string $resellerId): string
    {
        return Helpers\getResellerEmailFrom();
    }


    public function getEmailsByPermit(int|string $resellerId, string $permit): array
    {
        return Helpers\getEmailsByPermit($resellerId, $permit);
    }


    public function sendNotificationForEmployee(int $resellerId, array $templateData): bool
    {
        $resellerEmail = $this->getResellerEmailFrom($resellerId);
        $employeeEmails = $this->getEmailsByPermit($resellerId, 'tsGoodsReturn');

        if (empty($resellerEmail) || count($employeeEmails) === 0) {
            return false;
        }

        foreach ($employeeEmails as $employeeEmail) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $resellerEmail,
                    'emailTo' => $employeeEmail,
                    'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, NotificationEvent::CHANGE_RETURN_STATUS);
        }

        return true;
    }


    public function sendEmailToClient(Contractor $client, array $data, array $templateData): bool
    {
        $clientEmail = $client->getEmail();
        $resellerId = $data['resellerId'];

        $resellerEmail = $this->getResellerEmailFrom($resellerId);

        if (empty($resellerEmail)) {
            return false;
        }

        if (empty($clientEmail)) {
            return false;
        }

        MessagesClient::sendMessage([
            0 => [ // MessageTypes::EMAIL
                'emailFrom' => $resellerEmail,
                'emailTo' => $clientEmail,
                'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                'message' => __('complaintClientEmailBody', $templateData, $resellerId),
            ],
        ], $resellerId, $client->getId(), NotificationEvent::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);

        return true;
    }


    public function sendSmsToClient(Contractor $client, array $data, array $templateData): array
    {
        $clientMobile = $client->getMobile();
        $resellerId = $data['resellerId'];

        if (empty($clientMobile)) {
            return false;
        }

        $error = '';
        $result = NotificationManager::send(
            $resellerId,
            $client->getId(),
            NotificationEvent::CHANGE_RETURN_STATUS,
            (int)$data['differences']['to'],
            $templateData,
            $error
        );

        return ['isSent' => (bool) $result, 'message' => $error];
    }
}
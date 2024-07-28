<?php

namespace NW\WebService\References\Operations\Notification\Mailer;

use NW\WebService\References\Operations\Notification\Notification\Enum\NotificationEvent;
use NW\WebService\References\Operations\Notification\Validation\TemplateValidator;

/**
 * NotificationManager class
 */
class NotificationManager
{

    use MailableTrait;

    public function __construct()
    {
    }

    /**
     * Send employee notification.
     *
     * @return void
     * @throws \Exception
     */
    public function sendEmployeeNotification(array $data, array &$result): void
    {
        $resellerId = (int)($data['resellerId'] ?? 0);

        $templateData = (new NotificationTemplate($data))->getPreparedData();
        if (!(new TemplateValidator)->validate($templateData)) {
            return;
        }

        $emailFrom = $this->getResellerEmailFrom();
        // Получаем email сотрудников из настроек
        $emails = $this->getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvent::CHANGE_RETURN_STATUS->value);
                $result['notificationEmployeeByEmail'] = true;

            }
        }
    }

    /**
     * Send client notification.
     *
     * @return void
     * @throws \Exception
     */
    public function sendClientNotification(array $data, array &$result): void
    {
        $resellerId = (int)($data['resellerId'] ?? 0);

        $templateData = (new NotificationTemplate($data))->getPreparedData();
        if (!(new TemplateValidator)->validate($templateData)) {
            return;
        }

        if (!empty($emailFrom) && !empty($client->email)) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo' => $client->email,
                    'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvent::CHANGE_RETURN_STATUS->value, (int)$data['differences']['to']);
            $result['notificationClientByEmail'] = true;
        }

        if (!empty($client->mobile)) {
            $error = ''; // Get error
            $res = NotificationManager::send($resellerId, $client->id, NotificationEvent::CHANGE_RETURN_STATUS->value, (int)$data['differences']['to'], $templateData, $error);
            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($error)) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }
}

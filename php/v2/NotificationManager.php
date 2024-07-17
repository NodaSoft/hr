<?php

namespace NW\WebService\References\Operations\Notification;

use InvalidArgumentException;
use NW\WebService\References\Operations\Notification\Drivers\NotificationDriverFactory;
use NW\WebService\References\Operations\Notification\Drivers\NotificationDriverInterface;
use NW\WebService\References\Operations\Notification\Dto\NotificationData;
use NW\WebService\References\Operations\Notification\Notification\Enums\NotificationDriverEnum;

/**
 * Class NotificationManager
 *
 * Менеджер уведомлений, управляющий различными драйверами
 */
class NotificationManager
{
    private array $drivers = [];

    public function __construct(private readonly NotificationDriverFactory $driverFactory, private readonly TemplateService $templateService)
    {
        $this->drivers[NotificationDriverEnum::EMAIL->value] = $this->driverFactory->make(NotificationDriverEnum::EMAIL);
        $this->drivers[NotificationDriverEnum::SMS->value] = $this->driverFactory->make(NotificationDriverEnum::SMS);
    }

    public function driver(NotificationDriverEnum $driver): NotificationDriverInterface
    {
        if (!isset($this->drivers[$driver->value])) {
            throw new InvalidArgumentException("Driver [$driver->value] not supported.");
        }

        return $this->drivers[$driver->value];
    }

    /**
     * Отправляет уведомления сотрудникам
     */

    public function sendEmployeeNotifications(NotificationData $data, array &$result): void
    {
        $emailFrom = getResellerEmailFrom();
        $emails = getEmailsByPermit($data->resellerId, 'tsGoodsReturn');

        if (!empty($emailFrom) && count($emails) > 0) {

            $subject = $this->templateService->render('complaintClientEmailSubject', $data);
            $message = $this->templateService->render('complaintClientEmailBody', $data);

            foreach ($emails as $email) {

                $args = [
                    $emailFrom,
                    $email,
                    $subject,
                    $message,
                    $data->resellerId,
                    NotificationEvents::CHANGE_RETURN_STATUS
                ];

                $sent = $this->driver(NotificationDriverEnum::EMAIL)->send($args);
                if ($sent) {
                    $result['notificationEmployeeByEmail'] = true;
                }
            }

            $result['notificationEmployee'] = true;
        }
    }

    /**
     * Отправляет уведомления клиенту
     */
    public function sendClientNotifications(NotificationData $data, array &$result): void
    {
        $client = $this->getClientById($data->clientId);
        $emailFrom = $this->getResellerEmailFrom($data->resellerId);

        if (!empty($emailFrom) && !empty($client->email)) {
            $subject = $this->templateService->render('complaintClientEmailSubject', $data);
            $message = $this->templateService->render('complaintClientEmailBody', $data);

            $args = [
                $emailFrom,
                $client->email,
                $subject,
                $message,
                $data->resellerId,
                NotificationEvents::CHANGE_RETURN_STATUS
            ];

            $sent = $this->driver(NotificationDriverEnum::EMAIL)->send($args);
            if ($sent) {
                $result['notificationClientByEmail'] = true;
            }
        }

        if (!empty($client->mobile)) {
            $message = $this->templateService->render('complaintClientSms', $data);

            $args = [
                $client->mobile,
                $message
            ];

            $sent = $this->driver(NotificationDriverEnum::SMS)->send($args);
            $result['notificationClientBySms']['isSent'] = $sent;
            $result['notificationClientBySms']['message'] = $sent ? 'SMS sent successfully' : 'Failed to send SMS';
        }
    }

    private function getResellerEmailFrom(int $resellerId): ?string
    {
        // В реальном приложении здесь был бы код для получения email реселлера
        return 'reseller@example.com';
    }

    private function getEmailsByPermit(int $resellerId, string $permit): array
    {
        // В реальном приложении здесь был бы код для получения email адресов по разрешению
        return ['employee1@example.com', 'employee2@example.com'];
    }

    private function getClientById(int $clientId): object
    {
        // В реальном приложении здесь был бы код для получения данных клиента
        return (object)[
            'email' => 'client@example.com',
            'mobile' => '+1234567890',
        ];
    }
}
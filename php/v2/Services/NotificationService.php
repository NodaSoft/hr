<?php

namespace NW\WebService\References\Operations\Notification\Services;

use NW\WebService\References\Operations\Notification\Entities\Contractor;
use NW\WebService\References\Operations\Notification\Helpers\NotificationEvents;

use function NW\WebService\References\Operations\Notification\Helpers\getResellerEmailFrom;
use function NW\WebService\References\Operations\Notification\Helpers\getEmailsByPermit;

class NotificationService
{
  public function sendEmployeeNotifications(int $resellerId, array $templateData): bool
  {
    $emailFrom = $this->getResellerEmailFrom($resellerId);
    $emails    = $this->getEmailsByPermit($resellerId, 'tsGoodsReturn');

    if (empty($emailFrom) || empty($emails)) {
      return false;
    }

    foreach ($emails as $email) {
      $this->sendMessage([
        'emailFrom' => $emailFrom,
        'emailTo'   => $email,
        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
      ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
    }

    return true;
  }

  public function sendClientEmail(array $data, Contractor $client, array $templateData): bool
  {
    $emailFrom = $this->getResellerEmailFrom($data['resellerId']);

    if (empty($emailFrom) || empty($client->email)) {
      return false;
    }

    $this->sendMessage([
      'emailFrom' => $emailFrom,
      'emailTo'   => $client->email,
      'subject'   => __('complaintClientEmailSubject', $templateData, $data['resellerId']),
      'message'   => __('complaintClientEmailBody', $templateData, $data['resellerId']),
    ], $data['resellerId'], $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int) $data['differences']['to']);

    return true;
  }

  public function sendClientSms(array $data, Contractor $client, array $templateData): array
  {
    if (empty($client->mobile)) {
      return ['isSent' => false, 'message' => 'Client mobile number is empty'];
    }

    $error = '';
    $res   = $this->sendSmsNotification(
      $data['resellerId'],
      $client->id,
      NotificationEvents::CHANGE_RETURN_STATUS,
      (int) $data['differences']['to'],
      $templateData,
      $error
    );

    return [
      'isSent'  => $res,
      'message' => $error,
    ];
  }

  private function getResellerEmailFrom(int $resellerId): string
  {
    return getResellerEmailFrom();
  }

  private function getEmailsByPermit(int $resellerId, string $permit): array
  {
    return getEmailsByPermit($resellerId, $permit);
  }

  private function sendMessage(array $messageData, int $resellerId, int $clientId, string $event, ?int $differenceId = null): void
  {
    // Реализация должна быть заменена на реальную логику
  }

  private function sendSmsNotification(int $resellerId, int $clientId, string $event, int $statusId, array $templateData, string &$error): bool
  {
    // Реализация должна быть заменена на реальную логику
    return true;
  }
}
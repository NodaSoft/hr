<?php

namespace NW\WebService\References\Operations\Notification\Mailer;

use NW\WebService\References\Operations\Notification\Notification\Enum\NotificationType;

/**
 * NotificationTemplate class
 */
class NotificationTemplate
{

    protected array $data;

    public function __construct(array $data)
    {
        $this->data = $data;
    }

    /**
     * Returns prepared template data for notification.
     *
     * @return array
     */
    public function getPreparedData(): array
    {
        $cr = Employee::getById((int)($this->data['creatorId']));

        $et = Employee::getById((int)($this->data['expertId']));

        $client = Contractor::getById((int)($this->data['clientId']));

        $cFullName = $client->getFullName() ?: $client->name;

        $differences = $this->getDifferences(
            (int)($this->data['notificationType']),
            (int)($this->data['resellerId']),
            $this->data['differences']
        );

        return [
            'COMPLAINT_ID' => (int)($this->data['complaintId']),
            'COMPLAINT_NUMBER' => (string)($this->data['complaintNumber']),
            'CREATOR_ID' => (int)($this->data['creatorId']),
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => (int)($this->data['expertId']),
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => (int)($this->data['clientId']),
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => (int)($this->data['consumptionId']),
            'CONSUMPTION_NUMBER' => (string)($this->data['consumptionNumber']),
            'AGREEMENT_NUMBER' => (string)($this->data['agreementNumber']),
            'DATE' => (string)($this->data['date']),
            'DIFFERENCES' => $differences,
        ];
    }

    /**
     * Returns differences message.
     *
     * @param int $notificationType
     * @param int $recipientId
     * @param array $data
     * @return string
     */
    protected function getDifferences(int $notificationType, int $recipientId, array $data): string
    {
        $differences = '';

        if ($notificationType === NotificationType::NEW->value) {
            $differences = __('NewPositionAdded', null, $recipientId);
        } elseif ($notificationType === NotificationType::CHANGE->value
            && !empty($data)) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['from']),
                'TO' => Status::getName((int)$data['to']),
            ], $recipientId);
        }

        return $differences;
    }
}

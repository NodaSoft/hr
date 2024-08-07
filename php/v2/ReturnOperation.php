<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;
use NW\WebService\References\Operations\Notification\DTOs\EventsEnum;
use NW\WebService\References\Operations\Notification\DTOs\IndexDTO;
use NW\WebService\References\Operations\Notification\DTOs\NotificationEvents;
use NW\WebService\References\Operations\Notification\DTOs\SmsDTO;
use NW\WebService\References\Operations\Notification\Exceptions\InternalServerException;
use NW\WebService\References\Operations\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\Forms\DataForm;
use NW\WebService\References\Operations\Notification\Forms\IndexForm;
use NW\WebService\References\Operations\Notification\models\Contractor;
use NW\WebService\References\Operations\Notification\models\Employee;
use NW\WebService\References\Operations\Notification\models\ReferencesOperation;
use NW\WebService\References\Operations\Notification\models\Status;
use NW\WebService\References\Operations\Notification\Serializers\IndexSerializer;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;
    private function getDifferences(int $notificationType, array $differences, int $resellerId): string
    {
        if ($notificationType === self::TYPE_NEW) {
            return __('NewPositionAdded', null, $resellerId);
        }
        if ($notificationType === self::TYPE_CHANGE && !empty($differences)) {
            return __('PositionStatusHasChanged', [
                'FROM' => Status::findById((int)$differences['from'])?->name,
                'TO'   => Status::findById((int)$differences['to'])?->name,
            ], $resellerId);
        }
        return '';
    }

    private function prepareTemplateData(DataForm $form): array
    {
        $differences = $this->getDifferences($form->notificationType, $form->differences, $form->resellerId);

        return [
            'COMPLAINT_ID'       => $form->complaintId,
            'COMPLAINT_NUMBER'   => $form->complaintNumber,
            'CREATOR_ID'         => $form->creatorId,
            'CREATOR_NAME'       => $form->creator->getFullName(),
            'EXPERT_ID'          => $form->expertId,
            'EXPERT_NAME'        => $form->expert->getFullName(),
            'CLIENT_ID'          => $form->clientId,
            'CLIENT_NAME'        => $form->client->getFullName(),
            'CONSUMPTION_ID'     => $form->consumptionId,
            'CONSUMPTION_NUMBER' => $form->consumptionNumber,
            'AGREEMENT_NUMBER'   => $form->complaintNumber,
            'DATE'               => $form->date,
            'DIFFERENCES'        => $differences,
        ];
    }

    private function sendNotifications(DataForm $form, array $templateData): IndexDTO
    {
        $DTO = new IndexDTO(
            notificationEmployeeByEmail: false,
            notificationClientByEmail: false,
            notificationClientBySms: new SmsDTO(
                isSent: false,
                message: ''
            )
        );

        $emailFrom = getResellerEmailFrom($form->resellerId);
        $emails = getEmailsByPermit($form->resellerId, EventsEnum::TS_GOODS_RETURN);

        if ($emailFrom && !empty($emails)) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $form->resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $form->resellerId),
                    ],
                ], $form->resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $DTO->notificationClientByEmail = true;
            }
        }

        if ($form->notificationType === self::TYPE_CHANGE && !empty($form->differences['to'])) {
            if ($emailFrom && $form->client->email) {
                MessagesClient::sendMessage([
                    0 => [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $form->client->email,
                        'subject'   => __('complaintClientEmailSubject', $templateData, $form->resellerId),
                        'message'   => __('complaintClientEmailBody', $templateData, $form->resellerId),
                    ],
                ], $form->resellerId, $form->client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$form->differences['to']);
                $DTO->notificationEmployeeByEmail = true;
            }

            if ($form->client->mobile) {
                $error = '';
                $res = NotificationManager::send($form->resellerId, $form->client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$form->differences['to'], $templateData, $error);
                if ($res) {
                    $DTO->notificationClientBySms->isSent = true;
                }
                if ($error) {
                    $DTO->notificationClientBySms->message = $error;
                }
            }
        }

        return $DTO;
    }

    /**
     * @throws ValidationException
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');
        $form = new DataForm();
        $form->load($data);
        if (!$form->validate()) {
            throw new ValidationException($form->getErrors());
        }

        $templateData = $this->prepareTemplateData($form);

        $result = $this->sendNotifications($form, $templateData);

        return (new IndexSerializer($result))->jsonSerialize();
    }
}

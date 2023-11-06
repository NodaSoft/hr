<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class TsReturnOperation extends ReferencesOperation
{
    private $translator;
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    public function __construct(Translator $translator)
    {
        $this->translator = $translator;
    }

    /**
     * Выполняет операцию уведомления и возвращает результат.
     *
     * @return array Результат операции уведомления
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');
        $resellerId = (int)$data['resellerId'];
        $notificationType = (int)$data['notificationType'];

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $this->validateData($resellerId, $notificationType, $data);

        $differences = $this->getDifferencesText($notificationType, $data, $resellerId);

        $templateData = $this->buildTemplateData($data, $differences, $resellerId);

        $this->sendEmployeeNotifications($resellerId, $templateData, $result);
        $this->sendClientNotifications($notificationType, $data, $templateData, $result);

        return $result;
    }

    /**
     * @param $resellerId
     * @param $notificationType
     * @param $data
     * @return void
     * @throws Exception
     */
    private function validateData($resellerId, $notificationType, $data): void
    {
        if (empty($resellerId)) {
            throw new Exception('Empty resellerId', 400);
        }

        if (empty($notificationType)) {
            throw new Exception('Empty notificationType', 400);
        }

        $this->validateContractor($data, $resellerId);
        $this->validateEmployee($data);
    }

    /**
     * @param $data
     * @param $resellerId
     * @return void
     * @throws Exception
     */
    private function validateContractor($data, $resellerId): void
    {
        $client = Contractor::getById((int)$data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new Exception('Client not found!', 400);
        }
    }

    /**
     * @param $data
     * @return void
     * @throws Exception
     */
    private function validateEmployee($data): void
    {
        $this->validateEntity(Employee::getById((int)$data['creatorId']), 'Creator');
        $this->validateEntity(Employee::getById((int)$data['expertId']), 'Expert');
    }

    /**
     * @param $entity
     * @param $entityName
     * @return void
     * @throws Exception
     */
    private function validateEntity($entity, $entityName): void
    {
        if ($entity === null) {
            throw new Exception("$entityName not found!", 400);
        }
    }

    /**
     * @param $notificationType
     * @param $data
     * @param $resellerId
     * @return string
     */
    private function getDifferencesText($notificationType, $data, $resellerId): string
    {
        if ($notificationType === self::TYPE_NEW) {
            return $this->translator->translate('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            return $this->translator->translate('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO' => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }
        return '';
    }

    /**
     * @param $data
     * @param $differences
     * @param $resellerId
     * @return array
     */
    private function buildTemplateData($data, $differences, $resellerId): array
    {
        $client = Contractor::getById((int)$data['clientId']);
        $emailFrom = getResellerEmailFrom((int)$data['resellerId']);
        $cFullName = $client->getFullName();

        return [
            'COMPLAINT_ID' => (int)$data['complaintId'],
            'COMPLAINT_NUMBER' => (string)$data['complaintNumber'],
            'CREATOR_ID' => (int)$data['creatorId'],
            'CREATOR_NAME' => $this->getEmployeeFullName((int)$data['creatorId']),
            'EXPERT_ID' => (int)$data['expertId'],
            'EXPERT_NAME' => $this->getEmployeeFullName((int)$data['expertId']),
            'CLIENT_ID' => (int)$data['clientId'],
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$data['agreementNumber'],
            'DATE' => (string)$data['date'],
            'DIFFERENCES' => $differences,
        ];
    }

    /**
     * @param $employeeId
     * @return string
     */
    private function getEmployeeFullName($employeeId): string
    {
        $employee = Employee::getById($employeeId);
        return $employee ? $employee->getFullName() : '';
    }

    /**
     * @param $resellerId
     * @param $templateData
     * @param $result
     * @return void
     */
    private function sendEmployeeNotifications($resellerId, $templateData, &$result): void
    {
        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');

        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => $this->translator->translate('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => $this->translator->translate('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;
            }
        }
    }

    /**
     * @param $notificationType
     * @param $data
     * @param $templateData
     * @param $result
     * @return void
     */
    private function sendClientNotifications($notificationType, $data, $templateData, &$result): void
    {
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $resellerId = (int)$data['resellerId'];
            $client = Contractor::getById((int)$data['clientId']);
            $emailFrom = getResellerEmailFrom($resellerId);

            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $client->email,
                        'subject' => $this->translator->translate('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => $this->translator->translate('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            $this->sendClientSmsNotification($resellerId, $client, $data, $templateData, $result);
        }
    }

    /**
     * @param $resellerId
     * @param $client
     * @param $data
     * @param $templateData
     * @param $result
     * @return void
     */
    private function sendClientSmsNotification($resellerId, $client, $data, $templateData, &$result): void
    {
        if (!empty($client->mobile)) {
            $error = '';
            $res = NotificationManager::send(
                $resellerId,
                $client->id,
                NotificationEvents::CHANGE_RETURN_STATUS,
                (int)$data['differences']['to'],
                $templateData,
                $error
            );

            if ($res) {
                $result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($error)) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }
}

<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = $this->validateAndSanitizeInput((array)$this->getRequest('data'));
        $resellerId = $data['resellerId'];
        $notificationType = $data['notificationType'];
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        if (!$resellerId) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        if (!$notificationType) {
            throw new \Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById($data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        $clientFullName = $client->getFullName() ?: $client->name;

        $creator = Employee::getById($data['creatorId']);
        if ($creator === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $expert = Employee::getById($data['expertId']);
        if ($expert === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName($data['differences']['from']),
                'TO' => Status::getName($data['differences']['to']),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID' => $data['complaintId'],
            'COMPLAINT_NUMBER' => $data['complaintNumber'],
            'CREATOR_ID' => $data['creatorId'],
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $data['expertId'],
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_ID' => $data['clientId'],
            'CLIENT_NAME' => $clientFullName,
            'CONSUMPTION_ID' => $data['consumptionId'],
            'CONSUMPTION_NUMBER' => $data['consumptionNumber'],
            'AGREEMENT_NUMBER' => $data['agreementNumber'],
            'DATE' => $data['date'],
            'DIFFERENCES' => $differences,
        ];

        // Проверка наличия всех данных шаблона
        foreach ($templateData as $key => $value) {
            if (empty($value)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && !empty($emails)) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Уведомление клиента при смене статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo' => $client->email,
                        'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = '';
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $data['differences']['to'], $templateData, $error);
                $result['notificationClientBySms']['isSent'] = $res;
                $result['notificationClientBySms']['message'] = $error;
            }
        }

        return $result;
    }

    private function validateAndSanitizeInput(array $data): array
    {
        // Пример простой валидации и санитизации
        return [
            'resellerId' => isset($data['resellerId']) ? (int)filter_var($data['resellerId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'notificationType' => isset($data['notificationType']) ? (int)filter_var($data['notificationType'], FILTER_SANITIZE_NUMBER_INT) : null,
            'clientId' => isset($data['clientId']) ? (int)filter_var($data['clientId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'creatorId' => isset($data['creatorId']) ? (int)filter_var($data['creatorId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'expertId' => isset($data['expertId']) ? (int)filter_var($data['expertId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'complaintId' => isset($data['complaintId']) ? (int)filter_var($data['complaintId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'complaintNumber' => isset($data['complaintNumber']) ? htmlspecialchars($data['complaintNumber'], ENT_QUOTES, 'UTF-8') : '',
            'consumptionId' => isset($data['consumptionId']) ? (int)filter_var($data['consumptionId'], FILTER_SANITIZE_NUMBER_INT) : null,
            'consumptionNumber' => isset($data['consumptionNumber']) ? htmlspecialchars($data['consumptionNumber'], ENT_QUOTES, 'UTF-8') : '',
            'agreementNumber' => isset($data['agreementNumber']) ? htmlspecialchars($data['agreementNumber'], ENT_QUOTES, 'UTF-8') : '',
            'date' => isset($data['date']) ? htmlspecialchars($data['date'], ENT_QUOTES, 'UTF-8') : '',
            'differences' => isset($data['differences']) ? filter_var_array($data['differences'], FILTER_SANITIZE_FULL_SPECIAL_CHARS) : [],
        ];
    }
}

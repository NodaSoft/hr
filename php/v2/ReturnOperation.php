<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @return array
     * @throws Exception
     */
    public function doOperation(): array
    {
        $data = $this->getRequest('data');
        $resellerId = (int)($data['resellerId'] ?? 0);
        $clientId = (int)($data['clientId'] ?? 0);
        $creatorId = (int)($data['creatorId'] ?? 0);
        $expertId = (int)($data['expertId'] ?? 0);
        $notificationType = (int)($data['notificationType'] ?? 0);

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        if (empty($resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        if (empty($notificationType)) {
            throw new Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById((int)$resellerId);
        if ($reseller === null) {
            throw new Exception('Seller not found!', 400);
        }

        $client = Contractor::getById($clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->id !== $resellerId) {
            throw new Exception('Client not found!', 400);
        }

        $clientFullName = $client->getFullName();
        if (empty($clientFullName)) {
            $clientFullName = $client->name;
        }

        $creator = Employee::getById($creatorId);
        if ($creator === null) {
            throw new Exception('Creator not found!', 400);
        }

        $expert = Employee::getById($expertId);
        if ($expert === null) {
            throw new Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)($data['complaintId'] ?? 0),
            'COMPLAINT_NUMBER'   => (string)($data['complaintNumber'] ?? ''),
            'CREATOR_ID'         => $creatorId,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expertId,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $clientId,
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => (int)($data['consumptionId'] ?? 0),
            'CONSUMPTION_NUMBER' => (string)($data['consumptionNumber'] ?? ''),
            'AGREEMENT_NUMBER'   => (string)($data['agreementNumber'] ?? ''),
            'DATE'               => (string)($data['date'] ?? ''),
            'DIFFERENCES'        => $differences,
        ];

        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');

        if (!empty($emailFrom) && !empty($emails) && is_array($emails)) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);

                $result['notificationEmployeeByEmail'] = true;
            }
        }

        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    [
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $client->email,
                        'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);

                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = '';

                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                
                if ($res) {
                    $result['notificationClientBySms']['isSent'] = true;
                }

                if (!empty($error)) {
                    $result['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $result;
    }
}

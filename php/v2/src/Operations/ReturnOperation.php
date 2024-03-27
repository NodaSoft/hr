<?php

namespace App\Operations;

use App\Enum\ContractorType;
use App\Enum\Notification;
use App\Exceptions\EntityNotFoundException;
use App\Http\Request\BaseRequest;
use App\Models\Contractor;
use App\Models\Employee;
use App\Models\Seller;

class ReturnOperation
{
    /**
     * @throws EntityNotFoundException
     */
    public function newPosition(BaseRequest $request): array
    {
        $resellerId = $request->getDataField(BaseRequest::RESELLER_ID);
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];


        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new EntityNotFoundException('Reseller', $resellerId);
        }

        $clientId = $request->getDataField(BaseRequest::CLIENT_ID);
        $client = Contractor::getById($clientId);
        if (
            !$client instanceof Contractor ||
            $client->type !== ContractorType::CUSTOMER->value ||
            $client->seller->id !== $resellerId
        ) {
            throw new EntityNotFoundException('Client', $clientId);
        }

        $employeeId = $request->getDataField(BaseRequest::CREATOR_ID);
        $creator = Employee::getById($employeeId);
        if ($creator === null) {
            throw new EntityNotFoundException('Creator', $employeeId);
        }

        $expertId = $request->getDataField(BaseRequest::EXPERT_ID);
        $expert = Employee::getById($expertId);
        if ($expert === null) {
            throw new EntityNotFoundException('Expert', $employeeId);
        }

        $notificationType = $request->getDataField(BaseRequest::NOTIFICATION_TYPE);
        if ($notificationType === Notification::NEW->value) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === Notification::CHANGED && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO' => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID' => (int)$data['complaintId'],
            'COMPLAINT_NUMBER' => (string)$data['complaintNumber'],
            'CREATOR_ID' => (int)$data['creatorId'],
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => (int)$data['expertId'],
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => (int)$data['clientId'],
            'CLIENT_NAME' => $cFullName,
            'CONSUMPTION_ID' => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$data['agreementNumber'],
            'DATE' => (string)$data['date'],
            'DIFFERENCES' => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $client->email,
                        'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                        'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
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


    protected function getDifferencesArray()
    {

    }
}

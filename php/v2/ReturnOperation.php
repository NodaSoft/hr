<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    const TYPE_NEW = 1;
    const TYPE_CHANGE = 2;
    const EVENT_GOODS_RETURN = 'tsGoodsReturn';

    /**
     * @return array
     * @throws \Exception
     */
    public function doOperation(): array
    {
        try {
            $data = $this->getRequest('data');
            $requestData = new RequestDataDTO($data);

            $this->validateRequest($requestData);

            $reseller = $this->getResellerById($requestData->resellerId);
            $client = $this->getContractorById($requestData->clientId);

            $creator = $this->getEmployeeById($requestData->creatorId);
            $expert = $this->getEmployeeById($requestData->expertId);

            $differences = $this->getDifferences($requestData, $reseller);
            $templateData = $this->getTemplateData($requestData, $client, $creator, $expert, $differences);

            $this->validateTemplateData($templateData);

            $emailFrom = $this->getResellerEmailFrom();
            $emails = $this->getEmailsByPermit($reseller, $requestData->notificationType);


            $result = $this->initializeResultArray();

            $result = $this->sendEmployeeNotifications($result, $emailFrom, $emails, $templateData, $reseller);

            if ($this->isDataChanged($requestData)) {
                $result = $this->sendClientNotifications($result, $emailFrom, $client, $reseller, $templateData, $differences);
            }

            return $result;
        } catch (\Exception $e) {
            throw new \Exception($e->getMessage(), self::HTTP_BAD_REQUEST);
        }
    }

    /**
     * @param RequestDataDTO $requestData
     * @return void
     * @throws \Exception
     */
    private function validateRequest(RequestDataDTO $requestData): void
    {
        $errors = [];

        if (empty($requestData->resellerId)) {
            $errors[] = 'Empty resellerId';
        }

        if (empty($requestData->notificationType)) {
            $errors[] = 'Empty notificationType';
        }

        if (!empty($errors)) {
            throw new \Exception(implode(', ', $errors), 400);
        }
    }

    /**
     * @param int $resellerId
     * @return Seller
     * @throws \Exception
     */
    private function getResellerById(int $resellerId): Seller
    {
        if (!$reseller = Seller::getById($resellerId)) {
            throw new \Exception('Reseller not found!', 400);
        }

        return $reseller;
    }

    /**
     * @param int $clientId
     * @return Contractor
     * @throws \Exception
     */
    private function getContractorById(int $clientId): Contractor
    {
        if (!$client = Contractor::getById($clientId)) {
            throw new \Exception('Contractor not found!');
        }

        return $client;
    }

    /**
     * @param $employeeId
     * @return Employee
     * @throws \Exception
     */
    private function getEmployeeById($employeeId): Employee
    {
        if (!$employee = Employee::getById($employeeId)) {
            throw new \Exception('Employee not found!');
        }

        return $employee;
    }

    /**
     * @param RequestDataDTO $data
     * @param $reseller
     * @return array
     */
    private function getDifferences(RequestDataDTO $data, $reseller): array
    {
        $differences = [];

        switch ($data->notificationType) {
            case self::TYPE_NEW:
                $differences['NewPositionAdded'] = [
                    null,
                    $reseller->id
                ];
                break;
            case self::TYPE_CHANGE:
                $differences['PositionStatusHasChanged'] = [
                    'FROM' => Status::getName($data->differences['from']),
                    'TO' => Status::getName($data->differences['to']),
                    $reseller->id
                ];
                break;
        }

        return $differences;
    }

    /**
     * @param RequestDataDTO $data
     * @param $client
     * @param $creator
     * @param $expert
     * @param $differences
     * @return array
     */
    private function getTemplateData(RequestDataDTO $data, $client, $creator, $expert, $differences): array
    {
        return [
            'COMPLAINT_ID' => $data->complaintId,
            'COMPLAINT_NUMBER' => $data->complaintNumber,
            'CREATOR_ID' => $creator->id,
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $expert->id,
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_FULL_NAME' => $client->getFullName(),
            'CLIENT_ID' => $client->id,
            'CLIENT_NAME' => $client->name,
            'CONSUMPTION_ID' => $data->consumptionId,
            'CONSUMPTION_NUMBER' => $data->consumptionNumber,
            'AGREEMENT_NUMBER' => $data->agreementNumber,
            'DATE' => $data->date,
            'DIFFERENCES' => $differences,
        ];
    }

    /**
     * @param array $templateData
     * @return void
     * @throws \Exception
     */
    private function validateTemplateData(array $templateData): void
    {
        $emptyKeys = [];

        foreach ($templateData as $key => $value) {
            if (empty($value)) {
                $emptyKeys[] = $key;
            }
        }

        if (!empty($emptyKeys)) {
            $errorMessage = 'Empty template data: ' . implode(', ', array_keys($emptyKeys));
            throw new \Exception($errorMessage, 400);
        }
    }

    /**
     * @return string
     */
    private function getResellerEmailFrom(): string
    {
        return getResellerEmailFrom();
    }

    /**
     * @param Seller $reseller
     * @param string $event
     * @return string[]
     */
    private function getEmailsByPermit(Seller $reseller, string $event = self::EVENT_GOODS_RETURN): array
    {
        return getEmailsByPermit($reseller->id, $event);
    }

    /**
     * @return array
     */
    private function initializeResultArray(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false
            ]
        ];
    }

    /**
     * @param array $result
     * @param string $emailFrom
     * @param array $emails
     * @param array $templateData
     * @param Seller $reseller
     * @return array
     */
    private function sendEmployeeNotifications(array $result, string $emailFrom, array $emails, array $templateData, Seller $reseller): array
    {
        if (empty($emailFrom) && empty($emails)) {
            return $result;
        }

        foreach ($emails as $email) {
            $messageBody[MessageTypes::EMAIL] = [
                'emailFrom' => $emailFrom,
                'emailTo' => $email,
                'subject' => MessagesEmployee::getSubject('complaintEmployeeEmailSubject', $templateData, $reseller->id),
                'message' => MessagesEmployee::getMessage('complaintEmployeeEmailBody', $templateData, $reseller->id),
            ];

            MessagesEmployee::sendMessage($messageBody, $reseller->id, NotificationEvents::CHANGE_RETURN_STATUS);
        }

        $result['notificationEmployeeByEmail'] = true;

        return $result;
    }

    /**
     * @param array $result
     * @param string $emailFrom
     * @param Contractor $client
     * @param Seller $reseller
     * @param array $templateData
     * @param array $differences
     * @return array
     */
    private function sendClientNotifications(array $result, string $emailFrom, Contractor $client, Seller $reseller, array $templateData, array $differences): array
    {
        $differencesTo = $differences['PositionStatusHasChanged']['TO'];
        if (empty($differencesTo)) {
            return $result;
        }

        if (!empty($emailFrom) && !empty($client->email)) {
            $messageBody[MessageTypes::EMAIL] = [
                'emailFrom' => $emailFrom,
                'emailTo' => $client->email,
                'subject' => MessagesClient::getSubject('complaintClientEmailSubject', $templateData, $reseller->id),
                'message' => MessagesClient::getMessage('complaintClientEmailBody', $templateData, $reseller->id),
            ];

            MessagesClient::sendMessage($messageBody, $reseller->id, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo);

            $result['notificationClientByEmail'] = true;
        }

        if (!empty($client->mobile)) {
            $result['notificationClientBySms']['isSent'] = NotificationManager::send(int $reseller->id, int $client->id, string NotificationEvents::CHANGE_RETURN_STATUS, array $differencesTo, array $templateData);
        }

        return $result;
    }

    /**
     * @param RequestDataDTO $requestData
     * @return bool
     */
    private function isDataChanged(RequestDataDTO $requestData): bool
    {
        return $requestData->notificationType === self::TYPE_CHANGE;
    }
}

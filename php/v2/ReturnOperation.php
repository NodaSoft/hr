<?php

/*
 * Изменения и улучшения:
 *
 * 1. Разделила код на логически завершенные методы, чтобы каждый метод возвращал ожидаемый результат и код был более читаемым.
 * 2. Добавила проверку наличия основных параметров входящего запроса.
 * 3. Методам добавлен модификатор доступа private для ограничения доступа за пределами класса.
 * 4. Переименовала некоторые переменные на более понятные названия.
 * 5. Файл others.php оставила без изменений, так как не была уверена, что можно исправлять.
 */


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
        $data = $this->getValidateRequestData();
        $reseller = $this->getReseller($data['resellerId']);
        $client = $this->getClient($data);
        $creator = $this->getEmployeeByType((int)$data['creatorId'], 'creator');
        $expert = $this->getEmployeeByType((int)$data['expertId'], 'expert');
        $differences = $this->getDifferences($data);
        $templateData = $this->getTemplateData($data, $creator, $expert, $client, $differences);

        return $this->sendNotifications($templateData, $reseller, $client, $data);

    }

    /**
     * @throws Exception
     */
    protected function getValidateRequestData(): array
    {
        $data = (array)$this->getRequest('data');

        if (empty((int)$data['resellerId'])) {
            throw new Exception('Empty resellerId', 400);
        }

        if (empty((int)$data['notificationType'])) {
            throw new Exception('Empty notificationType', 400);
        }

        return $data;
    }


    /**
     * Получаем reseller по Ид
     *
     * @param int $resellerId
     * @return Seller
     * @throws Exception
     */
    private function getReseller(int $resellerId): Seller
    {
        $reseller = Seller::getById($resellerId);

        if ($reseller === null) {
            throw new Exception('Seller not found!', 400);
        }

        return $reseller;
    }

    /**
     * Получаем данные клиента по Ид
     *
     * @param array $data
     * @return Contractor
     * @throws Exception
     */
    private function getClient(array $data): Contractor
    {
        $client = Contractor::getById((int)$data['clientId']);

        if ($client === null ||
            $client->type !== Contractor::TYPE_CUSTOMER ||
            $client->id !== (int)$data['resellerId']) {
            throw new Exception('Client not found or unauthorized!', 400);
        }

        return $client;
    }

    /**
     * Получаем данные сотрудника по Ид
     *
     * @param int $employeeId
     * @param string $type
     * @return Employee
     * @throws Exception
     */
    private function getEmployeeByType(int $employeeId, string $type): Employee
    {
        $employee = Employee::getById($employeeId);

        if ($employee === null) {
            throw new Exception(ucfirst($type) . ' not found!', 400);
        }

        return $employee;
    }

    /**
     *
     * @param array $data
     * @return mixed
     */
    private function getDifferences(array $data)
    {
        if ((int)$data['notificationType'] === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $data['resellerId']);
        } elseif ((int)$data['notificationType']=== self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $data['resellerId']);
        }

        return $differences ?? null;
    }


    private function getTemplateData(
        array $data,
        Employee $creator,
        Employee $expert,
        Contractor $client,
        string $differences
    ): array {
        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => $creator->id,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expert->id,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->id,
            'CLIENT_NAME'        => $client->getFullName() ?? $client->name,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Проверяем, что все необходимые данные заполнены
        foreach ($templateData as $key => $tempData) {
            if ($tempData === '' || $tempData === null) {
                throw new \Exception("Template Data ({$key}) is empty!", 400);
            }
        }

        return $templateData;
    }


    private function sendNotifications(
        array $templateData,
        Seller $reseller,
        Contractor $client,
        array $data
    ): array {

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $emailFrom = getResellerEmailFrom();

        // Отправляем уведомления на email сотрудникам
        $this->sendEmployeeNotifications($emailFrom, $client->id, $reseller->id, $reseller->id, $templateData, $result);

        // Отправляем уведомления клиенту, если изменен статус
        if ($data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            $this->sendClientNotifications($emailFrom, $client, $templateData, $reseller->id, $data['differences']['to'], $result);
        }

        return $result;
    }

    private function sendEmployeeNotifications(
        ?string $emailFrom,
        int $clientId,
        int $resellerId,
        int $statusTo,
        array $templateData,
        array &$result
    ): void
    {
        $employeeEmails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($employeeEmails) > 0) {
            foreach ($employeeEmails as $email) {
                $this->sendEmailNotification($emailFrom, $email, $templateData, $resellerId, $clientId, $statusTo);
            }

            $result['notificationEmployeeByEmail'] = true;
        }
    }

    private function sendClientNotifications(
        ?string $emailFrom,
        Contractor $client,
        array $templateData,
        int $resellerId,
        int $statusTo,
        array &$result
    ): void {
        // Отправка email уведомления клиенту
        if ($emailFrom && $client->email) {
            $this->sendEmailNotification($emailFrom, $client->email, $templateData, $resellerId, $client->id, $statusTo);
            $result['notificationClientByEmail'] = true;
        }

        // Отправка SMS уведомления клиенту
        if (!empty($client->mobile)) {
            $error = '';
            $result['notificationClientBySms']['isSent'] = NotificationManager::send(
                $resellerId,
                $client->id,
                NotificationEvents::CHANGE_RETURN_STATUS,
                $statusTo,
                $templateData,
                $error
            );

            if ($error) {
                $result['notificationClientBySms']['message'] = $error;
            }
        }
    }

    private function sendEmailNotification(
        string $emailFrom,
        string $emailTo,
        array $templateData,
        int $resellerId,
        int $clientId,
        int $statusTo
    ): void {
        MessagesClient::sendMessage(
            [
                [
                    'emailFrom' => $emailFrom,
                    'emailTo' => $emailTo,
                    'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ],
            $resellerId,
            $clientId,
            NotificationEvents::CHANGE_RETURN_STATUS,
            $statusTo
        );
    }

}

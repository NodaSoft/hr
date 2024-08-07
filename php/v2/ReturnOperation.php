<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws NotFoundEntityException
     * @throws \Exception
     */
    public function doOperation(): array
    {
        /**
         * Т.к. рефакторинг происходит в отрыве от остальной части
         * программы, менять тип входящих и исходящих данных нельзя.
         * В дальнейшем следует перейти от использования массивов
         * к использованию объектов DTO. Это обеспечит безопасную
         * передачу данных от метода к методу.
         *
         * @todo Входящие и выходящие данные сделать DTO
         */
        $requestData = (array)$this->getRequest('data');
        $requestDataDTO = new ReturnOperationDTO(
            resellerId: (int)$requestData['resellerId'],
            clientId: (int)$requestData['clientId'],
            creatorId: (int)$requestData['creatorId'],
            expertId: (int)$requestData['expertId'],
            differencesFrom: (int)$requestData['differences']['from'],
            differencesTo: (int)$requestData['differences']['to'],
            notificationType: (int)$requestData['notificationType'],
            complaintId: (int)$requestData['complaintId'],
            consumptionId: (int)$requestData['consumptionId'],
            complaintNumber: (string)$requestData['complaintNumber'],
            consumptionNumber: (string)$requestData['consumptionNumber'],
            agreementNumber: (string)$requestData['agreementNumber'],
            date: (string)$requestData['date']
        );
        $this->validateRequestData($requestDataDTO);

        /**
         * Выглядит нелогичным, что notificationClientBySms - это массив,
         * а не булевое значение, как notificationEmployeeByEmail и notificationClientByEmail.
         *
         * Сообщение (message) содержит текст ошибки, имеет смысл переименовать в errorMessage и
         * вынести из notificationClientBySms, т.к. имеет отношение не только к СМС.
         *
         * @todo Пересмотреть структуру исходящих данных
         */
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        $reseller = Seller::getById($requestDataDTO->getResellerId());
        if ($reseller === null) {
            throw new NotFoundEntityException('Seller not found!', 400);
        }

        $client = Contractor::getById($requestDataDTO->getClientId());
        if (
            $client === null
            || $client->type !== Contractor::TYPE_CUSTOMER
            || $client->Seller->id !== $requestDataDTO->getResellerId()
        ) {
            throw new NotFoundEntityException('сlient not found!', 400);
        }

        $clientFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $clientFullName = $client->name;
        }

        $creator = Employee::getById($requestDataDTO->getCreatorId());
        if ($creator === null) {
            throw new NotFoundEntityException('Creator not found!', 400);
        }

        $expert = Employee::getById($requestDataDTO->getExpertId());
        if ($expert === null) {
            throw new NotFoundEntityException('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $requestDataDTO->getResellerId());
        } elseif ($notificationType === self::TYPE_CHANGE && $requestDataDTO->getDifferencesFrom() !== $requestDataDTO->getDifferencesTo()) {
            $differences = __('PositionStatusHasChanged', [
                    'FROM' => Status::getName($requestDataDTO->getDifferencesFrom()),
                    'TO'   => Status::getName($requestDataDTO->getDifferencesTo()),
                ], $requestDataDTO->getResellerId());
        }

        $templateData = [
            'COMPLAINT_ID'       => $requestDataDTO->getComplaintId(),
            'COMPLAINT_NUMBER'   => $requestDataDTO->getComplaintNumber(),
            'CREATOR_ID'         => $requestDataDTO->getCreatorId(),
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $requestDataDTO->getExpertId(),
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $requestDataDTO->getClientId(),
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => $requestDataDTO->getConsumptionId(),
            'CONSUMPTION_NUMBER' => $requestDataDTO->consumptionNumber,
            'AGREEMENT_NUMBER'   => $requestDataDTO->getAgreementNumber(),
            'DATE'               => $requestDataDTO->getDate(),
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($requestDataDTO->getResellerId());
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($requestDataDTO->getResellerId(), 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $email,
                           'subject'   => __('complaintEmployeeEmailSubject', $templateData, $requestDataDTO->getResellerId()),
                           'message'   => __('complaintEmployeeEmailBody', $templateData, $requestDataDTO->getResellerId()),
                    ],
                ], $requestDataDTO->getResellerId(), NotificationEvents::CHANGE_RETURN_STATUS);
                $result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($differencesTo)) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                           'subject'   => __('complaintClientEmailSubject', $templateData, $requestDataDTO->getResellerId()),
                           'message'   => __('complaintClientEmailBody', $templateData, $requestDataDTO->getResellerId()),
                    ],
                ], $requestDataDTO->getResellerId(), $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo);
                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($requestDataDTO->getResellerId(), $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo, $templateData, $error);
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

    /**
     * Валидация входящего запроса
     *
     * @param ReturnOperationDTO $operationDTO
     *
     * @throws ValidateRequestDataException
     */
    private function validateRequestData(ReturnOperationDTO $operationDTO): void
    {
        if ($operationDTO->getResellerId() === 0) {
            throw new ValidateRequestDataException('Empty resellerId', 400);
        }

        if ($operationDTO->getClientId() === 0) {
            throw new ValidateRequestDataException('Empty clientId', 400);
        }

        if ($operationDTO->getCreatorId() === 0) {
            throw new ValidateRequestDataException('Empty creatorId', 400);
        }

        if ($operationDTO->getExpertId() === 0) {
            throw new ValidateRequestDataException('Empty expertId', 400);
        }

        if ($operationDTO->getNotificationType() !== self::TYPE_NEW && $operationDTO->getNotificationType() !== self::TYPE_CHANGE) {
            throw new ValidateRequestDataException('Incorrect or empty notificationType', 400);
        }
    }
}

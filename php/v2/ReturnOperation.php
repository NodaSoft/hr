<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;
use NW\WebService\References\Operations\Notification\components\MessagesClient;
use NW\WebService\References\Operations\Notification\components\NotificationManager;
use NW\WebService\References\Operations\Notification\exceptions\EntityNotFoundException;
use NW\WebService\References\Operations\Notification\exceptions\ValidateRequestDataException;
use NW\WebService\References\Operations\Notification\models\Contractor;

class ReturnOperation extends ReferencesOperation
{
    public function __construct(
        private readonly NotificationManager $notificationManager,
        private readonly MessagesClient $messagesClient,
    )
    {
    }

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $requestData = (array)$this->getRequest('data');

        // It may be a builder, but I don't want to implement it.
        $requestDTO = new RequestDTO(
            resellerId: (int)$requestData['resellerId'],
            clientId: (int)$requestData['clientId'],
            creatorId: (int)$requestData['creatorId'],
            expertId: (int)$requestData['expertId'],
            differences: new Differences((int)$requestData['differences']['from'], (int)$requestData['differences']['to']),
            notificationType: NotificationType::tryFrom((int)$requestData['notificationType']),
            complaintId: (int)$requestData['complaintId'],
            consumptionId: (int)$requestData['consumptionId'],
            complaintNumber: (string)$requestData['complaintNumber'],
            consumptionNumber: (string)$requestData['consumptionNumber'],
            agreementNumber: (string)$requestData['agreementNumber'],
            date: (string)$requestData['date']
        );

        $this->validateRequestData($requestDTO);

        $resellerId  = $requestDTO->resellerId;
        $responseDTO = new ResponseDTO(
            notificationEmployeeByEmail: false,
            notificationClientByEmail: false,
            notificationClientBySms: new NotificationClientBySmsDTO(
                isSent: false,
                message: ''
            )
        );

        if (empty($requestDTO->resellerId)) {
            $responseDTO->notificationClientBySms->message = 'Empty resellerId';

            return $responseDTO->toArray();
        }

        $reseller = Seller::findById($resellerId);
        if ($reseller === null) {
            throw new EntityNotFoundException('Seller not found!', 400);
        }

        $client = Contractor::findById($requestDTO->clientId);
        if ($client === null || $client->type !== ContractorType::TYPE_CLIENT || $client->seller->id !== $resellerId) {
            throw new EntityNotFoundException('Client not found!', 400);
        }

        // Pointless code
        $clientFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $clientFullName = $client->name;
        }

        $templateData = $this->fetchTemplateData($requestDTO, $clientFullName, $resellerId);
        $emailFrom    = getResellerEmailFrom();

        $responseDTO->notificationClientBySms->isSent = $this->sendNotificationEmailToEmployee(
            $emailFrom,
            $templateData,
            $requestDTO,
        );

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($requestDTO->notificationType === NotificationType::TYPE_CHANGE && !empty($requestDTO->differences->to)) {
            $responseDTO->notificationClientByEmail = $this->sendNotificationEmailToClient(
                $client,
                $emailFrom,
                $templateData,
                $requestDTO
            );

            if (!empty($client->mobile)) {
                // It was undefined var $error. I suppose it must be provided as ref.
                // By the way, I would like to make response as an object with two statements: Success and Fail
                // So, I've deleted $error from parameter
                $response = $this->notificationManager->send(
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $requestDTO->differences->to,
                    $templateData
                );

                if ($response->error !== null) {
                    $responseDTO->notificationClientBySms->message = $response->error->message;
                } else {
                    $responseDTO->notificationClientBySms->isSent = true;
                }
            }
        }

        return $responseDTO->toArray();
    }

    /**
     * @throws ValidateRequestDataException
     */
    private function validateRequestData(RequestDTO $operationDTO): void
    {
        if ($operationDTO->clientId === 0) {
            throw new ValidateRequestDataException('Empty clientId', 400);
        }

        if (empty($operationDTO->notificationType)) {
            throw new ValidateRequestDataException('Empty notificationType', 400);
        }

        if ($operationDTO->creatorId === 0) {
            throw new ValidateRequestDataException('Empty creatorId', 400);
        }

        if ($operationDTO->expertId === 0) {
            throw new ValidateRequestDataException('Empty expertId', 400);
        }

        if ($operationDTO->notificationType !== NotificationType::TYPE_NEW && $operationDTO->notificationType !== NotificationType::TYPE_CHANGE) {
            throw new ValidateRequestDataException('Incorrect or empty notificationType', 400);
        }
    }

    private function getDifferences(RequestDTO $requestDTO, int $resellerId): string
    {
        $differences    = '';
        $hasDifferences = !empty($requestDTO->differences->from) && !empty($requestDTO->differences->to);

        if ($requestDTO->notificationType === NotificationType::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($requestDTO->notificationType === NotificationType::TYPE_CHANGE && $hasDifferences) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::from($requestDTO->differences->from)->getTextName(),
                'TO'   => Status::from($requestDTO->differences->to)->getTextName(),
            ], $resellerId);
        }

        return $differences;
    }

    /**
     * @throws EntityNotFoundException
     * @throws Exception
     */
    private function fetchTemplateData(RequestDTO $requestDTO, string $clientFullName, int $resellerId): array
    {
        $creator = Employee::findById($requestDTO->creatorId);
        if ($creator === null) {
            throw new EntityNotFoundException('Creator not found!', 400);
        }

        $expert = Employee::findById($requestDTO->expertId);
        if ($expert === null) {
            throw new EntityNotFoundException('Expert not found!', 400);
        }

        $templateData = [
            'COMPLAINT_ID'       => $requestDTO->complaintId,
            'COMPLAINT_NUMBER'   => $requestDTO->complaintNumber,
            'CREATOR_ID'         => $requestDTO->creatorId,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $requestDTO->expertId,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $requestDTO->clientId,
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => $requestDTO->consumptionId,
            'CONSUMPTION_NUMBER' => $requestDTO->complaintNumber,
            'AGREEMENT_NUMBER'   => $requestDTO->agreementNumber,
            'DATE'               => $requestDTO->date,
            'DIFFERENCES'        => $this->getDifferences($requestDTO, $resellerId),
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ($key) is empty!", 500);
            }
        }

        return $templateData;
    }

    private function sendNotificationEmailToEmployee(string $emailFrom, array $templateData, RequestDTO $requestDTO): bool
    {
        $isSent = false;
        // Получаем email сотрудников из настроек
        $emails     = getEmailsByPermit($requestDTO->resellerId, 'tsGoodsReturn'); // It may be another case in Notification Event, but I don't know any context
        $resellerId = $requestDTO->resellerId;

        foreach ($emails as $email) {
            $this->messagesClient->sendMessage([
                MessageTypes::EMAIL->value => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId), // IDK what it is
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId), // IDK what it is
                ],
            ], $resellerId, $requestDTO->clientId, NotificationEvents::CHANGE_RETURN_STATUS);

            $isSent = true;
        }

        return $isSent;
    }

    private function sendNotificationEmailToClient(Contractor $client, string $emailFrom, array $templateData, RequestDTO $requestDTO): bool
    {
        $isSent     = false;
        $resellerId = $requestDTO->resellerId;

        if (!empty($emailFrom) && !empty($client->email)) {
            $this->messagesClient->sendMessage([
                MessageTypes::EMAIL->value => [
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->email,
                    'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId), // IDK what it is
                    'message'   => __('complaintClientEmailBody', $templateData, $resellerId), // IDK what it is
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $requestDTO->differences->to);

            $isSent = true;
        }

        return $isSent;
    }
}

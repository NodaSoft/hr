<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Entities\Contractor;
use NW\WebService\References\Operations\Notification\Entities\Seller;
use NW\WebService\References\Operations\Notification\Entities\Employee;
use NW\WebService\References\Operations\Notification\Entities\Status;
use NW\WebService\References\Operations\Notification\Exceptions\CreatorNotExistException;
use NW\WebService\References\Operations\Notification\Exceptions\ExpertNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\SellerNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\ClientNotFoundException;
use NW\WebService\References\Operations\Notification\Services\NotificationService;


class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    public function __construct(
        private NotificationService $notificationService
    )
    {
    }

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array)$this->getRequest('data');
        $resellerId = $data['resellerId'];
        $notificationType = (int)$data['notificationType'];

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        if (empty($resellerId)) {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $result;
        }

        if (empty($notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }

        $reseller = $this->getReseller($resellerId);
        $client = $this->getClient((int)$data['clientId']);
        $creator = $this->getCreator((int)$data['creatorId']);
        $expert = $this->getExpert((int)$data['expertId']);

        $templateData = $this->getTemplateData($data, $client, $creator, $expert);

        $result['notificationEmployeeByEmail'] = $this->notificationService
            ->sendNotificationForEmployee($resellerId, $templateData);

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($this->isStatusChanged($data)) {

            $result['notificationClientByEmail'] = $this->notificationService
                ->sendEmailToClient($client, $data, $templateData);

            $result['notificationClientBySms'] = $this->notificationService
                ->sendSmsToClient($client, $data, $templateData);
        }

        return $result;
    }


    private function getReseller(int $resellerId): Contractor
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new SellerNotFoundException();
        }

        return $reseller;
    }


    private function getClient(int $clientId): Contractor
    {
        $client = Contractor::getById($clientId);
        if (
            $client === null
            || $client->getType() !== Contractor::TYPE_CUSTOMER
            || $client->getSeller()->getId() !== $clientId
        ) {
            throw new ClientNotFoundException();
        }

        return $client;
    }


    private function getCreator(int $creatorId): Employee
    {
        $creator = Employee::getById((int)$creatorId);
        if ($creator === null) {
            throw new CreatorNotExistException();
        }

        return $creator;
    }


    private function getExpert(int $expertId): Employee
    {
        $expert = Employee::getById((int)$expertId);
        if ($expert === null) {
            throw new ExpertNotFoundException();
        }

        return $expert;
    }


    private function getDifferences(array $data): string
    {
        $differences = '';
        $notificationType = $data['notificationType'];
        $resellerId = $data['resellerId'];

        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO' => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        return $differences;

    }


    private function getTemplateData(
        array $data,
        Contractor $client,
        Employee $creator,
        Employee $expert
    ): array
    {
        $templateData = [
            'COMPLAINT_ID' => (int)$data['complaintId'],
            'COMPLAINT_NUMBER' => (string)$data['complaintNumber'],
            'CREATOR_ID' => $creator->getId(),
            'CREATOR_NAME' => $creator->getFullName() ?? $creator->getName(),
            'EXPERT_ID' => $expert->getId(),
            'EXPERT_NAME' => $expert->getFullName() ?? $expert->getName(),
            'CLIENT_ID' => $client->getId(),
            'CLIENT_NAME' => $client->getFullName() ?? $client->getName(),
            'CONSUMPTION_ID' => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$data['agreementNumber'],
            'DATE' => (string)$data['date'],
            'DIFFERENCES' => $this->getDifferences($data),
        ];

        foreach ($templateData as $key => $templateValue) {
            if (empty($templateValue)) {
                throw new \TemplateDataIsEmptyException($key);
            }
        }

        return $templateData;
    }


    private function isStatusChanged(array $data): bool
    {
        return $data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences']['to']);
    }

}

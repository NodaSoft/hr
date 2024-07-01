<?php

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Exceptions\ValidationException;
use NW\WebService\References\Operations\Notification\Services\NotificationService;
use NW\WebService\References\Operations\Notification\Validators\DataValidator;
use NW\WebService\References\Operations\Notification\Entities\{Contractor, Seller, Employee};
use NW\WebService\References\Operations\Notification\Helpers\ReferencesOperation;

class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    private array $result = [
        'notificationEmployeeByEmail' => false,
        'notificationClientByEmail'   => false,
        'notificationClientBySms'     => [
            'isSent'  => false,
            'message' => '',
        ],
    ];

    private DataValidator $validator;
    private NotificationService $notificationService;

    public function __construct(DataValidator $validator, NotificationService $notificationService)
    {
        $this->validator           = $validator;
        $this->notificationService = $notificationService;
    }

    /**
     * @throws ValidationException
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');
        $this->validator->validateData($data);

        $reseller = $this->getReseller($data['resellerId']);
        $client   = $this->getClient($data['clientId'], $data['resellerId']);
        $creator  = $this->getEmployee($data['creatorId'], 'Creator');
        $expert   = $this->getEmployee($data['expertId'], 'Expert');

        $templateData = $this->prepareTemplateData($data, $client, $creator, $expert);

        $this->sendEmployeeNotifications($data['resellerId'], $templateData);

        if ($this->isStatusChangeNotification($data)) {
            $this->sendClientNotifications($data, $client, $templateData);
        }

        return $this->result;
    }

    private function getReseller(int $resellerId): Seller
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new ValidationException('Seller not found!', 400);
        }
        return $reseller;
    }

    private function getClient(int $clientId, int $resellerId): Contractor
    {
        $client = Contractor::getById($clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new ValidationException('Client not found!', 400);
        }
        return $client;
    }

    private function getEmployee(int $employeeId, string $role): Employee
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            throw new ValidationException("$role not found!", 400);
        }
        return $employee;
    }

    private function prepareTemplateData(array $data, Contractor $client, Employee $creator, Employee $expert): array
    {
        $templateData = [
            'COMPLAINT_ID'       => (int) $data['complaintId'],
            'COMPLAINT_NUMBER'   => (string) $data['complaintNumber'],
            'CREATOR_ID'         => (int) $data['creatorId'],
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => (int) $data['expertId'],
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => (int) $data['clientId'],
            'CLIENT_NAME'        => $client->getFullName() ?: $client->name,
            'CONSUMPTION_ID'     => (int) $data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string) $data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string) $data['agreementNumber'],
            'DATE'               => (string) $data['date'],
            'DIFFERENCES'        => $this->getDifferences($data),
        ];

        $this->validator->validateTemplateData($templateData);

        return $templateData;
    }

    private function getDifferences(array $data): string
    {
        if ($data['notificationType'] === self::TYPE_NEW) {
            return __('NewPositionAdded', null, $data['resellerId']);
        } elseif ($data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences'])) {
            return __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) $data['differences']['from']),
                'TO'   => Status::getName((int) $data['differences']['to']),
            ], $data['resellerId']);
        }
        return '';
    }

    private function isStatusChangeNotification(array $data): bool
    {
        return $data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences']['to']);
    }

    private function sendEmployeeNotifications(int $resellerId, array $templateData): void
    {
        $this->result['notificationEmployeeByEmail'] = $this->notificationService->sendEmployeeNotifications($resellerId, $templateData);
    }

    private function sendClientNotifications(array $data, Contractor $client, array $templateData): void
    {
        $this->result['notificationClientByEmail'] = $this->notificationService->sendClientEmail($data, $client, $templateData);
        $this->result['notificationClientBySms']   = $this->notificationService->sendClientSms($data, $client, $templateData);
    }
}
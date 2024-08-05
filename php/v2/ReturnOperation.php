<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use Exception;
use RuntimeException;
use NW\WebService\Exceptions\BadRequestException;
use NW\WebService\Exceptions\NotFoundException;
use NW\WebService\References\Operations\Email\EmailMessage;

class ReturnOperation extends ReferencesOperation
{
    protected const TYPE_NEW = 1;
    protected const TYPE_CHANGE = 2;

    protected const NOTIFICATION_TYPES = [
        self::TYPE_NEW,
        self::TYPE_CHANGE
    ];

    private array $data = [];
    private ?Contractor $client = null;
    private ?Employee $creator = null;
    private ?Employee $expert = null;
    private ?array $templateData = null;

    /**
     * @throws BadRequestException
     * @throws NotFoundException
     * @throws Exception
     */
    public function doOperation(): array
    {
        $this->data = (array)$this->getRequest('data');

        $this->validateInputData();

        $reseller = $this->getReseller((int)$this->data['resellerId']);
        $this->client = $this->getClient((int)$this->data['clientId'], $reseller->id);
        $this->creator = $this->getEmployee($this->data['creatorId'], 'Creator');
        $this->expert = $this->getEmployee($this->data['expertId'], 'Expert');

        $this->prepareTemplateData();
        $this->validateTemplateData();

        return $this->sendNotifications($reseller->id);
    }

    /**
     * @throws BadRequestException
     */
    private function validateInputData(): void
    {
        $requiredFields = ['resellerId', 'notificationType', 'clientId', 'creatorId', 'expertId'];

        foreach ($requiredFields as $field) {
            if (empty($this->data[$field])) {
                throw new BadRequestException("Empty or invalid {$field}");
            }
        }

        if (!in_array($this->data['notificationType'], self::NOTIFICATION_TYPES, true)) {
            throw new BadRequestException('Invalid notificationType');
        }

        if (
            $this->data['notificationType'] === self::TYPE_CHANGE &&
            (!isset($this->data['differences']['from'], $this->data['differences']['to']))
        ) {
            throw new BadRequestException('Differences data is missing or invalid for TYPE_CHANGE');
        }
    }


    /**
     * @param int $resellerId
     * @return Seller
     * @throws NotFoundException
     */
    private function getReseller(int $resellerId): Seller
    {
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new NotFoundException('Reeller not found!');
        }

        return $reseller;
    }

    /**
     * @throws NotFoundException
     */
    private function getClient(int $clientId, int $resellerId): Contractor
    {
        $client = Contractor::getById($clientId);

        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || ($client->Seller && $client->Seller->id !== $resellerId)) {
            throw new NotFoundException('Client not found!');
        }

        return $client;
    }

    /**
     * @throws NotFoundException
     */
    private function getEmployee(int $employeeId, string $role): Employee
    {
        $employee = Employee::getById($employeeId);
        if ($employee === null) {
            throw new NotFoundException("{$role} not found!");
        }

        return $employee;
    }

    private function prepareTemplateData(): void
    {
        $differences = $this->getDifferencesMessage((int)$this->data['notificationType']);

        $this->templateData = [
            'COMPLAINT_ID' => (int)($this->data['complaintId'] ?? 0),
            'COMPLAINT_NUMBER' => (string)($this->data['complaintNumber'] ?? ''),
            'CREATOR_ID' => (int)$this->data['creatorId'],
            'CREATOR_NAME' => $this->creator->getFullName(),
            'EXPERT_ID' => (int)$this->data['expertId'],
            'EXPERT_NAME' => $this->expert->getFullName(),
            'CLIENT_ID' => (int)$this->data['clientId'],
            'CLIENT_NAME' => $this->client->getFullName() ?: $this->client->name,
            'CONSUMPTION_ID' => (int)($this->data['consumptionId'] ?? 0),
            'CONSUMPTION_NUMBER' => (string)($this->data['consumptionNumber'] ?? ''),
            'AGREEMENT_NUMBER' => (string)($this->data['agreementNumber'] ?? ''),
            'DATE' => (string)($this->data['date'] ?? ''),
            'DIFFERENCES' => $differences,
        ];
    }

    private function getDifferencesMessage(int $notificationType): string
    {
        switch ($notificationType) {
            case self::TYPE_NEW:
                return __('NewPositionAdded', null, (int)$this->data['resellerId']);

            case self::TYPE_CHANGE:
                if (!empty($this->data['differences'])) {
                    return __('PositionStatusHasChanged', [
                        'FROM' => Status::getName((int)$this->data['differences']['from']),
                        'TO' => Status::getName((int)$this->data['differences']['to']),
                    ], (int)$this->data['resellerId']);
                }
                break;
        }

        throw new RuntimeException('Invalid notification type or missing differences.');
    }

    /**
     * @throws Exception
     */
    private function validateTemplateData(): void
    {
        foreach ($this->templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \RuntimeException("Template Data ({$key}) is empty!", 500);
            }
        }
    }

    private function getResellerEmail(int $resellerId): ?string
    {
        return getResellerEmailFrom($resellerId);
    }

    private function getEmployeeEmails(int $resellerId): array
    {
        return getEmailsByPermit($resellerId, 'tsGoodsReturn');
    }

    private function sendNotifications(int $resellerId): array
    {
        $result = [
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        $emailFrom = $this->getResellerEmail($resellerId);

        $result['notificationEmployeeByEmail'] = $this->notifyEmployees($resellerId, $emailFrom);

        if ((int)$this->data['notificationType'] === self::TYPE_CHANGE && !empty($this->data['differences']['to'])) {
            $result = $this->notifyClient($resellerId, $emailFrom, $result);
        }

        return $result;
    }

    private function sendClientSms(int $resellerId, array &$result): void
    {
        $error = '';

        $sendResult = NotificationManager::send($resellerId, $this->client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$this->data['differences']['to'], $this->templateData, $error);
        if ($sendResult) {
            $result['notificationClientBySms']['isSent'] = true;
        }

        if (!empty($error)) {
            $result['notificationClientBySms']['message'] = $error;
        }
    }

    private function notifyEmployees(int $resellerId, ?string $emailFrom): bool
    {
        $emails = $this->getEmployeeEmails($resellerId);
        if (!empty($emailFrom) && !empty($emails)) {
            foreach ($emails as $email) {
                $emailMessage = new EmailMessage(
                    $emailFrom,
                    $email,
                    'complaintEmployeeEmailSubject',
                    'complaintEmployeeEmailBody',
                    $resellerId,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $this->templateData
                );
                $emailMessage->send();
            }
            return true;
        }
        return false;
    }

    private function notifyClient(int $resellerId, ?string $emailFrom, array $result): array
    {
        if (!empty($emailFrom) && !empty($this->client->email)) {
            $emailMessage = new EmailMessage(
                $emailFrom,
                $this->client->email,
                'complaintClientEmailSubject',
                'complaintClientEmailBody',
                $resellerId,
                NotificationEvents::CHANGE_RETURN_STATUS,
                $this->templateData,
                $this->client->id,
                (int)$this->data['differences']['to']
            );
            $emailMessage->send();

            $result['notificationClientByEmail'] = true;
        }

        if (!empty($this->client->mobile)) {
            $this->sendClientSms($resellerId, $result);
        }

        return $result;
    }
}

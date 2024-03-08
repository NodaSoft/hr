<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    private array $templateData;
    private ?int $resellerId = null;
    private ?int $notificationType = null;
    private ?array $data;
    private ?Contractor $client;
    private ?Employee $creator;
    private ?Employee $expert;

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $this->data = $this->getRequest('data');

        $this->initParams();

        return $this->notifyAll();
    }

    private function emptyResult(): array
    {
        return [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }

    /**
     * @throws Exception
     */
    private function getDifferences(): string
    {
        if ($this->notificationType === self::TYPE_NEW) {
            $templateName = 'NewPositionAdded';
            $templateParams = null;
        } elseif ($this->notificationType === self::TYPE_CHANGE && !empty($this->data['differences'])) {
            $templateName = 'PositionStatusHasChanged';
            $templateParams = [
                'FROM' => Status::getName($this->data['differences']['from']),
                'TO' => Status::getName($this->data['differences']['to'])
            ];
        }

        if (empty($templateName) || empty($templateParams)) {
            return '';
        }

        return __($templateName, $templateParams, $this->resellerId);
    }

    private function notificationClientByEmail(string $from, string $to): bool
    {
        try {
            $this->sendEmail($from, $to);
        } catch (\Throwable $throwable) {
            return false;
        }
        return true;
    }

    private function sendEmail(string $from, string $to, ?int $type = MessageTypes::EMAIL)
    {
        MessagesClient::sendMessage(
            [
                $type => [
                    'emailFrom' => $from,
                    'emailTo' => $to,
                    'subject' => __('complaintEmployeeEmailSubject', $this->templateData, $this->resellerId),
                    'message' => __('complaintEmployeeEmailBody', $this->templateData, $this->resellerId),
                ],
            ],
            $this->resellerId,
            NotificationEvents::CHANGE_RETURN_STATUS
        );
    }

    /**
     * @return bool
     * \throwable
     */
    private function sendSms() : bool
    {
        return NotificationManager::send(
            $this->resellerId,
            $this->client->id,
            NotificationEvents::CHANGE_RETURN_STATUS,
            $this->data['differences']['to'],
            $this->templateData,
        );
    }

    /**
     * @throws Exception
     */
    private function getTemplateData() : array
    {
        $prepareString = fn($s) => trim($s);

        $templateData = [
            'COMPLAINT_ID' => $this->toIntOrNull($this->data['complaintId']),
            'COMPLAINT_NUMBER' => $prepareString($this->data['complaintNumber']),
            'CREATOR_ID' => $this->creator->id,
            'CREATOR_NAME' => $this->creator->getFullName(),
            'EXPERT_ID' => $this->expert->id,
            'EXPERT_NAME' => $this->expert->getFullName(),
            'CLIENT_ID' => $this->client->id,
            'CLIENT_NAME' => $this->client->getFullName(),
            'CONSUMPTION_ID' => $this->toIntOrNull($this->data['consumptionId']),
            'CONSUMPTION_NUMBER' => $prepareString($this->data['consumptionNumber']),
            'AGREEMENT_NUMBER' => $prepareString($this->data['agreementNumber']),
            'DATE' => $prepareString($this->data['date']),
            'DIFFERENCES' => $this->getDifferences(),
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", self::HTTP_BAD_REQUEST);
            }
        }

        return $templateData;
    }

    /**
     * @throws Exception
     */
    private function initParams(): void
    {
        if (!$this->resellerId = $this->toIntOrNull($this->data['resellerId'])) {
            throw new Exception('Wrong resellerId', self::HTTP_BAD_REQUEST);
        }

        $this->notificationType = $this->data['notificationType'];

        if (!in_array($this->notificationType, [self::TYPE_NEW, self::TYPE_CHANGE])) {
            throw new Exception('Wrong notificationType', self::HTTP_BAD_REQUEST);
        }

        $this->client = Contractor::getById($this->data['clientId']);

        if ( $this->client->type !== Contractor::TYPE_CUSTOMER
            || $this->client->Seller->id !== $this->resellerId
        ) {
            throw new Exception('Client not found!', self::HTTP_BAD_REQUEST);
        }

        $this->creator = Employee::getById($this->data['creatorId']);
        $this->expert = Employee::getById($this->data['expertId']);
        $this->templateData = $this->getTemplateData();
    }

    private function notifyAll(): array
    {
        $result = $this->emptyResult();

        $emailFrom = getResellerEmailFrom($this->resellerId);

        foreach (getEmailsByPermit($this->resellerId, 'tsGoodsReturn') as  $email) {
            $result['notificationEmployeeByEmail'][$email] = $this->notificationClientByEmail($emailFrom, $email);
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($this->notificationType === self::TYPE_CHANGE && !empty($this->data['differences']['to'])) {

            if (!empty($emailFrom) && !empty($this->client->email)) {
                $result['notificationClientByEmail'] =
                    $this->notificationClientByEmail($emailFrom, $this->client->email);
            }

            if (!empty($this->client->mobile)) {
                try {
                    $result['notificationClientBySms']['isSent'] = $this->sendSms();
                } catch (\Throwable $throwable) {
                    $result['notificationClientBySms']['message'] = $throwable->getMessage();
                }
            }
        }

        return $result;
    }

    private function toIntOrNull(mixed $param): ?int
    {
        return is_numeric($param) ? (int)$param : null;
    }
}

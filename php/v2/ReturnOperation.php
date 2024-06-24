<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const int TYPE_NEW = 1;
    public const int TYPE_CHANGE = 2;
    private string $differences = '';
    private int $resellerId;
    private int $notificationType;
    private array $data;
    private Contractor $contractor;
    private Contractor $reseller;
    private Contractor $creator;
    private Contractor $expert;
    private int $clientId;
    private int $creatorId;
    private int $expertId;
    private array $result;
    private array $templateData;

    public function getDifferences(): string
    {
        return $this->differences;
    }

    public function setDifferences(): void
    {
        if ($this->notificationType === static::TYPE_NEW) {
            $this->differences = __('NewPositionAdded', null, $this->resellerId);
        } elseif ($this->notificationType === static::TYPE_CHANGE && !empty($this->data['differences'])) {
            $this->differences = __('PositionStatusHasChanged', [
                'FROM' => Status::from((int)$this->data['differences']['from'])->name,
                'TO' => Status::from((int)$this->data['differences']['to'])->name,
            ], $this->resellerId);
        }
    }

    /**
     * @return array
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $this->setValues();
        $this->validate();

        $this->result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];

        if (!$this->resellerId) {
            $this->result['notificationClientBySms']['message'] = 'Empty resellerId';
            return $this->result;
        }

        $this->reseller = Seller::getById($this->resellerId);

        $this->setContractor();

        $this->creator = Employee::getById($this->creatorId);

        $this->expert = Employee::getById($this->expertId);

        $this->setDifferences();

        $this->setTemplateData();

        // Получаем email сотрудников из настроек
        $emailFrom = getResellerEmailFrom($this->resellerId);
        $this->notify($emailFrom);
        if ($this->notificationType === static::TYPE_CHANGE && !empty($this->data['differences']['to'])) {
            $this->notifyStatusChange($emailFrom);
        }

        return $this->result;
    }

    /**
     * @return void
     */
    private function setValues(): void
    {
        $this->data = $this->getRequest('data');
        $this->resellerId = (int)$this->data['resellerId'];
        $this->notificationType = (int)$this->data['notificationType'];
        $this->clientId = (int)$this->data['clientId'];
        $this->creatorId = (int)$this->data['creatorId'];
        $this->expertId = (int)$this->data['expertId'];
    }

    /**
     * @return void
     * @throws \Exception
     */
    public function validate(): void
    {
        if (!$this->notificationType) {
            throw new \Exception('Empty notificationType', 400);
        }

        if (!$this->clientId) {
            throw new \Exception('Client ID is empty!', 400);
        }

        if (!$this->creatorId) {
            throw new \Exception('Creator ID is empty!', 400);
        }

        if (!$this->expertId) {
            throw new \Exception('Expert ID is empty!', 400);
        }
    }

    /**
     * Шлём клиентское уведомление, только если произошла смена статуса
     * @param string $emailFrom
     * @param $error
     * @return void
     */
    public function notifyStatusChange(string $emailFrom): void
    {
        if ($emailFrom && $this->contractor->email) {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo' => $this->contractor->email,
                    'subject' => __('complaintClientEmailSubject', $this->templateData, $this->resellerId),
                    'message' => __('complaintClientEmailBody', $this->templateData, $this->resellerId),
                ],
            ],
                $this->resellerId,
                $this->contractor->getId(),
                NotificationEvents::CHANGE_RETURN_STATUS->value,
                (int)$this->data['differences']['to']
            );
            $this->result['notificationClientByEmail'] = true;
        }

        if ($this->contractor->mobile) {
            $error = '';
            $res = NotificationManager::send(
                $this->resellerId, $this->contractor->getId(),
                NotificationEvents::CHANGE_RETURN_STATUS->value,
                (int)$this->data['differences']['to'],
                $this->templateData,
                $error
            );
            if ($res) {
                $this->result['notificationClientBySms']['isSent'] = true;
            }
            if ($error) {
                $this->result['notificationClientBySms']['message'] = $error;
            }
        }
    }

    /**
     * @param string $emailFrom
     * @return void
     */
    public function notify(string $emailFrom): void
    {
        $emails = getEmailsByPermit($this->resellerId, 'tsGoodsReturn');
        if ($emailFrom && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo' => $email,
                        'subject' => __('complaintEmployeeEmailSubject', $this->templateData, $this->resellerId),
                        'message' => __('complaintEmployeeEmailBody', $this->templateData, $this->resellerId),
                    ],
                ], $this->resellerId, NotificationEvents::CHANGE_RETURN_STATUS->value);
                $this->result['notificationEmployeeByEmail'] = true;
            }
        }
    }

    /**
     * @return void
     * @throws \Exception
     */
    public function setContractor(): void
    {
        $this->contractor = Contractor::getById($this->clientId);
        if (
            $this->contractor->getType() !== Contractor::TYPE_CUSTOMER ||
            $this->contractor->Seller->getId() !== $this->resellerId
        ) {
            throw new \Exception('Client not found!', 400);
        }
    }

    /**
     * @return void
     * @throws \Exception
     */
    public function setTemplateData(): void
    {
        $this->templateData = [
            'COMPLAINT_ID' => (int)$this->data['complaintId'],
            'COMPLAINT_NUMBER' => (string)$this->data['complaintNumber'],
            'CREATOR_ID' => $this->creatorId,
            'CREATOR_NAME' => $this->creator->getFullName(),
            'EXPERT_ID' => $this->expertId,
            'EXPERT_NAME' => $this->expert->getFullName(),
            'CLIENT_ID' => $this->clientId,
            'CLIENT_NAME' => $this->contractor->getFullName(),
            'CONSUMPTION_ID' => (int)$this->data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$this->data['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$this->data['agreementNumber'],
            'DATE' => (string)$this->data['date'],
            'DIFFERENCES' => $this->getDifferences(),
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($this->templateData as $key => $tempData) {
            if (!$tempData) {
                throw new \Exception("Template Data ($key) is empty!", 500);
            }
        }
    }
}

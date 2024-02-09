<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    public Seller $reseller;
    public int $notificationType = 0;
    public NotificationData $data;
    protected $templateData = [];

    /**
     * В данном случае реализация паттерна Command (Operation) целью которого является сделать
     * непрямым взаимодействие исполнителя (MessagesClient и NotificationManager) и вызывающего класса (контроллера)
     * Реализовано не полностью тк :
     * - исполнитель лезет в данные запроса Request (а значит привязан к слою контроллера), а также в нем есть HTTP коды,
     * которые могут иметь какое-то значение только в рамках http, лучше кидать эксепшены бизнес-логики, 
     * уже в слое контроллер их преобразовывать для ответа клииенту, текущее решение снижает гибкость использования операции 
     * то есть например осуществить данную операцию из командной строки уже невозможно (исправлено)
     * - происходит два разных действия с точки зрения бизнеса (Админское уведомление и Клиентское уведомление), что странно (не исправлено)
     * И есть некоторые вопросы, оставшиеся неисправленными
     * - надо бы валидацию данных (not found for example) унести в контроллер, здесь оставить только бизнес-проверки
     *  - что надо делать если одному админу письмо ушло, а на втором упало, надо ли слать клиенту
     * @throws \Exception
     */
    public function doOperation(): array
    {
        if ($this->notificationType == 0  || !in_array($this->notificationType, [self::TYPE_CHANGE, self::TYPE_NEW]) ) {
            throw new BusinessException('Empty notificationType');
        }
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
        if (empty($this->reseller)) {
            throw new BusinessException('Empty resellerId!');
        }

        // Эта операция должна пройти в контроллере при валидации данных
        // $reseller = Seller::getById((int)$resellerId);
        // if ($reseller === null) {
        //     throw new \Exception('Seller not found!', 400);
        // }

        $client = Contractor::getById($this->data->clientId);
        if (
            $client === null
            || $client->type !== Contractor::TYPE_CUSTOMER
            || ($client->Seller && $client->Seller->id !== $this->reseller->id)
            ) {
            throw new BusinessException('сlient not found!');
        }

        $cr = Employee::getById($this->data->creatorId);
        if ($cr === null) {
            throw new BusinessException('Creator not found!');
        }

        $expert = Employee::getById($this->data->expertId);
        if ($expert === null) {
            throw new BusinessException('Expert not found!');
        }

        $differences = '';
        if ($this->notificationType === self::TYPE_NEW) {
            $differences = ResellerView::make('NewPositionAdded', null, $this->reseller->id);
        } elseif ($this->notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = ResellerView::make('PositionStatusHasChanged', [
                    'FROM' => Status::getName($this->data->differencesFrom),
                    'TO'   => Status::getName($this->data->differencesTo),
                ], $this->reseller->id);
        }

        $this->templateData = [
            'COMPLAINT_ID'       => $this->data->complaintId,
            'COMPLAINT_NUMBER'   => $this->data->complaintNumber,
            'CREATOR_ID'         => $cr->id,
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => $expert->id,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->id,
            'CLIENT_NAME'        => $client->getFullName() ??  $client->name,
            'CONSUMPTION_ID'     => $this->data->consumptionId,
            'CONSUMPTION_NUMBER' => $this->data->consumptionNumber,
            'AGREEMENT_NUMBER'   => $this->data->agreementNumber,
            'DATE'               => $this->data->date,
            'DIFFERENCES'        => $differences,
        ];
        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($this->templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new BusinessErrorException("Template Data ({$key}) is empty!");
            }
        }

        $result['notificationEmployeeByEmail'] = $this->adminNotification();

        if ($this->notificationType === self::TYPE_CHANGE && !empty($this->data->differencesTo)) {
            $result['notificationClientByEmail'] = $this->clientEmailNotification($client);
            $error = '';
            $result['notificationClientBySms']['isSent'] = $this->clientMobileNotification($client, $error);
        }

        return $result;
    }

    protected function adminNotification() {
        $result = false;
        $emailFrom = getResellerEmailFrom($this->reseller->id);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($this->reseller->id, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                try {
                    MessagesClient::sendMessage([
                        MessagesClient::EMAIL => [
                               'emailFrom' => $emailFrom,
                               'emailTo'   => $email,
                               'subject'   => ResellerView::make('complaintEmployeeEmailSubject', $this->templateData, $this->reseller->id),
                               'message'   => ResellerView::make('complaintEmployeeEmailBody', $this->templateData, $this->reseller->id),
                        ],
                    ], $this->reseller->id, NotificationEvents::CHANGE_RETURN_STATUS);
                    $result = true;
                } catch(\Exception $e) {
                }
            }
        }
        return $result;
    }

    protected function clientEmailNotification($client) {
        $result = false;
        if (!empty($emailFrom) && !empty($client->email)) {
            MessagesClient::sendMessage([
                MessagesClient::EMAIL => [
                       'emailFrom' => $emailFrom,
                       'emailTo'   => $client->email,
                       'subject'   => ResellerView::make('complaintClientEmailSubject', $this->templateData, $this->reseller->id),
                       'message'   => ResellerView::make('complaintClientEmailBody', $this->templateData, $this->reseller->id),
                ],
            ], $this->reseller->id, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $this->data->differencesTo);
            $result = true;
        }
        return $result;
    }

    protected function clientMobileNotification($client, &$error)
    {
        $result = false;
        if (!empty($client->mobile)) {
            $res = NotificationManager::send($this->reseller->id, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $this->data->differencesTo, $this->templateData, $error);
            if ($res) {
                $result = true;
            }
        }
        return $result;
    }
}

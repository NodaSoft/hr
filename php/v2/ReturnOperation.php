<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

/*
 * Рассылка уведомлений об изменениях у продавца "чего-то"
 * Обобщил и упростил, возможно какие-то моменты сходу не заметил
 */
class ReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;
    public const TYPES = [
        self::TYPE_NEW,
        self::TYPE_CHANGE
    ];

    public const FIELDS_NOT_EMPTY = [
        'complaintId',
        'complaintNumber',
        'expertId',
        'creatorId',
        'clientId',
        'sellerId',
        'notificationType',
        'consumptionId',
        'consumptionNumber',
        'agreementNumber',
        'date',
        'differences',
    ];

    private array $data = [];
    private array $templateData;
    private int $notificationType;
    private string $differences;
    private Seller $seller;
    private Customer $client;
    private Employee $creator;
    private Employee $expert;

    /**
     * Базовая проверка на заполненность полей
     *
     * @throws Exception
     */
    protected function checkFields(): void
    {
        $this->data = $this->getRequest('data');
        foreach (self::FIELDS_NOT_EMPTY as $field) {
            if (empty($this->data[$field])) {
                throw new Exception(sprintf('Empty %s field', $field), 400);
            }
        }
    }

    /**
     * Обрабатываем и проверяем тип уведомления
     *
     * @throws Exception
     */
    protected function getNotificationType(): void
    {
        if (!in_array((int)$this->data['notificationType'], self::TYPES)) {
            throw new Exception('Unknown notificationType %s field', 400);
        }

        $this->notificationType = (int)$this->data['notificationType'];
    }

    /**
     * Универсальная функция получения и проверки юзеров
     *
     * @throws Exception
     */
    protected function getUser(string $field, string $class, $errorMsg): Contractor|Seller|Employee|Customer
    {
        if (class_exists($class)) {
            $user = $class::getById((int)$this->data[$field]);
        }

        if (empty($user)) {
            throw new Exception($errorMsg, 400);
        }

        return $user;
    }

    /**
     * Подготавливаем необходимых пользователей
     *
     * @return void
     * @throws Exception
     */
    protected function prepareUsers(): void
    {
        $this->seller = $this->getUser('sellerId', Seller::class, 'Seller not found!');
        $this->client = $this->getUser('clientId', Customer::class, 'Client not found!');
        // дополнительная проверка
        if ($this->client->Seller->getId() !== $this->seller->getId()) {
            throw new Exception('Client and seller not linked!', 400);
        }

        $this->creator = $this->getUser('creatorId', Employee::class, 'Creator not found!');
        $this->expert = $this->getUser('expertId', Employee::class, 'Expert not found!');
    }

    /**
     * Проверяем и формируем строку различий
     *
     * @throws Exception
     */
    protected function getDifferences(): void
    {
        if ($this->notificationType === self::TYPE_NEW) {
            $this->differences = __('NewPositionAdded', null, $this->seller->getId());

            return;
        }

        if (
            $this->notificationType === self::TYPE_CHANGE
            && !empty($this->data['differences']['from'])
            && !empty($this->data['differences']['to'])
        ) {
            $fromStatusName = !empty($data['differences']['from']) ?
                Status::getName((int)$data['differences']['from']) : null;
            $toStatusName = !empty($data['differences']['to']) ?
                Status::getName((int)$data['differences']['to']): null;
            if ($fromStatusName && $toStatusName) {
                $this->differences = __('PositionStatusHasChanged', [
                    'FROM' => $fromStatusName ,
                    'TO'   => $toStatusName,
                ], $this->seller->getId());

                return;
            }
        }

        throw new Exception('Unknown differences', 400);
    }

    /**
     * Проверяем и формируем данные для шаблонов писем
     *
     * @throws Exception
     */
    protected function getTemplateData(): void
    {
        /*
         * По-хорошему, наверно надо более корректно получать и проверять данные $this->data
         */
        $templateData = [
            'COMPLAINT_ID'       => (int)$this->data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$this->data['complaintNumber'],
            'CREATOR_ID'         => $this->creator->getId(),
            'CREATOR_NAME'       => $this->creator->getFullName(),
            'EXPERT_ID'          => $this->expert->getId(),
            'EXPERT_NAME'        => $this->expert->getFullName(),
            'CLIENT_ID'          => $this->client->getId(),
            'CLIENT_NAME'        => $this->client->getFullName(),
            'CONSUMPTION_ID'     => (int)$this->data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$this->data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$this->data['agreementNumber'],
            'DATE'               => (string)$this->data['date'],
            'DIFFERENCES'        => $this->differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $this->templateData = $templateData;
    }

    /**
     * Выполняем операцию рассылки уведомлений
     * Можно было бы сделать так, чтобы функция всегда выдавала нужную структуру ответа + ошибки, без генерирования исключений,
     * но допускаю, что этого не требуется
     *
     * @throws Exception
     */
    public function doOperation(): array
    {
        // проверяем, все ли поля переданы
        $this->checkFields();
        // проверяем тип уведомлений
        $this->getNotificationType();
        // определяем всех пользователй
        $this->prepareUsers();
        // получает строку различий
        $this->getDifferences();
        // подготавливаем данные для шаблона
        $this->getTemplateData();

        // формируем ответ
        $result = [
            'notificationEmployeeByEmail' => $this->sendEmailEmployee(), // шлем сотрудникам
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($this->notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            // почта - клиенту
            $result['notificationClientByEmail'] = $this->sendEmailClient();
            // sms - клиенту
            $result['notificationClientBySms'] = $this->sendSMSClient();
        }

        return $result;
    }

    /*
     * Отправка СМС клиенту
     */
    protected function sendSMSClient(): array
    {
        $result = [
            'isSent' => false,
            'message' => '',
        ];

        if (!empty($this->client->getMobile())) {
            $error = null;
            $res = NotificationManager::send(
                $this->seller->getId(),
                $this->client->getId(),
                NotificationEvents::CHANGE_RETURN_STATUS,
                (int)$this->data['differences']['to'],
                $this->templateData,
                $error
            );
            if ($res) {
                $result['isSent'] = true;
            }
            if (!empty($error)) {
                $result['message'] = $error;
            }
        }

        return $result;
    }

    /*
     * Отправка почты клиенту
     */
    protected function sendEmailClient(): bool
    {
        $emailFrom = getSellerEmailFrom();
        if ($this->client->getEmail()) {
            return $this->sendEmails(
                $emailFrom,
                [$this->client->getEmail()],
                (int)$this->data['differences']['to']
            );
        }

        return false;
    }

    /*
     * Отправка почты сотрудникам
     */
    protected function sendEmailEmployee(): bool
    {
        $emailFrom = getSellerEmailFrom();
        // Получаем email сотрудников из настроек.
        $emails = getEmailsByPermit($this->seller->getId(), 'tsGoodsReturn');
        if ($emails) {
            return $this->sendEmails(
                $emailFrom,
                $emails
            );
        }

        return false;
    }

    /**
     * Отправляем почту, общая функция
     *
     * @param string $emailFrom
     * @param array $toEmails
     * @param int|null $differences
     * @return bool
     */
    protected function sendEmails(
        string $emailFrom,
        array  $toEmails,
        int    $differences = null
    ): bool
    {
        if ($emailFrom && $toEmails) {
            $sellerId = $this->seller->getId();
            $templateData = $this->templateData;
            $messages = array_map(static function ($email) use ($emailFrom, $templateData, $sellerId) {
                return [
                    // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $email,
                    'subject'   => __('complaintEmployeeEmailSubject', $templateData, $sellerId),
                    'message'   => __('complaintEmployeeEmailBody', $templateData, $sellerId),
                ];
            }, $toEmails);

            return MessagesClient::sendMessage(
                $messages,
                $sellerId,
                NotificationEvents::CHANGE_RETURN_STATUS,
                $differences
            );
        }

        return false;
    }
}

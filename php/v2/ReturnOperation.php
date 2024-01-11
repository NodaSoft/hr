<?php
/*
 * это первичный рефакторинг
 * конечно многое зависит от дальнейшей реализации, тут много заглушек функций
 * файл others.php нужно разбивать на отдельные файлы и прорабатывать вглубь
 * что касается ReturnOperation - основную функцию конечно лучше переименовать, возможно есть смысл засунуть основную последовательность действий
 * в construct
 * Довольно много вопросов по валидации REQUEST, который используется для даты - лучше всю валидации вынести в отдельный интерфейс, а не располагать в том же классе
 * Константы TYPE_* можно вынести в enum, но неизвестно какую версию пхп используем. С синтаксисом по этому поводу тоже есть спорные моменты
 * Возникают вопросы почему в части случаев выбрасывается исключения, а для пустого $this->resellerId используется result['notificationClientBySms']['message'] = 'Empty resellerId';
 * возможно исключение должны быть везде или наоборот при исключении должно уходить уведомление всегда, изначальная постановка задачи неизвестна - поэтому много вопросов как это дальше рефакторить
 *
 * Теперь про назначение - очевидно это отправка уведомлений клиентам и сотрудниками, предусмотрено два канала - электропочта и смс
 *
 * */
namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{

    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    private $data   = [];
    private $result = [];

    private $notificationType;
    private $resellerId;
    private $client;
    private $reseller;
    private $creator;
    private $expert;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $this->setData();

        $this->resellerId = (int) $this->data['resellerId'];
        if (!$this->resellerId) {
            $this->result['notificationClientBySms']['message'] = 'Empty resellerId';

            return $this->result;
        }

        $this->notificationType = (int) $this->data['notificationType'];
        $this->setResult();

        $this->setParams();

        $templateData = $this->setTemplateData();

        $emailFrom = getResellerEmailFrom($this->resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($this->resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                $this->sendEmail($emailFrom, $email, $templateData, 'complaintEmployeeEmailSubject', 'complaintEmployeeEmailBody');
                $this->result['notificationEmployeeByEmail'] = true;

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($this->notificationType === self::TYPE_CHANGE && !empty($this->data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($this->client->email)) {
                $this->sendEmail($emailFrom, $this->client->email, $templateData, 'complaintClientEmailSubject', 'complaintClientEmailBody', (int) $this->data['differences']['to']);
                $this->result['notificationClientByEmail'] = true;
            }

            $this->sendSms($templateData);
        }

        return $this->result;
    }

    private function sendEmail($emailFrom, $emailTo, $templateData, $subject, $message, $differenceTo = 0)
    {
        MessagesClient::sendMessage([
            0 => [ // MessageTypes::EMAIL
                'emailFrom' => $emailFrom,
                'emailTo'   => $emailTo,
                'subject'   => __($subject, $templateData, $this->resellerId),
                'message'   => __($message, $templateData, $this->resellerId),
            ],
        ], $this->resellerId, $this->client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differenceTo);
    }

    private function sendSms(array $templateData): void
    {
        // will be better use try catch for error
        if (!empty($this->client->mobile)) {
            $res = NotificationManager::send($this->resellerId, $this->client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int) $this->data['differences']['to'], $templateData);
            if ($res['status']) {
                $this->result['notificationClientBySms']['isSent'] = true;
            }
            if (!empty($res['error'])) {
                $this->result['notificationClientBySms']['message'] = $res['error'];
            }
        }
    }

    private function setData()
    {
        $this->data = (array) $this->getRequest('data');
    }

    private function setResult()
    {
        $this->result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];
    }

    private function setParams(): void
    {
        if (!$this->notificationType) {
            throw new \Exception('Empty notificationType', 400);
        }

        $this->reseller = Seller::getById($this->resellerId);
        if ($this->reseller == 0) {
            throw new \Exception('Seller not found!', 400);
        }

        $this->client = Contractor::getById((int) $this->data['clientId']);
        if ($this->client == 0 || $this->client->type !== Contractor::TYPE_CUSTOMER || $this->client->Seller->id !== $this->resellerId) {
            throw new \Exception('сlient not found!', 400);
        }

        $this->creator = Employee::getById((int) $this->data['creatorId']);
        if ($this->creator == 0) {
            throw new \Exception('Creator not found!', 400);
        }

        $this->expert = Employee::getById((int) $this->data['expertId']);
        if ($this->expert == 0) {
            throw new \Exception('Expert not found!', 400);
        }
    }

    private function setTemplateData(): array
    {
        $differences = '';
        if ($this->notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $this->resellerId);
        } elseif ($this->notificationType === self::TYPE_CHANGE && !empty($this->data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) $this->data['differences']['from']),
                'TO'   => Status::getName((int) $this->data['differences']['to']),
            ], $this->resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int) $this->data['complaintId'],
            'COMPLAINT_NUMBER'   => (string) $this->data['complaintNumber'],
            'CREATOR_ID'         => (int) $this->data['creatorId'],
            'CREATOR_NAME'       => $this->creator->getFullName(),
            'EXPERT_ID'          => (int) $this->data['expertId'],
            'EXPERT_NAME'        => $this->expert->getFullName(),
            'CLIENT_ID'          => (int) $this->data['clientId'],
            'CLIENT_NAME'        => empty($this->client->getFullName()) ? $this->client->name : $this->client->getFullName(),
            'CONSUMPTION_ID'     => (int) $this->data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string) $this->data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string) $this->data['agreementNumber'],
            'DATE'               => (string) $this->data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        return $templateData;
    }

}

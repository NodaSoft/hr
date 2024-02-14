<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;
use Throwable;
use NW\WebService\References\Operations\Notification\{
    Models\Seller,
    Models\Contractor,
    Models\Employee,
    Models\Status,
    Services\MessagesClient,
    Services\NotificationManager,
    Enums\NotificationEvents,
};

/**
 * Краткое резюме.
 * В данной реализации используется обращение к Request, было бы правильнее сделать валидацию этих данных на уровне выши и передавать в эту функцию уже класс через DTO
 * Это позволит использовать модель что минимизирует возможность получение ошибок при написании кода. При использовании ключей массива этот риск сохраняется
 * Тк это уровень бизнес-логики и нет приямого взаимодействия с контекстом, то нет смысла возвращать коды ошибок (исправлено)
 * Функции отправки сообщений (почта и смс) требуется обернуть в try-catch (добавлено)
 * Валидация типа $notificationType сделана в начале чтобы код не выполнял запросов к хранилищу за получением Seller, Contractor и Employee, тк далее все равно не пройдет валидацию
 * В коде было несколько ошибок, думаю, что оставленных целенаправлено: тип возвращаемой функции, отсутствие объявления используемых моделей, не переданный параметр $client->id.
 * Наверняка чото-то еще, чего я уже не нашел)
 * Ну и к коду наверно будет не лишним добавить форматирование по стандарту PSR
 */
class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        // Неправильно из слоя бизнес-логики лезть в Request. На вышестоящем уровне было бы правильнее провести валидацию данные и через DTO вернуть класс
        $data = (array)$this->getRequest('data');
        $resellerId = (int)$data['resellerId'];
        $notificationType = (int)$data['notificationType'];
        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        if ($resellerId === 0) {
            throw new Exception('Empty resellerId');
        }

        // Вообще в данном случае достаточно проверки на notificationType, но, если допускается, что могут быть какие-то другие типы, вроде TYPE_DELETE и TYPE_DISABLED, то проверку на тип стоит убрать.
        // Правда тогда код пройдет дальше, но все равно не пройдет валидацию $templateData
        if ($notificationType == 0  || !in_array($notificationType, [self::TYPE_CHANGE, self::TYPE_NEW]) ) {
            throw new Exception('Empty notificationType');
        }

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new Exception('Seller not found!');
        }

        $client = Contractor::getById((int)$data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new Exception('сlient not found!');
        }

        // Данную логику было правильнее перенести в модель и далее использовать только $client->getFullName()
        $clientFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $clientFullName = $client->name;
        }

        $creator = Employee::getById((int)$data['creatorId']);
        if ($creator === null) {
            throw new Exception('Creator not found!');
        }

        $expert = Employee::getById((int)$data['expertId']);
        if ($expert === null) {
            throw new Exception('Expert not found!');
        }

        $differences = '';
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => (int)$data['creatorId'],
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $clientFullName,
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new Exception("Template Data ({$key}) is empty!");
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');

        $allEmailsSentResult = false;
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                // Вот тут нужно определиться с возвращаемым результатом: какой будет валидным: если хоть одна отправка отработала корректно, то отдаем true или все должны отработать без ошибок и только тогда - true
                // В первом случае в коде ниже вместо AND нужно поставить OR
                // должно получиться вот так: $allEmailsSentResult = $allEmailsSentResult || this->sendEmail(
                $allEmailsSentResult = $allEmailsSentResult && this->sendEmail(
                    $emailFrom,
                    $email,
                    __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS
                );

            }
        }
        $result['notificationEmployeeByEmail'] = $allEmailsSentResult;

        $result['notificationClientBySms']['isSent'] = false;
        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && isset($data['differences']['to'])) {
            $sentResult = false;
            if (!empty($emailFrom) && !empty($client->email)) {
                $sentResult = this->sendEmail(
                    $emailFrom,
                    $client->email,
                    __('complaintClientEmailSubject', $templateData, $resellerId),
                    __('complaintClientEmailBody', $templateData, $resellerId),
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$data['differences']['to']
                );
            }
            $result['notificationClientByEmail'] = $sentResult;

            if (!empty($client->mobile)) {
                $error = '';
                try {
                    $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                    if ($res) {
                        $result['notificationClientBySms']['isSent'] = $res;
                    }
                    if (!empty($error)) {
                        $result['notificationClientBySms']['message'] = $error;
                    }
                } catch (Throwable $e) {
                    $result['notificationClientBySms']['message'] = $e->getMessage();
                }
            }
        }

        return $result;
    }

    /**
     * @param string $from
     * @param string $to
     * @param string $subject
     * @param string $body
     * @param int $resellerId
     * @param int $clientId
     * @param int $notificationStatus
     * @param int|null $statusCodeTo
     * @return bool
     */
    private function sendEmail(
        string $from,
        string $to,
        string $subject,
        string $body,
        int $resellerId,
        int $clientId,
        int $notificationStatus,
        int $statusCodeTo = null
    )
    {
        // В исходном коде в вызове данной функции (в цикле) либо пропущен параметр $client->id либо в самой функции проводится анализ типа переменной передаваемой в 3й параметр, а остальные являются не обязательными...
        // Предполагаю, что логика описанная во втором подходе добавит излишнюю логику в код поэтому отталкиваюсь от того, что параметр просто пропущен
        // Вообще если знать реализацию MessagesClient::sendMessage, то этоу реализацию переписать не долго
        $sentResult = true;
        try {
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $from,
                    'emailTo' => $to,
                    'subject' => $subject,
                    'message' => $body,
                ],
            ], $resellerId, $clientId, $notificationStatus, $statusCodeTo);

        } catch(Throwable $e) {
            $sentResult = false;
        }

        return $sentResult;
    }
}

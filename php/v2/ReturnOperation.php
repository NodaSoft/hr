<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

// Структуры данных

// Базовый класс для DTO (я обычно использую spatie/laravel-data для запросов/валидации/переноса между сервисами/в качестве ресурса)
class Data implements \ArrayAccess
{
    public function offsetExists(mixed $offset): bool
    {
        return property_exists($this, $offset);
    }

    public function offsetGet(mixed $offset): mixed
    {
        if ($this->offsetExists($offset)) {
            return $this->$offset;
        } else {
            return null;
        }
    }

    public function offsetSet(mixed $offset, mixed $value): void
    {
        $this->$offset = $value;
    }

    public function offsetUnset(mixed $offset): void
    {
        unset($this->$offset);
    }
}

class ClientNotificationBySmsDto extends Data
{
    public $isSent = false;
    public $message = '';
}

class OperationResultDto extends Data
{
    public $notificationEmployeeByEmail = false;
    public $notificationClientByEmail = false;
    public $notificationClientBySms = new ClientNotificationBySmsDto;
}

/**
 * Тут таки следует использовать PSR-12 для именования свойств
 */
class TemplateDataDto extends Data
{
    public int $COMPLAINT_ID = 0;
    public string $COMPLAINT_NUMBER = '';
    public int $CREATOR_ID = 0;
    public string $CREATOR_NAME = '';
    public int $EXPERT_ID = 0;
    public string $EXPERT_NAME = '';
    public int $CLIENT_ID = 0;
    public string $CLIENT_NAME = '';
    public int $CONSUMPTION_ID = 0;
    public string $CONSUMPTION_NUMBER = '';
    public string $AGREEMENT_NUMBER = '';
    public string $DATE = '';
    public string $DIFFERENCES = '';
}

// простой helper по аналогии с ларавелевским во избежании ошибок доступа к отсутствующим элементам массива при получении элементов массива
// это не всегда нужно и не всегда оправдавданно, но в данном случае я думаю имеет смысл
// БЕЗ .dot нотации
function arr_get(mixed $array, string $key, mixed $default = null): mixed
{
    if (!is_array($array)) {
        $array = (array) $array;
    }
    $result = $array[$key] ?? $default;

    return $result;
}

class TsReturnOperation extends ReferencesOperation
{
    // Можно вынести в отдельный backed enum
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(array $data): Data // нужен более узкий тип возврата
    {
        // Данные надо получать уже в валидированном и типизированном формате,
        // иначе получается code mess с приведением типов

        $resellerId = (int) arr_get($data, 'resellerId');
        $notificationType = (int) arr_get($data, 'notificationType');
        $clientId = (int) arr_get($data, 'clientId');

        // Инит DTO с результатом
        $result = new OperationResultDto;

        // в случае пустого реселлера возвращаем результат операции с ошибкой
        if (empty($resellerId)) {
            // почему тут ошибка в notificationClientBySms->message
            $result->notificationClientBySms->message = 'Empty reseller';
            return $result;
        }

        // В случае пустого $notificationType, выбрасываем ошибку
        if (empty($notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }

        // Реселлер обязателен
        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        // Должен присутствовать клиент с типом TYPE_CUSTOMER (который почему-то единственен)
        // Получение реселерера у клиента не реализовано
        $client = Contractor::getById($clientId);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->seller->id !== $resellerId) {
            throw new \Exception('сlient not found!', 400);
        }

        // Создатель обязателен (как Господь)
        $creator = Employee::getById((int) arr_get($data, 'creatorId'));
        if ($creator === null) {
            throw new \Exception('Creator not found!', 400);
        }

        // Эксперт аналогично обязателен
        $expert = Employee::getById((int) arr_get($data, 'expertId'));
        if ($expert === null) {
            throw new \Exception('Expert not found!', 400);
        }

        // получим разницу
        $differences = $this->assembleDifferences($data, $resellerId, $notificationType);

        // соберем данные для шаблона
        $templateData = $this->assembleTemplateData($data, $creator, $expert, $client, $differences);

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            // Так норм, особо можно не улучшать (можно через array_filter/сравнение размеров массивов, но нет смысла)
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        // Пробуем отослать письмо от реселеллера
        if ($this->sendResellerEmail($resellerId, $templateData)) {
            $result->notificationClientByEmail = true;
        };

        // Шлём клиентское уведомление только если произошла смена статуса
        if ($this->isStatusChanged($notificationType, $data)) {
            // Шлем сообщение по мылу
            $result->notificationClientByEmail = $this->notifyClientByEmail($client, $templateData, $resellerId);

            // Шлем СМС на мобильный номер если есть
            if (!empty($client->mobile)) {
                try {
                    $this->notifyClientBySms(
                        $resellerId,
                        $client,
                        $data,
                        $templateData
                    );

                    $result->notificationClientBySms->isSent = true;
                }
                catch (\Exception $e) {
                    $result->notificationClientBySms->isSent = false;
                    $result->notificationClientBySms->message = $e->getMessage();
                }
            }
        }

        return $result;
    }

    private function isStatusChanged(int $notificationType, array $data): bool
    {
        return ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to']));
    }

    private function notifyClientBySms(
        int $resellerId,
        Contractor $client,
        array $data,
        TemplateDataDto $templateData
    ): bool {
        $error = null;

        $res = NotificationManager::send(
            $resellerId,
            $client->id,
            NotificationEvents::CHANGE_RETURN_STATUS,
            (int) $data['differences']['to'],
            $templateData,
            $error);

        if (!empty($error)) {
            throw new \Exception($error);
        }

        return (bool) $res;
    }

    private function notifyClientByEmail(
        Contractor $client,
        TemplateDataDto $templateData,
        int $resellerId
        ): bool
    {
        if (!empty($emailFrom) && !empty($client->email)) {
            $emailFrom = getResellerEmailFrom($resellerId);

            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo' => $client->email,
                    'subject' => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int) $data['differences']['to']);

            return true;
        }

        return false;
    }

    private function sendResellerEmail(int $resellerId, TemplateDataDto $templateData): bool
    {
        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            // подготовим массив сообщений об отсылке
            $toSend = [];
            foreach ($emails as $email) {
                $toSend[] = [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo' => $email,
                    'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ];
            }
            MessagesClient::sendMessage($toSend, $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);

            return true;
        }

        return false;
    }

    /**
     * Возвращает строку с изменением позиции (если есть)
     *
     * Либо строка "Добавлена новая позиция", либо строка
     * "Позиция изменена с "XXX" на "XXX"
     * В противном случае возвращает пустую строку
     *
     * @param array $data
     * @param integer $resellerId
     * @param integer $notificationType
     * @return string
     */
    private function assembleDifferences(
        array $data,
        int $resellerId,
        int $notificationType
    ): string {
        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === self::TYPE_CHANGE && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int) ($data['differences']['from'] ?? 0)),
                'TO' => Status::getName((int) ($data['differences']['to'] ?? 0)),
            ], $resellerId);
        } else {
            $differences = '';
        }

        return $differences;
    }

    private function assembleTemplateData(
        array $data,
        Employee $creator,
        Employee $expert,
        Contractor $client,
        string $differences,
    ): TemplateDataDto {
        $dto = new TemplateDataDto;
        $dto->COMPLAINT_ID = (int) arr_get($data, 'complaintId');
        $dto->COMPLAINT_NUMBER = (string) arr_get($data, 'complaintNumber');
        $dto->CREATOR_ID = (int) arr_get($data, 'creatorId');
        $dto->CREATOR_NAME = $creator->getFullName();
        $dto->EXPERT_ID = (int) arr_get($data, 'expertId');
        $dto->EXPERT_NAME = $expert->getFullName();
        $dto->CLIENT_ID = (int) arr_get($data, 'clientId');
        $dto->CLIENT_NAME = $client->getFullName();
        $dto->CONSUMPTION_ID = (int) arr_get($data, 'consumptionId');
        $dto->CONSUMPTION_NUMBER = (string) arr_get($data, 'consumptionNumber');
        $dto->AGREEMENT_NUMBER = (string) arr_get($data, 'agreementNumber');
        $dto->DATE = (string) arr_get($data, 'date');
        $dto->DIFFERENCES = $differences;

        return $dto;
    }
}

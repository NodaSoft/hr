<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        /**
         * @vphilippov:
         * Вместо жесткого приведения типа к (array) лучше проверить тип праметра:
         * так легче отлаживать код, если вместо массива будет передано что-то другое.
         */
        $data = $this->getRequest('data');

        if (!is_array($data)) {
            throw new \Exception('Parameter [data] is not specified or has invalid value');
        }

        $resellerId = (int)($data['resellerId'] ?? 0);
        $notificationType = $data['notificationType'] ?? null;

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        if (empty($resellerId)) {
            /**
             * @vphilippov:
             * Не понятно, почему здесь возврат результата при некорректном параметре, в то время как в остальных
             * подобных случаях бросается Exception. Возможно такая стоит задача, но больше похоже на логическую ошибку.
             * Вместо return стоит вставить throw new \Exception(...)
             */
            $result['notificationClientBySms']['message'] = 'Parameter [resellerId] is empty or has invalid value';
            return $result;
        }

        if (empty($notificationType) || !in_array($notificationType, [self::TYPE_NEW, self::TYPE_CHANGE])) {
            throw new \Exception('Parameter [notificationType] is empty or has invalid value', 400);
        }

        /**
         * @vphilippov:
         * Проверка результата getById на значение null имеет смысл только после доработки самого метода,
         * который в своем первоначальном виде не мог возвращать что-то, кроме объекта типа Contractor.
         *
         * Где-то логическая ошибка: либо метод getById должен бросать исключение, когда объект не найден
         * (тогда проверка результата на null не нужна), либо должен уметь возвращать null.
         * Выбран второй вариант чтобы сообщения Exception были более содержательными.
         */

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById($data['clientId'] ?? 0);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!', 400);
        }

        /*
         * @vphilippov:
         * Выражение не имеет смысла, т.к. если пуст fullName, то name и подавно

            $cFullName = $client->getFullName();
            if (empty($client->getFullName())) {
                $cFullName = $client->name;
            }
        */

        $creator = Employee::getById($data['creatorId'] ?? 0);
        if ($creator === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $expert = Employee::getById($data['expertId'] ?? 0);
        if ($expert === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $dataDifferences = $data['differences'] ?? [];  // expect array

        if (!is_array($dataDifferences)) {
            throw new \Exception('Parameter [differences] has invalid value');
        }

        if (isset($dataDifferences['from']) && !Status::checkId($dataDifferences['from'])) {
            throw new \Exception('Parameter [differences][from] has invalid value');
        }

        if (isset($dataDifferences['to']) && !Status::checkId($dataDifferences['to'])) {
            throw new \Exception('Parameter [differences][to] has invalid value');
        }

        $differences = '';

        /**
         * @vphilippov:
         * Функции __() в коде не объявлено, корректность вызова нельзя проверить.
         * Судя по всему, имеется ввиду какой-то шаблонизатор, возвращающий строку.
         * Поэтому исходим из того, что в вызовах этой функции нет ошибок
         */

        if ($notificationType === self::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);

        } else if ($notificationType === self::TYPE_CHANGE && !empty($dataDifferences)) {
            // Считаем, что свойства [from] и [to] являются обязательными.
            // Корректность значений проверена выше
            if (isset($dataDifferences['from']) && isset($dataDifferences['to'])) {
                // Не лишним будет так же проверить, что данные действительно поменялись
                if ($dataDifferences['from'] != $dataDifferences['to']) {
                    $differences = __('PositionStatusHasChanged', [
                        'FROM' => Status::getName($dataDifferences['from']),
                        'TO'   => Status::getName($dataDifferences['to']),
                    ],
                        $resellerId
                    );
                }
            } else {
                throw new \Exception('Parameter [differences] must have [from] and [to] properties');
            }
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => $creator->id,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expert->id,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->id,
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        foreach ($templateData as $key => $tempData) {
            if (empty($tempData)) {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }

        $emailFrom = getResellerEmailFrom($resellerId);

        // Получаем email сотрудников из настроек
        $emails = getEmailsByPermit($resellerId, TS_GOODS_RETURN);

        /**
         * @vphilippov:
         * Классы MessagesClient и NotificationManager в коде не представлены,
         * поэтому исходим из того, что вызовы их методов не содержат ошибок.
         */

        foreach ($emails as $email) {
            MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $email,
                           'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ],
                $resellerId,
                $notificationType === self::TYPE_NEW
                    ? NotificationEvents::NEW_RETURN_STATUS
                    : NotificationEvents::CHANGE_RETURN_STATUS
            );

            $result['notificationEmployeeByEmail'] = true;
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && isset($dataDifferences['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                        0 => [ // MessageTypes::EMAIL
                               'emailFrom' => $emailFrom,
                               'emailTo'   => $client->email,
                               'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                               'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                        ],
                    ],
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $dataDifferences['to']
                );

                $result['notificationClientByEmail'] = true;
            }

            if (!empty($client->mobile)) {
                $error = null;
                $res = NotificationManager::send(
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $dataDifferences['to'],
                    $templateData,
                    $error           // Видимо передача по ссылке
                );

                if ($res) {
                    $result['notificationClientBySms']['isSent'] = true;
                }
                if (!empty($error)) {
                    $result['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $result;
    }
}

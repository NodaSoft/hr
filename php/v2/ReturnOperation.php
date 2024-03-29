<?php
// PHP 8.1
namespace NW\WebService\References\Operations\Notification;

interface MessagesClientInterface {
    public function sendMessage(
        array $array,
        mixed $resellerId,
        $id,
        string $CHANGE_RETURN_STATUS,
        int $param
    ): void;
}
class MessagesClient implements MessagesClientInterface  {
    public function sendMessage(
        array $array,
        mixed $resellerId,
        $id,
        string $eventName,
        ?int $param = null
    ): void 
    {
        
    }
}

interface NotificationManagerInterface {
    public function send(
        mixed $resellerId,
        $id,
        string $CHANGE_RETURN_STATUS,
        int $param,
        array $templateData
    ): bool;

    public function lastError(): ?string;
}

class NotificationManager implements NotificationManagerInterface {
    public function send(
        mixed $resellerId,
        $id,
        string $CHANGE_RETURN_STATUS,
        int $param,
        array $templateData
    ): bool {
        return true;
    }

    public function lastError(): ?string
    {
        return null;
    }


}

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW = 1;
    public const TYPE_CHANGE = 2;

    public function __construct(
        private MessagesClient $messagesClient,
        private NotificationManager $smsNotificationService
    )
    {}

    /**
     * @return array
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data               = $this->getRequestData();
        $resellerId         = (int) $data['resellerId'];
        $notificationType   = (int) $data['notificationType'];
        $result             = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        if (empty($resellerId))
        {
            $result['notificationClientBySms']['message'] = 'Empty resellerId';

            return $result;
        }

        // Проверим что ресселер существует
        $this->validateResseler($resellerId);
        // Получим клиента/контрактора
        $client = $this->getContractor($data['clientId'], $resellerId);
        // Получим создателя
        $cr = $this->getCreator($data['creatorId']);
        // Получим эксперта
        $et = $this->getExpert($data['expertId']);

        if ($notificationType === self::TYPE_NEW)
        {
            $differences = $this->buildTemplate('NewPositionAdded', null, $resellerId);
        }
        elseif ($notificationType === self::TYPE_CHANGE)
        {
            $differences = $this->buildTemplate('PositionStatusHasChanged', [
                'FROM' => Status::getName($data['differences']['from']),
                'TO'   => Status::getName($data['differences']['to']),
            ], $resellerId);
        }
        else
        {
            // На  случай если в getRequestData поменяем код , пропустим новый тип , а реализовать забудем
            throw new \Exception("Не поддерживаемый notificationType = $notificationType" , 400);
        }

        $templateData = [
            'COMPLAINT_ID'       => (int)$data['complaintId'],
            'COMPLAINT_NUMBER'   => (string)$data['complaintNumber'],
            'CREATOR_ID'         => (int)$data['creatorId'],
            'CREATOR_NAME'       => $cr->getFullName(),
            'EXPERT_ID'          => (int)$data['expertId'],
            'EXPERT_NAME'        => $et->getFullName(),
            'CLIENT_ID'          => (int)$data['clientId'],
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int)$data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string)$data['agreementNumber'],
            'DATE'               => (string)$data['date'],
            'DIFFERENCES'        => $differences,
        ];

        // Если хоть одна переменная для шаблона не задана, то не отправляем уведомления
        // Скроем побочку от foreach
        $this->validateTemplateData($templateData);

        $emailFrom = $this->getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
        $emails = $this->getEmailsByPermit($resellerId);

        if($emailFrom)
        {
            foreach ($emails as $email)
            {
                $this->messagesClient->sendMessage([
                    [ // MessageTypes::EMAIL
                      'emailFrom' => $emailFrom,
                      'emailTo'   => $email,
                      'subject'   => $this->buildTemplate(
                          'complaintEmployeeEmailSubject',
                          $templateData,
                          $resellerId
                      ),
                      'message'   => $this->buildTemplate(
                          'complaintEmployeeEmailBody',
                          $templateData,
                          $resellerId
                      ),
                    ]
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS);

                $result['notificationEmployeeByEmail'] = true;
            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE)
        {
            if ($client->hasEmail())
            {
                $this->messagesClient->sendMessage(
                    [
                        [ // MessageTypes::EMAIL
                          'emailFrom' => $emailFrom,
                          'emailTo'   => $client->email,
                          'subject'   => $this->buildTemplate(
                              'complaintClientEmailSubject',
                              $templateData,
                              $resellerId
                          ),
                          'message'   => $this->buildTemplate(
                              'complaintClientEmailBody',
                              $templateData,
                              $resellerId
                          ),
                        ]
                    ],
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$data['differences']['to']
                );
                $result['notificationClientByEmail'] = true;
            }

            if ($client->isMobile())
            {
                $isSent = $this->smsNotificationService->send(
                    $resellerId,
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    (int)$data['differences']['to'],
                    $templateData
                );
                if ($isSent)
                {
                    $result['notificationClientBySms']['isSent'] = true;
                }
                else
                {
                    $result['notificationClientBySms']['message'] = $this->smsNotificationService->lastError() ?? 'ошибка';
                }
            }
        }

        return $result;
    }

    // Шаблонизатор/форматтирование?
    private function buildTemplate(
        string $string,
        ?array $templateData = null,
        ?int $resellerId = null
    ): string
    {
        return '';
    }

    /**
     * @param int $resellerId
     *
     * @throws \Exception
     */
    private function validateResseler(int $resellerId): void
    {
        $reseller = Seller::getById($resellerId);
        if (empty($reseller))
        {
            throw new \Exception('Seller not found!', 400);
        }
    }

    private function getContractor(int $id, int $resellerId): Contractor
    {
        $contractor = Contractor::getById($id);
        if (empty($contractor))
        {
            throw new \Exception('сlient not found!', 400);
        }
        if ($contractor->type !== Contractor::TYPE_CUSTOMER)
        {
            throw new \Exception('неверный тип у Contractor', 400);
        }
        if ($contractor->Seller->id !== $resellerId)
        {
            throw new \Exception('Contractor seller id != resellerId', 400);
        }

        if(!$contractor->hasFullName())
        {
            throw new \Exception('Contractor пустой FullName', 400);
        }

        return $contractor;
    }

    private function getCreator(int $id): Employee
    {
        $ct = Employee::getById($id);
        if (empty($cr))
        {
            throw new \Exception('Creator not found!', 400);
        }

        if(!$cr->hasFullName())
        {
            throw new \Exception('Creator (cr) пустой FullName', 400);
        }

        return $ct;
    }

    private function getExpert(int $id): Employee
    {
        $et = Employee::getById($id);
        if (empty($et))
        {
            throw new \Exception('Expert not found!', 400);
        }

        if(!$et->hasFullName())
        {
            throw new \Exception('EXPERT (et) пустой FullName', 400);
        }

        return $et;
    }

    /**
     * Гарантируем корректные входные данные
     *
     * @return array
     * @throws \Exception
     */
    private function getRequestData(): array
    {
        $data = $this->getRequest('data');
        if(!is_array($data))
        {
            throw new \Exception("data должен быть array (или объект для json)", 400);
        }

        foreach ([
            'complaintId',
            'complaintNumber',
            'creatorId',
            'expertId',
            'clientId',
            'consumptionId',
            'consumptionNumber',
            'agreementNumber',
            'date',
            'resellerId',
            'notificationType',
        ] as $k => $v)
        {
            if(empty($data[$k]))
            {
                throw new \Exception("Неверные входные данные. Нужно заполнить значение $k", 400);
            }
        }

        if(!in_array($data['notificationType'], [self::TYPE_NEW, self::TYPE_CHANGE]))
        {
            throw new \Exception('Неверный notificationType в запросе', 400);
        }

        if($data['notificationType'] === self::TYPE_CHANGE && !empty($data['differences']))
        {
            if(empty($data['differences']['from']))
            {
                throw new \Exception('data.differences.from пустой', 400);
            }
            if(empty($data['differences']['to']))
            {
                throw new \Exception('data.differences.to пустой', 400);
            }

            // Проверим корректность статусов
            if(!Status::isValid($data['differences']['from']))
            {
                throw new \Exception('data.differences.from не поддерживается', 400);
            }
            if(!Status::isValid($data['differences']['to']))
            {
                throw new \Exception('data.differences.to не поддерживается', 400);
            }
        }
        elseif($data['notificationType'] === self::TYPE_CHANGE && empty($data['differences']['to']))
        {
            throw new \Exception('требуется differences.to ', 400);
        }

        return $data;
    }

    /**
     * Скроем побочные эффекты от foreach
     * @param array $templateData
     *
     * @throws \Exception
     */
    private function validateTemplateData(array $templateData): void
    {
        foreach ($templateData as $key => $tempData)
        {
            if (empty($tempData))
            {
                throw new \Exception("Template Data ({$key}) is empty!", 500);
            }
        }
    }

    /**
     * @param int $resellerId
     *
     * @return string
     */
    private function getResellerEmailFrom(int $resellerId): string
    {
        return getResellerEmailFrom($resellerId);
    }

    private function getEmailsByPermit(int $resellerId): array
    {
        return getEmailsByPermit($resellerId, 'tsGoodsReturn');
    }
}

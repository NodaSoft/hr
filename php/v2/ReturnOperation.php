<?php

namespace NW\WebService\References\Operations\Notification;

class TsReturnOperation extends ReferencesOperation
{
    /**
     * Здесь лучше Enum
     * но не знаю какая версия PHP подразумевалась
     */
    const TYPE_NEW    = 1;
    const TYPE_CHANGE = 2;

    /**
     * Название конечно...))
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $data = (array) $this->getRequest('data');
        self::validateRequestData($data, [
            'resellerId',
            'clientId',
            'notificationType',
            'creatorId',
            'expertId',
            'complaintId',
            'complaintNumber',
            'consumptionId',
            'consumptionNumber',
            'agreementNumber',
            'date',
        ]);

        $result = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail'   => false,
            'notificationClientBySms'     => [
                'isSent'  => false,
                'message' => '',
            ],
        ];

        $notificationType  = (int) $data['notificationType'];

        $reseller = Seller::getById((int) $data['resellerId']);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int) $data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $reseller->id) {
            throw new \Exception('Client not found!', 400);
        }

        $creator = Employee::getById((int) $data['creatorId']);
        if ($creator === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $expert = Employee::getById((int) $data['expertId']);
        if ($expert === null) {
            throw new \Exception('Expert not found!', 400);
        }


        $templateData = [
            'COMPLAINT_ID'       => (int)    $data['complaintId'],
            'COMPLAINT_NUMBER'   => (string) $data['complaintNumber'],
            'CREATOR_ID'         => $creator->id,
            'CREATOR_NAME'       => $creator->getFullName(),
            'EXPERT_ID'          => $expert->id,
            'EXPERT_NAME'        => $expert->getFullName(),
            'CLIENT_ID'          => $client->id,
            'CLIENT_NAME'        => $client->getFullName(),
            'CONSUMPTION_ID'     => (int)    $data['consumptionId'],
            'CONSUMPTION_NUMBER' => (string) $data['consumptionNumber'],
            'AGREEMENT_NUMBER'   => (string) $data['agreementNumber'],
            'DATE'               => (string) $data['date'],
            'DIFFERENCES'        => self::getDifferences($notificationType, $reseller->id, $data['differences'] ?? []),
        ];
        self::validateTemplateData($templateData);


        $emailFrom = getResellerEmailFrom($reseller->id);

        $result['notificationEmployeeByEmail'] = self::sendEmailsToEmployee($emailFrom, $templateData, $reseller->id);


        $differencesTo = null;
        if (isset($data['differences']) && isset($data['differences']['to'])) {
            $differencesTo = (int) $data['differences']['to'];
        }
        if ($notificationType === self::TYPE_CHANGE && $differencesTo) {
            $result['notificationClientByEmail'] = self::sendEmailToClient(
                $emailFrom,
                $client,
                $templateData,
                $differencesTo
            );

            try {
                $result['notificationClientBySms']['isSent'] = (bool) self::sendSmsToClient($client, $templateData, $differencesTo);
            } catch (\Exception $e) {
                $result['notificationClientBySms']['message'] = $e->getMessage();
            }
        }

        return $result;
    }

    private static function getDifferences($notificationType, $resellerId, array $differences = []): string
    {
        switch ($notificationType)
        {
            case self::TYPE_NEW:
                return __('NewPositionAdded', null, $resellerId);

            case self::TYPE_CHANGE:
                if (isset($differences['from']) && isset($differences['to'])) {
                    return __('PositionStatusHasChanged', [
                        'FROM' => Status::getName((int) $differences['from']),
                        'TO'   => Status::getName((int) $differences['to']),
                    ], $resellerId);
                } else {
                    return '';
                }

            default:
                return '';
        }
    }

    private static function sendEmailsToEmployee($emailFrom, $templateData, $resellerId): bool
    {
        $emails = getEmailsByPermit($resellerId, 'tsGoodsReturn');
        if (!empty($emailFrom) && count($emails) > 0) {
            foreach ($emails as $email) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                        'emailFrom' => $emailFrom,
                        'emailTo'   => $email,
                        'subject'   => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                        'message'   => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, NotificationEvents::CHANGE_RETURN_STATUS);
            }

            return true;
        }

        return false;
    }

    private static function sendEmailToClient($emailFrom, Contractor $client, $templateData, int $differencesTo): bool
    {
        if (!empty($emailFrom) && !empty($client->email)) {
            $resellerId = $client->Seller->id;
            MessagesClient::sendMessage([
                0 => [ // MessageTypes::EMAIL
                    'emailFrom' => $emailFrom,
                    'emailTo'   => $client->email,
                    'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                    'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                ],
            ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo);
            return true;
        }

        return false;
    }

    /**
     * @throws \Exception
     */
    private static function sendSmsToClient(Contractor $client, $templateData, int $differencesTo)
    {
        if (!empty($client->mobile)) {
            $error = null;
            $result = NotificationManager::send(
                $client->Seller->id,
                $client->id,
                NotificationEvents::CHANGE_RETURN_STATUS,
                $differencesTo,
                $templateData,
                $error
            );
            if ($error) {
                throw new \Exception($error);
            }

            return $result;
        }

        return null;
    }

    /**
     * @throws \Exception
     */
    private static function validateRequestData($requestData, array $requiredKeys)
    {
        if ($emptyKeys = self::getEmptyKeys($requestData, $requiredKeys)) {
            $emptyValuesImploded = implode(', ', $emptyKeys);
            throw new \Exception("Request Data ({$emptyValuesImploded}) is empty!", 400);
        }
    }

    /**
     * @throws \Exception
     */
    private static function validateTemplateData($templateData)
    {
        if ($emptyKeys = self::getEmptyKeys($templateData)) {
            $emptyValuesImploded = implode(', ', $emptyKeys);
            throw new \Exception("Template Data ({$emptyValuesImploded}) is empty!", 500);
        }
    }

    private static function getEmptyKeys(array $data, array $requiredKeys = []): array
    {
        if ($requiredKeys) {
            $data = array_filter($data, function ($key) use($requiredKeys) {
                return in_array($key, $requiredKeys);
            }, ARRAY_FILTER_USE_KEY);
        }

        return array_keys(
            array_filter($data, function ($value) {
                return empty($value);
            })
        );
    }
}

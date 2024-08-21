<?php

namespace NW\WebService\References\Operations\Notification;

use http\Exception\InvalidArgumentException;

class TsReturnOperation extends ReferencesOperation
{
    protected static function initOperationModelsFromRequest(array $request): array
    {
        $models = [
            'resellerId' => Seller::class,
            'clientId' => Contractor::class,
            'creatorId' => Employee::class,
            'expertId' => Employee::class,
        ];
        foreach ($models as $key => &$model) {
            // все проверки результата getById надо выносить в конкретные реализации в моделях
            $model = call_user_func_array([$model, 'getById'], $request[$key]);
            if (is_null($model)) {
                throw new \Exception("Model {$model} by id {$key} not found", 400);
            }
        }
        unset($model);
        return $models;
    }

    protected static function compileTemplateDataFromRequest(array $request, array $models, string $differences): array
    {
        /** @var Employee $client */
        $creator = $models['creatorId'];
        /** @var Employee $expert */
        $expert = $models['expertId'];
        /** @var Contractor $client */
        $client = $models['creatorId'];
        return [
            'COMPLAINT_ID' => $request['complaintId'],
            'COMPLAINT_NUMBER' => (string)$request['complaintNumber'],
            'CREATOR_ID' => $request['creatorId'],
            'CREATOR_NAME' => $creator->getFullName(),
            'EXPERT_ID' => $request['expertId'],
            'EXPERT_NAME' => $expert->getFullName(),
            'CLIENT_ID' => $request['clientId'],
            'CLIENT_NAME' => $client->getFullName(),
            'CONSUMPTION_ID' => $request['consumptionId'],
            'CONSUMPTION_NUMBER' => (string)$request['consumptionNumber'],
            'AGREEMENT_NUMBER' => (string)$request['agreementNumber'],
            'DATE' => (string)$request['date'],
            'DIFFERENCES' => $differences,
        ];
    }

    protected static function notifyByEmail(array &$operationResult, array $request, array $models, array $template): void
    {
        $emailFrom = getResellerEmailFrom($request['resellerId']);
        $events = [
            'employee' => [
                'subject' => 'complaintEmployeeEmailSubject',
                'message' => 'complaintEmployeeEmailBody',
                'emails' => getEmailsByPermit($request['resellerId'], 'tsGoodsReturn'),
                'diff' => null,
                'flag' => 'notificationEmployeeByEmail'
            ],
            'client' => [
                'subject' => 'complaintClientEmailSubject',
                'message' => 'complaintClientEmailBody',
                'emails' => [$models['clientId']->email],
                'diff' => $request['differences']['to'],
                'flag' => 'notificationClientByEmail'
            ]
        ];
        if (!empty($emailFrom)) {
            foreach ($events as $subject => $event) {
                // сбросим цикл, чтобы не слать клиенту, когда нет изменений
                if ($subject === 'client' && $request['notificationType'] !== self::TYPE_CHANGE) {
                    break;
                }
                $sent = [];
                foreach ($event['emails'] as $email) {
                     $sent[] = MessagesClient::sendMessage([
                        [
                            'emailFrom' => $emailFrom,
                            'emailTo' => $email,
                            'subject' => __($event['subject'], $template, $request['resellerId']),
                            'message' => __($event['message'], $template, $request['resellerId']),
                        ],
                    ], $request['resellerId'], NotificationEvents::CHANGE_RETURN_STATUS, $event['diff']);
                }
                // если хотя бы что-то ушло - уже хорошо
                if (in_array(true, $sent)) {
                    $operationResult[$event['flag']] = true;
                }
            }

        }
    }

    /**
     * @throws \Exception
     */
    public function doOperation(): array
    {
        $operationResult = [
            'notificationEmployeeByEmail' => false,
            'notificationClientByEmail' => false,
            'notificationClientBySms' => [
                'isSent' => false,
                'message' => '',
            ],
        ];
        $request = $this->getRequest('data');
        $models = self::initOperationModelsFromRequest($request);
        $differences = call_user_func(function (array $request) {
            switch ($request['notificationType']) {
                case self::TYPE_NEW:
                    return __('NewPositionAdded', null, $request['resellerId']);
                case self::TYPE_CHANGE:
                    return __('PositionStatusHasChanged', [
                        'FROM' => Status::getName($request['differences']['from']),
                        'TO' => Status::getName($request['differences']['to']),
                    ], $request['resellerId']);
                default:
                    throw new InvalidArgumentException(
                        'Unsupported notificationType: ' . $request['notificationType'],
                        422
                    );
            }
        }, $request);
        $templateData = self::compileTemplateDataFromRequest($request, $models, $differences);
        self::notifyByEmail($operationResult, $request, $models, $templateData);
        // не вижу смысла выносить в отдельный метод на уведомление по мобиле, и так все понятно
        if ($request['notificationType'] === self::TYPE_CHANGE) {
            $client = $models['clientId'];
            $error = [];
            if (isset($client->mobile)) {
                $operationResult['notificationClientBySms']['isSent'] = NotificationManager::send(
                    $request['resellerId'],
                    $client->id,
                    NotificationEvents::CHANGE_RETURN_STATUS,
                    $request['differences']['to'],
                    $templateData,
                    $error
                );
                if (!empty($error)) {
                    $operationResult['notificationClientBySms']['message'] = $error;
                }
            }
        }

        return $operationResult;
    }
}

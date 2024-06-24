<?php

namespace NW\WebService\References\Operations\Notification;

use Exception;

class ReturnOperation extends ReferencesOperation
{
    /**
     * @throws Exception
     */
    public function doOperation(): array
    {
        $request = ReturnOperationRequest::fromArray($this->getRequest('data'));

        if ($request->notificationType === NotificationTypeEnum::TYPE_CHANGE->value && !$request->differences) {
            throw new Exception('Differences must be non-empty if notificationType = 2', 400);
        }

        if (is_null(Seller::getById($request->resellerId))) {
            throw new Exception('Seller not found!', 400);
        }

        if (
            is_null($client = Contractor::getById($request->clientId))
            || $client->type !== Contractor::TYPE_CUSTOMER
            || $client->Seller->id !== $request->resellerId
        ) {
            throw new Exception('Client not found!', 400);
        }

        if (is_null($cr = Employee::getById($request->creatorId))) {
            throw new Exception('Creator not found!', 400);
        }

        if (is_null($et = Employee::getById($request->expertId))) {
            throw new Exception('Expert not found!', 400);
        }

        $emailFrom = getResellerEmailFrom($request->resellerId);
        $emails = getEmailsByPermit($request->resellerId, 'tsGoodsReturn');
        $notificationClientByEmail = false;
        $notificationEmployeeByEmail = false;
        $notificationClientBySms = ['isSent' => false, 'message' => ''];

        $templateData = [
            'COMPLAINT_ID' => $request->complaintId,
            'COMPLAINT_NUMBER' => $request->complaintNumber,
            'CREATOR_ID' => $cr->id,
            'CREATOR_NAME' => $cr->getFullName(),
            'EXPERT_ID' => $et->id,
            'EXPERT_NAME' => $et->getFullName(),
            'CLIENT_ID' => $client->id,
            'CLIENT_NAME' => $client->getFullName(),
            'CONSUMPTION_ID' => $request->consumptionId,
            'CONSUMPTION_NUMBER' => $request->consumptionNumber,
            'AGREEMENT_NUMBER' => $request->agreementNumber,
            'DATE' => $request->date,
            'DIFFERENCES' => $this->getDifferences(
                $request->notificationType,
                $request->resellerId,
                $request->differences
            ),
        ];

        if (!empty($emailFrom) && !empty($emails)) {
            $notificationEmployeeByEmail = $this->sendEmail(
                $emailFrom,
                $emails,
                $templateData,
                $request->resellerId
            );
        }

        if (
            $request->notificationType === NotificationTypeEnum::TYPE_CHANGE->value
            && !empty($request->differences['to'])
        ) {
            if (!empty($emailFrom) && !empty($client->email)) {
                $notificationClientByEmail = $this->sendEmail(
                    $emailFrom,
                    $client->email,
                    $templateData,
                    $request->resellerId,
                    $client->id,
                    $request->differences['to']
                );
            }

            if (!empty($client->phone)) {
                $notificationClientBySms = $this->sendSMS(
                    $request->resellerId,
                    $request->clientId,
                    $request->differences['to'],
                    $templateData
                );
            }
        }

        return [
            'notificationEmployeeByEmail' => $notificationEmployeeByEmail,
            'notificationClientByEmail' => $notificationClientByEmail,
            'notificationClientBySms' => $notificationClientBySms,
        ];
    }

    /**
     * @throws Exception
     */
    private function getDifferences(int $notificationType, int $resellerId, ?array $differences): string
    {
        return match (true) {
            $notificationType === NotificationTypeEnum::TYPE_NEW->value => __('NewPositionAdded', null, $resellerId),
            $notificationType === NotificationTypeEnum::TYPE_CHANGE->value && !empty($data['differences']) =>
            __('PositionStatusHasChanged', [
                'FROM' => StatusEnum::tryFrom($differences['from']),
                'TO' => StatusEnum::tryFrom($differences['to']),
            ], $resellerId),
            default => throw new \Exception("Template Data (DIFFERENCES) is empty!", 500)
        };
    }

    private function sendEmail(
        string       $emailFrom,
        array|string $emailsTo,
        array        $templateData,
        int          $resellerId,
        ?int         $clientId = null,
        ?int         $differencesTo = null,
    ): bool
    {
        $emails = is_array($emailsTo) ? $emailsTo : [$emailsTo];
        $params = [$resellerId, $clientId, NotificationEvents::CHANGE_RETURN_STATUS, $differencesTo];
        $result = false;

        foreach ($emails as $email) {
            try {
                MessagesClient::sendMessage([[
                    'emailFrom' => $emailFrom,
                    'emailTo' => $email,
                    'subject' => __('complaintEmployeeEmailSubject', $templateData, $resellerId),
                    'message' => __('complaintEmployeeEmailBody', $templateData, $resellerId),
                ]], ...array_filter($params));

                $result = true;
            } catch (Exception $e) {
                Log::notice($e);
            }
        }

        return $result;
    }

    private function sendSMS(int $resellerId, int $clientId, int $differencesTo, array $templateData): array
    {
        $message = '';

        $isSent = NotificationManager::send(
            $resellerId,
            $clientId,
            NotificationEvents::CHANGE_RETURN_STATUS,
            $differencesTo,
            $templateData,
            $message
        );

        return ['isSent' => $isSent, 'message' => $message];
    }
}
<?php

namespace NW\WebService\References\Operations\Notification;

use NodaSoft\Factory\Dto\TsReturnDtoFactory;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;

class TsReturnOperation extends ReferencesOperation
{
    public const TYPE_NEW    = 1;
    public const TYPE_CHANGE = 2;

    /** @var TsReturnOperationResult */
    private $result;

    public function __construct()
    {
        $this->result = new TsReturnOperationResult();
    }

    /**
     * @throws \Exception
     * @return TsReturnOperationResult
     */
    public function doOperation(): ReferencesOperationResult
    {
        $data = (array)$this->getRequest('data');
        $factory = new TsReturnDtoFactory();
        $dto = $factory->makeTsReturnDto($data);
        $resellerId = $data['resellerId'];
        $notificationType = (int)$data['notificationType'];

        if (empty((int)$resellerId)) {
            $this->result->setClientSmsErrorMessage('Empty resellerId');
            return $this->result;
        }

        if (empty((int)$notificationType)) {
            throw new \Exception('Empty notificationType', 400);
        }

        $reseller = Seller::getById((int)$resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!', 400);
        }

        $client = Contractor::getById((int)$data['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('сlient not found!', 400);
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->name;
        }

        $cr = Employee::getById((int)$data['creatorId']);
        if ($cr === null) {
            throw new \Exception('Creator not found!', 400);
        }

        $et = Employee::getById((int)$data['expertId']);
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
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

        $dto->setCreatorName($cr->getFullName());
        $dto->setExpertName($et->getFullName());
        $dto->setClientName($cFullName);
        $dto->setDifferences($differences);

        if (! $dto->isValid()) {
            $emptyKey = $dto->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!", 500);
        }

        $templateData = $dto->toArray();

        $emailFrom = getResellerEmailFrom($resellerId);
        // Получаем email сотрудников из настроек
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
                $this->result->markEmployeeEmailSent();

            }
        }

        // Шлём клиентское уведомление, только если произошла смена статуса
        if ($notificationType === self::TYPE_CHANGE && !empty($data['differences']['to'])) {
            if (!empty($emailFrom) && !empty($client->email)) {
                MessagesClient::sendMessage([
                    0 => [ // MessageTypes::EMAIL
                           'emailFrom' => $emailFrom,
                           'emailTo'   => $client->email,
                           'subject'   => __('complaintClientEmailSubject', $templateData, $resellerId),
                           'message'   => __('complaintClientEmailBody', $templateData, $resellerId),
                    ],
                ], $resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to']);
                $this->result->markClientEmailSent();
            }

            if (!empty($client->mobile)) {
                $res = NotificationManager::send($resellerId, $client->id, NotificationEvents::CHANGE_RETURN_STATUS, (int)$data['differences']['to'], $templateData, $error);
                if ($res) {
                    $this->result->markClientSmsSent();
                }
                if (!empty($error)) {
                    $this->result->setClientSmsErrorMessage($error);
                }
            }
        }

        return $this->result;
    }
}

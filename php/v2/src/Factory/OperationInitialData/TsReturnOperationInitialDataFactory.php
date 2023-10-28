<?php

namespace NodaSoft\Factory\OperationInitialData;

use NW\WebService\References\Operations\Notification\Contractor;
use NW\WebService\References\Operations\Notification\Employee;
use NW\WebService\References\Operations\Notification\Seller;

use NodaSoft\Factory\Dto\TsReturnDtoFactory;
use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\OperationInitialData\TsReturnOperationInitialData;
use NW\WebService\References\Operations\Notification\Status;
use NW\WebService\References\Operations\Notification\TsReturnOperation;
use function NW\WebService\References\Operations\Notification\__;

class TsReturnOperationInitialDataFactory implements OperationInitialDataFactory
{
    /**
     * @return TsReturnOperationInitialData
     */
    public function makeInitialData(array $params): OperationInitialData
    {

        //todo: set error codes 400 and 500 as it was

        $resellerId = (int) $params['resellerId'];
        $notificationType = (int) $params['notificationType'];

        if (empty($resellerId)) {
            throw new \Exception('Empty resellerId');
        }

        if (empty($notificationType)) {
            throw new \Exception('Empty notificationType');
        }

        $reseller = Seller::getById($resellerId);
        if ($reseller === null) {
            throw new \Exception('Seller not found!');
        }

        $client = Contractor::getById((int) $params['clientId']);
        if ($client === null || $client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $resellerId) {
            throw new \Exception('Client not found!');
        }

        $cFullName = $client->getFullName();
        if (empty($client->getFullName())) {
            $cFullName = $client->name;
        }

        $cr = Employee::getById((int) $params['creatorId']);
        if ($cr === null) {
            throw new \Exception('Creator not found!');
        }

        $et = Employee::getById((int) $params['expertId']);
        if ($et === null) {
            throw new \Exception('Expert not found!', 400);
        }

        $differences = '';
        if ($notificationType === TsReturnOperation::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $resellerId);
        } elseif ($notificationType === TsReturnOperation::TYPE_CHANGE
            && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $resellerId);
        }



        $templateFactory = new TsReturnDtoFactory();
        $messageTemplate = $templateFactory->makeTsReturnDto($params);
        $messageTemplate->setCreatorName($cr->getFullName());
        $messageTemplate->setExpertName($et->getFullName());
        $messageTemplate->setClientName($cFullName);
        $messageTemplate->setDifferences($differences);
        if (! $messageTemplate->isValid()) {
            $emptyKey = $messageTemplate->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!");
        }

        $data = new TsReturnOperationInitialData();
        $data->setMessageTemplate($messageTemplate);
        $data->setResellerId($reseller->id);
        $data->setNotificationType($notificationType);
    }
}

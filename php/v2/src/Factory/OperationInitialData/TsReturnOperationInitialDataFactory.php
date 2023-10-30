<?php

namespace NodaSoft\Factory\OperationInitialData;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Mapper\ClientMapper;
use NodaSoft\DataMapper\Mapper\EmployeeMapper;
use NodaSoft\DataMapper\Mapper\ResellerMapper;

use NodaSoft\Factory\Dto\TsReturnDtoFactory;
use NodaSoft\OperationParams\OperationParams;
use NodaSoft\OperationParams\TsReturnOperationParams;
use NodaSoft\OperationInitialData\OperationInitialData;
use NodaSoft\OperationInitialData\TsReturnOperationInitialData;
use NW\WebService\References\Operations\Notification\Status;
use NW\WebService\References\Operations\Notification\TsReturnOperation;
use function NW\WebService\References\Operations\Notification\__;

class TsReturnOperationInitialDataFactory implements OperationInitialDataFactory
{
    /**
     * @param TsReturnOperationParams $params
     * @return TsReturnOperationInitialData
     */
    public function makeInitialData(OperationParams $params): OperationInitialData
    {
        //todo: set error codes 400 and 500 as it was

        $notificationType = $params->getNotificationType();

        try {
            $reseller = $this->getReseller($params->getResellerId());
            $client = $this->getClient($params->getClientId(), $reseller);
            $creator = $this->getCreator($params->getCreatorId());
            $expert = $this->getExpert($params->getExpertId());
        } catch (\Exception $e) {
            throw new \Exception($e->getMessage());
        }

        $differences = '';
        if ($notificationType === TsReturnOperation::TYPE_NEW) {
            $differences = __('NewPositionAdded', null, $params->getResellerId());
        } elseif ($notificationType === TsReturnOperation::TYPE_CHANGE
            && !empty($data['differences'])) {
            $differences = __('PositionStatusHasChanged', [
                'FROM' => Status::getName((int)$data['differences']['from']),
                'TO'   => Status::getName((int)$data['differences']['to']),
            ], $params->getResellerId());
        }

        $templateFactory = new TsReturnDtoFactory();
        $messageTemplate = $templateFactory->makeTsReturnDto($params);
        $messageTemplate->setCreatorName($creator->getFullName());
        $messageTemplate->setExpertName($expert->getFullName());
        $messageTemplate->setClientName($client->getFullName());
        $messageTemplate->setDifferences($differences);

        if (! $messageTemplate->isValid()) {
            $emptyKey = $messageTemplate->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!");
        }

        $data = new TsReturnOperationInitialData();
        $data->setMessageTemplate($messageTemplate);
        $data->setReseller($reseller);
        $data->setNotificationType($notificationType);

        return $data;
    }

    public function getReseller(int $resellerId): Reseller
    {
        $resellerMapper = new ResellerMapper();
        $reseller = $resellerMapper->getById($resellerId);
        if (is_null($reseller)) {
            throw new \Exception('Reseller not found!');
        }
        return $reseller;
    }

    public function getClient(int $clientId, Reseller $reseller): Client
    {
        $clientMapper = new ClientMapper();
        $client = $clientMapper->getById($clientId);
        if (is_null($client)
            || $client->isCustomer()
            || $client->hasReseller($reseller)
        ) {
            throw new \Exception('Client not found!');
        }
        return $client;
    }

    public function getCreator(int $creatorId): Employee
    {
        $employeeMapper = new EmployeeMapper(); // todo: duplication of EmployeeMapper initialization
        $creator = $employeeMapper->getById($creatorId);
        if (is_null($creator)) {
            throw new \Exception('Creator not found!');
        }
        return $creator;
    }

    public function getExpert(int $expertId): Employee
    {
        $employeeMapper = new EmployeeMapper(); // todo: duplication of EmployeeMapper initialization
        $expert = $employeeMapper->getById($expertId);
        if (is_null($expert)) {
            throw new \Exception('Expert not found!');
        }
        return $expert;
    }
}

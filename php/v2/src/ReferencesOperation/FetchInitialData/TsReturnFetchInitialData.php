<?php

namespace NodaSoft\ReferencesOperation\FetchInitialData;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ClientMapper;
use NodaSoft\DataMapper\Mapper\EmployeeMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
use NodaSoft\DataMapper\Mapper\ResellerMapper;
use NodaSoft\Factory\Dto\TsReturnDtoFactory;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Params\TsReturnOperationParams;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\TsReturnInitialData;

class TsReturnFetchInitialData implements FetchInitialData
{
    /** @var MapperFactory */
    private $mapperFactory;

    /**
     * @param TsReturnOperationParams $params
     * @return TsReturnInitialData
     */
    public function fetch(ReferencesOperationParams $params): InitialData
    {
        //todo: set error codes 400 and 500 as it was

        try {
            $reseller = $this->getReseller($params->getResellerId());
            $client = $this->getClient($params->getClientId(), $reseller);
            $creator = $this->getCreator($params->getCreatorId());
            $expert = $this->getExpert($params->getExpertId());
            $employees = $this->getEmployees($params->getResellerId());
        } catch (\Exception $e) {
            throw new \Exception($e->getMessage());
        }

        $notification = $this->getNotification($params->getNotificationType());

        $templateFactory = new TsReturnDtoFactory();
        $messageTemplate = $templateFactory->makeTsReturnDto($params);
        $messageTemplate->setCreatorName($creator->getFullName());
        $messageTemplate->setExpertName($expert->getFullName());
        $messageTemplate->setClientName($client->getFullName());
        $messageTemplate->setDifferences($notification->composeMessage($params));

        if (! $messageTemplate->isValid()) {
            $emptyKey = $messageTemplate->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!");
        }

        $data = new TsReturnInitialData();
        $data->setMessageTemplate($messageTemplate);
        $data->setReseller($reseller);
        $data->setNotification($notification);
        $data->setDifferencesFrom($params->getDifferencesFrom());
        $data->setDifferencesTo($params->getDifferencesTo());
        $data->setClient($client);
        $data->setEmployees($employees);

        return $data;
    }

    public function getReseller(int $resellerId): Reseller
    {
        /** @var ResellerMapper $resellerMapper */
        $resellerMapper = $this->mapperFactory->getMapper('Reseller');
        $reseller = $resellerMapper->getById($resellerId);
        if (is_null($reseller)) {
            throw new \Exception('Reseller not found!');
        }
        return $reseller;
    }

    public function getClient(int $clientId, Reseller $reseller): Client
    {
        /** @var ClientMapper $clientMapper */
        $clientMapper = $this->mapperFactory->getMapper('Client');
        $client = $clientMapper->getById($clientId); //todo: replace condition with getter by filter if it's needed
        if (is_null($client)
            || ! $client->isCustomer()
            || ! $client->hasReseller($reseller)) {
            throw new \Exception('Client not found!');
        }
        return $client;
    }

    public function getCreator(int $creatorId): Employee
    {
        /** @var EmployeeMapper $employeeMapper */
        $employeeMapper = $this->mapperFactory->getMapper('Employee'); // todo: duplication of EmployeeMapper initialization
        $creator = $employeeMapper->getById($creatorId);
        if (is_null($creator)) {
            throw new \Exception('Creator not found!');
        }
        return $creator;
    }

    public function getExpert(int $expertId): Employee
    {
        /** @var EmployeeMapper $employeeMapper */
        $employeeMapper = $this->mapperFactory->getMapper('Employee'); // todo: duplication of EmployeeMapper initialization
        $expert = $employeeMapper->getById($expertId);
        if (is_null($expert)) {
            throw new \Exception('Expert not found!');
        }
        return $expert;
    }

    public function setMapperFactory(MapperFactory $mapperFactory): void
    {
        $this->mapperFactory = $mapperFactory;
    }

    /**
     * @param int $resellerId
     * @return Employee[]
     */
    public function getEmployees(int $resellerId): array
    {
        /** @var EmployeeMapper $employeeMapper */
        $employeeMapper = $this->mapperFactory->getMapper('Employee'); // todo: duplication of EmployeeMapper initialization
        return $employeeMapper->getAllByReseller($resellerId);
    }

    /**
     * @param TsReturnOperationParams $params
     * @return Notification
     * @throws \Exception
     */
    public function getNotification(int $id): ?Notification
    {
        /** @var NotificationMapper $notificationMapper */
        $notificationMapper = $this->mapperFactory->getMapper('Notification');
        return $notificationMapper->getById($id);
    }
}

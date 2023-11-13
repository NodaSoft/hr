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
use NodaSoft\GenericDto\Factory\GenericDtoFactory;
use NodaSoft\GenericDto\Dto\ReturnOperationStatusChangedMessageBodyList;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use NodaSoft\ReferencesOperation\Params\ReturnOperationStatusChangedParams;
use NodaSoft\ReferencesOperation\InitialData\InitialData;
use NodaSoft\ReferencesOperation\InitialData\ReturnOperationStatusChangedInitialData;

class ReturnOperationStatusChangedFetchInitialData implements FetchInitialData
{
    /** @var MapperFactory */
    private $mapperFactory;

    public function setMapperFactory(MapperFactory $mapperFactory): void
    {
        $this->mapperFactory = $mapperFactory;
    }

    /**
     * @param ReturnOperationStatusChangedParams $params
     * @return ReturnOperationStatusChangedInitialData
     */
    public function fetch(ReferencesOperationParams $params): InitialData
    {
        try {
            $reseller = $this->getReseller($params->getResellerId());
            $client = $this->getClient($params->getClientId(), $reseller);
            $creator = $this->getEmployee($params->getCreatorId());
            $expert = $this->getEmployee($params->getExpertId());
            $notification = $this->getNotification($params->getNotificationType());
        } catch (\Exception $e) {
            throw new \Exception("An entity was not found.", 400, $e);
        }


        $templateFactory = new GenericDtoFactory();
        /** @var ReturnOperationStatusChangedMessageBodyList $messageTemplate */
        $messageTemplate = $templateFactory->fillDtoParams(
            new ReturnOperationStatusChangedMessageBodyList(),
            $params
        );
        $messageTemplate->setCreatorName($creator->getFullName());
        $messageTemplate->setExpertName($expert->getFullName());
        $messageTemplate->setClientName($client->getFullName());
        $messageTemplate->setStatement($notification->composeMessage($params));

        if (! $messageTemplate->isValid()) {
            $emptyKey = $messageTemplate->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!", 500);
        }

        $data = new ReturnOperationStatusChangedInitialData();
        $data->setMessageTemplate($messageTemplate);
        $data->setReseller($reseller);
        $data->setNotification($notification);
        $data->setDifferencesFrom($params->getDifferencesFrom());
        $data->setDifferencesTo($params->getDifferencesTo());
        $data->setClient($client);
        $data->setEmployees($reseller->getEmployees());

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

    public function getEmployee(int $creatorId): Employee
    {
        /** @var EmployeeMapper $employeeMapper */
        $employeeMapper = $this->mapperFactory->getMapper('Employee'); // todo: duplication of EmployeeMapper initialization
        $creator = $employeeMapper->getById($creatorId);
        if (is_null($creator)) {
            throw new \Exception('Employee not found!');
        }
        return $creator;
    }

    /**
     * @param int $resellerId
     * @return Employee[]
     */
    public function getEmployees(int $resellerId): array
    {
        /** @var EmployeeMapper $employeeMapper */
        $employeeMapper = $this->mapperFactory->getMapper('Employee'); // todo: duplication of EmployeeMapper initialization
        $employee = $employeeMapper->getAllByReseller($resellerId);
        if (is_null($employee)) {
            throw new \Exception('Employee not found!');
        }
        return $employee;
    }

    public function getNotification(int $id): Notification
    {
        /** @var NotificationMapper $notificationMapper */
        $notificationMapper = $this->mapperFactory->getMapper('Notification');
        $notification = $notificationMapper->getById($id);
        if (is_null($notification)) {
            throw new \Exception('Notification not found!');
        }
        return $notification;
    }
}

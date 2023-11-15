<?php

namespace NodaSoft\ReferencesOperation\FetchInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ComplaintMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
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
        /** @var NotificationMapper $notificationMapper */
        $notificationMapper = $this->mapperFactory->getMapper('Notification');
        $notification = $notificationMapper->getByName('complaint status changed');

        if (is_null($notification)) {
            throw new \Exception('Notification was not found!', 500);
        }

        /** @var ComplaintMapper $complaintMapper */
        $complaintMapper = $this->mapperFactory->getMapper('Complaint');
        $complaint = $complaintMapper->getById($params->getComplaintId());

        if (is_null($complaint)) {
            throw new \Exception('Complaint was not found!', 400);
        }

        $reseller = $complaint->getReseller();
        $client = $complaint->getClient();
        $creator = $complaint->getCreator();
        $expert = $complaint->getExpert();
        $employees = $reseller->getEmployees();
        $currentStatus = $complaint->getStatus();
        $previousStatus = $complaint->getPreviousStatus();

        $templateFactory = new GenericDtoFactory();
        /** @var ReturnOperationStatusChangedMessageBodyList $messageTemplate */
        $messageTemplate = $templateFactory->fillDtoParams(
            new ReturnOperationStatusChangedMessageBodyList(),
            $params
        );
        $messageTemplate->setCreatorName($creator->getFullName());
        $messageTemplate->setExpertName($expert->getFullName());
        $messageTemplate->setClientName($client->getFullName());
        $messageTemplate->setCurrentStatus($currentStatus->getName());
        $messageTemplate->setPreviousStatus($previousStatus->getName());

        if (! $messageTemplate->isValid()) {
            $emptyKey = $messageTemplate->getEmptyKeys()[0];
            throw new \Exception("Template Data ({$emptyKey}) is empty!", 500);
        }

        $data = new ReturnOperationStatusChangedInitialData();
        $data->setMessageTemplate($messageTemplate);
        $data->setReseller($reseller);
        $data->setNotification($notification);
        $data->setClient($client);
        $data->setEmployees($employees);

        return $data;
    }
}

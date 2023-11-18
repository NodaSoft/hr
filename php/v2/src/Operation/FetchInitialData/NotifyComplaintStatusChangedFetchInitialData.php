<?php

namespace NodaSoft\Operation\FetchInitialData;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ComplaintMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
use NodaSoft\Operation\InitialData\InitialData;
use NodaSoft\Operation\InitialData\NotifyComplaintStatusChangedInitialData;
use NodaSoft\Request\Request;

class NotifyComplaintStatusChangedFetchInitialData implements FetchInitialData
{
    /** @var MapperFactory */
    private $mapperFactory;

    public function setMapperFactory(MapperFactory $mapperFactory): void
    {
        $this->mapperFactory = $mapperFactory;
    }

    /**
     * @param Request $request
     * @return NotifyComplaintStatusChangedInitialData
     */
    public function fetch(Request $request): InitialData
    {
        $complaintId = $request->get('complaintId');

        if (! is_int($complaintId) || $complaintId <= 0) {
            throw new \Exception('Complaint id required', 400);
        }

        /** @var NotificationMapper $notificationMapper */
        $notificationMapper = $this->mapperFactory->getMapper('Notification');
        $notification = $notificationMapper->getByName('complaint status changed');

        if (is_null($notification)) {
            throw new \Exception('Notification was not found!', 500);
        }

        /** @var ComplaintMapper $complaintMapper */
        $complaintMapper = $this->mapperFactory->getMapper('Complaint');
        $complaint = $complaintMapper->getById($complaintId);

        if (is_null($complaint)) {
            throw new \Exception('Complaint was not found!', 400);
        }

        $data = new NotifyComplaintStatusChangedInitialData();
        $data->setComplaint($complaint);
        $data->setNotification($notification);

        return $data;
    }
}

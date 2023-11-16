<?php

namespace NodaSoft\Operation\InitialData;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationNewMessageContentList;

class NotifyComplaintNewInitialData implements InitialData
{
    /** @var ReturnOperationNewMessageContentList */
    private $messageContentList;

    /** @var Reseller */
    private $reseller;

    /** @var Notification */
    private $notification;

    /** @var EmployeeCollection */
    private $employees;

    public function getMessageContentList(): ReturnOperationNewMessageContentList
    {
        return $this->messageContentList;
    }

    public function setMessageContentList(ReturnOperationNewMessageContentList $messageContentList): void
    {
        $this->messageContentList = $messageContentList;
    }

    public function getReseller(): Reseller
    {
        return $this->reseller;
    }

    public function setReseller(Reseller $reseller): void
    {
        $this->reseller = $reseller;
    }

    public function getNotification(): Notification
    {
        return $this->notification;
    }

    public function setNotification(Notification $notification): void
    {
        $this->notification = $notification;
    }

    public function getEmployees(): EmployeeCollection
    {
        return $this->employees;
    }

    public function setEmployees(EmployeeCollection $employees): void
    {
        $this->employees = $employees;
    }
}

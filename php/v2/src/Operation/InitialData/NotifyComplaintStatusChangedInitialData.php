<?php

namespace NodaSoft\Operation\InitialData;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationStatusChangedMessageBodyList;

class NotifyComplaintStatusChangedInitialData implements InitialData
{
    /** @var ReturnOperationStatusChangedMessageBodyList */
    private $messageTemplate;

    /** @var Reseller */
    private $reseller;

    /** @var Notification */
    private $notification;

    /** @var Client */
    private $client;

    /** @var EmployeeCollection */
    private $employees;

    public function getMessageTemplate(): ReturnOperationStatusChangedMessageBodyList
    {
        return $this->messageTemplate;
    }

    public function setMessageTemplate(
        ReturnOperationStatusChangedMessageBodyList $messageTemplate
    ): void {
        $this->messageTemplate = $messageTemplate;
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

    public function getClient(): Client
    {
        return $this->client;
    }

    public function setClient(Client $client): void
    {
        $this->client = $client;
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

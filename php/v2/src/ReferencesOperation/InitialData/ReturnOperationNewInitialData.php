<?php

namespace NodaSoft\ReferencesOperation\InitialData;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationNewMessageBodyList;

class ReturnOperationNewInitialData implements InitialData
{
    /** @var ReturnOperationNewMessageBodyList */
    private $messageTemplate;

    /** @var Reseller */
    private $reseller;

    /** @var Notification */
    private $notification;

    /** @var Client */
    private $client;

    /** @var Employee[] */
    private $employees;

    public function getMessageTemplate(): ReturnOperationNewMessageBodyList
    {
        return $this->messageTemplate;
    }

    public function setMessageTemplate(ReturnOperationNewMessageBodyList $messageTemplate): void
    {
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

    /**
     * @return Employee[]
     */
    public function getEmployees(): array
    {
        return $this->employees;
    }

    /**
     * @param Employee[] $employees
     * @return void
     */
    public function setEmployees(array $employees): void
    {
        $this->employees = $employees;
    }
}

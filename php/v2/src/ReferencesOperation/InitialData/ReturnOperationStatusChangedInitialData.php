<?php

namespace NodaSoft\ReferencesOperation\InitialData;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationStatusChangedMessageBodyList;

class ReturnOperationStatusChangedInitialData implements InitialData
{
    /** @var ReturnOperationStatusChangedMessageBodyList */
    private $messageTemplate;

    /** @var Reseller */
    private $reseller;

    /** @var Notification */
    private $notification;

    /** @var ?int */
    private $differencesFrom;

    /** @var ?int */
    private $differencesTo;

    /** @var Client */
    private $client;

    /** @var Employee[] */
    private $employees;

    public function getMessageTemplate(): ReturnOperationStatusChangedMessageBodyList
    {
        return $this->messageTemplate;
    }

    public function setMessageTemplate(ReturnOperationStatusChangedMessageBodyList $messageTemplate): void
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

    public function getDifferencesFrom(): ?int
    {
        return $this->differencesFrom;
    }

    public function setDifferencesFrom(?int $differencesFrom): void
    {
        $this->differencesFrom = $differencesFrom;
    }

    public function getDifferencesTo(): ?int
    {
        return $this->differencesTo;
    }

    public function setDifferencesTo(?int $differencesTo): void
    {
        $this->differencesTo = $differencesTo;
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

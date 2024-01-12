<?php

declare(strict_types=1);

namespace NW\WebService\References\Operations\Notification;

use NW\WebService\References\Operations\Notification\Exceptions\ClientNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\CreatorNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\EmptyNotificationTypeException;
use NW\WebService\References\Operations\Notification\Exceptions\EmptyResellerException;
use NW\WebService\References\Operations\Notification\Exceptions\ExpertNotFoundException;
use NW\WebService\References\Operations\Notification\Exceptions\SellerNotFoundException;

class ReturnOperationDTO
{
    public mixed $notificationType;
    public Seller $reseller;
    public Employee $creator;
    public Employee $expert;
    public Contractor $client;

    public int $newStatus;

    /**
     * @param  array  $data
     * @throws ClientNotFoundException
     * @throws CreatorNotFoundException
     * @throws EmptyNotificationTypeException
     * @throws EmptyResellerException
     * @throws ExpertNotFoundException
     * @throws SellerNotFoundException
     */
    public function __construct(array $data)
    {
        $resellerId = $this->getResellerId($data['resellerId']);
        $this->reseller = $this->getReseller($resellerId);
        $this->notificationType = $this->getNotificationType($data);
        $this->client = $this->getClient($data);
        $this->creator = $this->getCreator($data);
        $this->expert = $this->getExpert($data);
    }

    /**
     * @param  array  $data
     * @return mixed
     * @throws EmptyResellerException
     */
    private function getResellerId(array $data): mixed
    {

        return !empty((int) $data['resellerId']) ? (int) $data['resellerId'] : throw new EmptyResellerException('Empty resellerId',
            400);

    }

    /**
     * @param  int  $resellerId
     * @return Seller
     * @throws SellerNotFoundException
     */
    private function getReseller(int $resellerId): Seller
    {
        return Seller::getById($resellerId) ?? throw new SellerNotFoundException('Seller not found!', 400);
    }

    /**
     * @param  array  $data
     * @return int|mixed
     * @throws EmptyNotificationTypeException
     */
    private function getNotificationType(array $data): mixed
    {
        return !empty((int) $data['notificationType']) ? (int) $data['notificationType'] : throw new EmptyNotificationTypeException('Empty notificationType',
            400);
    }

    /**
     * @param  array  $data
     * @return Contractor
     * @throws ClientNotFoundException
     */
    private function getClient(array $data): Contractor
    {
        $client = Contractor::getById((int) $data['clientId']);

        if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->reseller->id) {
            throw new ClientNotFoundException('сlient not found!', 400);
        }

        return $client;
    }

    /**
     * @param  array  $data
     * @return Employee
     * @throws CreatorNotFoundException
     */
    private function getCreator(array $data): Employee
    {
        $client = Employee::getById((int) $data['creatorId']);

        if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->reseller->id) {
            throw new CreatorNotFoundException('Creator not found!', 400);
        }

        return $client;
    }

    /**
     * @param  array  $data
     * @return Employee
     * @throws ExpertNotFoundException
     */
    private function getExpert(array $data): Employee
    {
        $client = Employee::getById((int) $data['expertId']);

        if ($client->type !== Contractor::TYPE_CUSTOMER || $client->Seller->id !== $this->reseller->id) {
            throw new ExpertNotFoundException('Expert not found!', 400);
        }

        return $client;
    }
}

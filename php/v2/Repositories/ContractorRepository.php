<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace Repositories;

use NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException;
use NW\WebService\References\Operations\Notification\Exceptions\EntityNotExistException;
use NW\WebService\References\Operations\Notification\Helpers\Support;
use NW\WebService\References\Operations\Notification\Interfaces\ContractorInterface;
use Models\Contractor;

/**
 * Class ContractorRepository
 * Implements Repository Pattern for Entity
 */
class ContractorRepository
{
    /** @var \NW\WebService\References\Operations\Notification\Interfaces\ContractorInterface[] */
    private array $contractors = [];

    /**
     * @param int $contractorId
     * @param bool $forceReload
     * @return \Models\Contractor
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\EntityNotExistException
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    public function getById(int $contractorId, bool $forceReload = false): ContractorInterface
    {
        if (empty($contractorId)) {
            throw new InvalidArgumentsException(Support::__('Entity ID should not be empty'));
        }

        if ($forceReload || empty($this->contractors[$contractorId])) {
            $isFound = $contractorId < 1000; // mock DB request
            if (!$isFound) {
                throw new EntityNotExistException(Support::__('Entity not found :id', [':id' => $contractorId]));
            }

            /**
             * @todo fill out Contractor with real DB Data
             */
            $this->contractors[$contractorId] = new Contractor(
                id: $contractorId,
                name: 'Contractor Name',
                email: 'contractor@example.com',
                type: ContractorInterface::TYPE_CUSTOMER,
                mobile: '+79260000000'
            );
        }

        return $this->contractors[$contractorId];
    }
}
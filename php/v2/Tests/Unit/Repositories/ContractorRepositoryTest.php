<?php

/**
 * This file is part of the Notification package responsible for handling TS Goods Return operations
 *
 * @package  NW\WebService\References\Operations\Notification
 * @author   Dmitrii Fionov <dfionov@gmail.com>
 */

declare(strict_types=1);

namespace Tests\Unit\Repositories;

use NW\WebService\References\Operations\Notification\Exceptions\EntityNotExistException;
use NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException;
use PHPUnit\Framework\TestCase;
use Repositories\ContractorRepository;

/**
 * Class ContractorRepositoryTest
 * Test for @see ContractorRepository
 */
class ContractorRepositoryTest extends TestCase
{
    /** @var \Repositories\ContractorRepository */
    private ContractorRepository $contractorRepository;

    /**
     * @return void
     */
    protected function setUp(): void
    {
        parent::setUp();

        $this->contractorRepository = new ContractorRepository();
    }

    /**
     * @throws \NW\WebService\References\Operations\Notification\Exceptions\InvalidArgumentsException
     */
    public function testGetById(): void
    {
        $this->expectException(EntityNotExistException::class);
        $this->contractorRepository->getById(1001, true);

        $this->expectException(InvalidArgumentsException::class);
        $this->contractorRepository->getById(1001, true);
    }
}
<?php

use Nodasoft\Testapp\DTO\MessageDifferenceDto;
use Nodasoft\Testapp\DTO\SendNotificationDTO;
use Nodasoft\Testapp\Entities\Client\Client;
use Nodasoft\Testapp\Entities\Employee\Employee;
use Nodasoft\Testapp\Entities\Seller\Seller;
use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Enums\NotificationType;
use Nodasoft\Testapp\Repositories\Client\ClientRepositoryInterface;
use Nodasoft\Testapp\Repositories\Employee\EmployeeRepositoryInterface;
use Nodasoft\Testapp\Repositories\Seller\SellerRepositoryInterface;
use Nodasoft\Testapp\Services\SendNotification\Base\GetDifferencesInterface;
use Nodasoft\Testapp\Services\SendNotification\Base\ReferencesOperation;
use Nodasoft\Testapp\Services\SendNotification\TsReturnOperation;
use PHPUnit\Framework\MockObject\Exception;
use PHPUnit\Framework\TestCase;

class TsReturnOperationTest extends TestCase
{
    private ClientRepositoryInterface $clientRepository;
    private EmployeeRepositoryInterface $employeeRepository;
    private SellerRepositoryInterface $sellerRepository;
    private GetDifferencesInterface $differentService;
    private ReferencesOperation $operation;

    /**
     * @throws Exception
     */
    protected function setUp(): void
    {
        $this->clientRepository = $this->createMock(ClientRepositoryInterface::class);
        $this->employeeRepository = $this->createMock(EmployeeRepositoryInterface::class);
        $this->sellerRepository = $this->createMock(SellerRepositoryInterface::class);
        $this->differentService = $this->createMock(GetDifferencesInterface::class);

        $this->operation = new TsReturnOperation(
            $this->clientRepository,
            $this->employeeRepository,
            $this->sellerRepository,
            $this->differentService
        );
    }

    /**
     * @throws \Exception
     */
    public function testDoOperation(): void
    {
        $dto = new SendNotificationDTO(
            1,
            NotificationType::TYPE_NEW,
            1,
            1,
            1,
            1,
            '12345',
            1,
            '67890',
            '54321',
            '2023-01-01',
            new MessageDifferenceDto(1, 2)
        );

        $creator = new Employee(
            1,
            'some name',
            'some address',
            ContactorType::TYPE_EMPLOYEE
        );

        $expert = new Employee(
            1,
            'some name',
            'some address',
            ContactorType::TYPE_EMPLOYEE
        );

        $seller = new Seller(
            1,
            'some name',
            'some address',
            ContactorType::TYPE_SELLER
        );

        $client = new Client(
            1,
            'some name',
            1,
            'some email',
            'some mobile',
            ContactorType::TYPE_CUSTOMER
        );

        $this->sellerRepository->method('getById')->willReturn($seller);
        $this->clientRepository->method('getById')->willReturn($client);
        $this->employeeRepository->method('getById')->willReturnOnConsecutiveCalls($creator, $expert);
        $this->differentService->method('getDifference')->willReturn('Some differences');

        $result = $this->operation->doOperation($dto);

        $this->assertIsArray($result);
        $this->assertArrayHasKey('notificationEmployeeByEmail', $result);
        $this->assertArrayHasKey('notificationClientByEmail', $result);
        $this->assertArrayHasKey('notificationClientBySms', $result);
    }
}

<?php

namespace Tests\Unit\ReferencesOperation\Operation;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ClientMapper;
use NodaSoft\DataMapper\Mapper\EmployeeMapper;
use NodaSoft\DataMapper\Mapper\ResellerMapper;
use NodaSoft\ReferencesOperation\Factory\TsReturnOperationFactory;
use NodaSoft\ReferencesOperation\Operation\ReferencesOperation;
use NodaSoft\Request\HttpRequest;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;
use PHPUnit\Framework\TestCase;

/**
 * @group phpunit-excluded
 */
class ReferencesOperationTest extends TestCase
{
    protected function setUp(): void
    {
        require_once '/var/www/html/others.php';
    }

    /**
     * @dataProvider operationDataProvider
     */
    public function testDoOperation(
        array $data,
        bool $isEmployeeNotifiedEmail,
        bool $isClientNotifiedEmail,
        bool $isClientNotifiedSms,
        string $errorMessage
    ): void {
        $_REQUEST['data'] = $data;
        $request = new HttpRequest();
        $factory = new TsReturnOperationFactory();
        $mapperFactory = $this->getMapperFactoryMock();
        $tsReturnOperation = new ReferencesOperation(
            $factory,
            $request,
            $mapperFactory
        );
        /** @var TsReturnOperationResult $result */
        $result = $tsReturnOperation->doOperation();

        $this->assertSame($isEmployeeNotifiedEmail, $result->getEmployeeEmail()->isSent());
        $this->assertSame($isClientNotifiedEmail, $result->getClientEmail()->isSent());
        $this->assertSame($isClientNotifiedSms, $result->getClientSms()->isSent());
        $this->assertSame($errorMessage, $result->getClientSms()->getErrorMessage());
    }

    public function operationDataProvider(): \Generator
    {
        yield 'main' => [$this->getMainData(), true, true, true, ''];
    }

    private function getMainData(): array
    {
        return [
            'resellerId' => 86, //int
            'notificationType' => 2, //int
            'clientId' => 27, //int
            'creatorId' => 12, //int
            'expertId' => 7, //int
            'differences' => [
                'from' => 1, //int
                'to' => 1, //int
            ],
            'complaintId' => 1, //int
            'complaintNumber' => 1, //string
            'consumptionId' => 1, //int
            'consumptionNumber' => 1, //string
            'agreementNumber' => 1, //string
            'date' => 1, //string
        ];
    }

    private function getMapperFactoryMock(): MapperFactory
    {
        $mapperFactory = $this->createMock(MapperFactory::class);
        $employeeMapper = $this->createMock(EmployeeMapper::class);
        $clientMapper = $this->createMock(ClientMapper::class);
        $resellerMapper = $this->createMock(ResellerMapper::class);

        $expert = $this->createMock(Employee::class);
        $creator = $this->createMock(Employee::class);
        $client = $this->createMock(Client::class);
        $reseller = $this->createMock(Reseller::class);

        $expert->setId(7);
        $expert->method('getFullName')->willReturn("Boris 7");
        $creator->setId(12);
        $creator->method('getFullName')->willReturn("Sarah 12");
        $client->setId(27);
        $client->method('getFullName')->willReturn("Anna 27");
        $client->method('getEmail')->willReturn("Anna.27@gmail.com");
        $client->method('getCellphoneNumber')->willReturn(555989898);
        $reseller->setId(86);


        $employeeMapper
            ->method('getById')
            ->willReturnMap([
                [12, $creator],
                [7, $expert],
            ]);

        $clientMapper
            ->method('getById')
            ->with(27)
            ->willReturn($client);

        $resellerMapper
            ->method('getById')
            ->with(86)
            ->willReturn($reseller);

        $mapperFactory
            ->method('getMapper')
            ->willReturnMap([
                ['Reseller', $resellerMapper],
                ['Client', $clientMapper],
                ['Employee', $employeeMapper],
            ]);

        return $mapperFactory;
    }
}

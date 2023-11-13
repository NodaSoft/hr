<?php

namespace Tests\Unit\ReferencesOperation\Operation;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ClientMapper;
use NodaSoft\DataMapper\Mapper\EmployeeMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
use NodaSoft\DataMapper\Mapper\ResellerMapper;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Message\Client\EmailClient;
use NodaSoft\Message\Client\SmsClient;
use NodaSoft\Message\Message;
use NodaSoft\Message\Messenger;
use NodaSoft\ReferencesOperation\Factory\ReturnOperationStatusChangedFactory;
use NodaSoft\ReferencesOperation\Operation\ReferencesOperation;
use NodaSoft\ReferencesOperation\Result\ReturnOperationStatusChangedResult;
use NodaSoft\Request\HttpRequest;
use PHPUnit\Framework\TestCase;

/**
 * @group phpunit-excluded
 */
class ReferencesOperationTest extends TestCase
{
    const COMPLAINT_STATUS = ['new' => 1, 'changed' => 2];

    /**
     * @dataProvider operationDataProvider
     */
    public function testDoOperation(
        array $data,
        bool $isEmployee1NotifiedEmail,
        bool $isEmployee2NotifiedEmail,
        bool $isClientNotifiedEmail,
        bool $isClientNotifiedSms,
        string $errorMessage
    ): void {
        $_REQUEST['data'] = $data;
        $dependencies = $this->mockDependencies();
        $request = new HttpRequest();
        $factory = new ReturnOperationStatusChangedFactory();
        $mapperFactory = $this->getMapperFactoryMock();
        $tsReturnOperation = new ReferencesOperation(
            $dependencies,
            $factory,
            $request,
            $mapperFactory
        );
        /** @var ReturnOperationStatusChangedResult $result */
        $result = $tsReturnOperation->doOperation();
        $employeeEmails = $result->getEmployeeEmails()->getList();

        $this->assertSame($isEmployee1NotifiedEmail, $employeeEmails[0]->isSent(), "employee1's email");
        $this->assertSame($isEmployee2NotifiedEmail, $employeeEmails[1]->isSent(), "employee2's email");
        $this->assertSame($isClientNotifiedEmail, $result->getClientEmail()->isSent(), "client's email");
        $this->assertSame($isClientNotifiedSms, $result->getClientSms()->isSent(), "client's sms");
        $this->assertSame($errorMessage, $result->getClientSms()->getErrorMessage(), "client's sms error message");
    }

    public function operationDataProvider(): \Generator
    {
        yield 'main' => [$this->getMainData(), true, true, true, true, ''];
    }

    private function getMainData(): array
    {
        return [
            'resellerId' => 86, //int
            'notificationType' => self::COMPLAINT_STATUS['changed'], //int
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
        $notificationMapper = $this->createMock(NotificationMapper::class);

        $expert = new Employee();
        $creator = new Employee();
        $employee1 = new Employee();
        $employee2 = new Employee();
        $client = new Client();
        $reseller = new Reseller();
        $notificationNew = new Notification();
        $notificationChanged = new Notification();

        $expert->setId(7);
        $expert->setName("Boris");
        $expert->setEmail("Boris@mail.com");

        $creator->setId(12);
        $creator->setName('Sarah');

        $employee1->setId(64);
        $employee1->setName("Mark");
        $employee1->setEmail("mark@mail.ru");

        $employee2->setId(65);
        $employee2->setName("Nana");
        $employee2->setEmail("nana@mail.ru");

        $reseller->setId(86);
        $reseller->setEmail('reseller@mail.ru');

        $client->setId(27);
        $client->setName('Anna');
        $client->setEmail("Anna.27@gmail.com");
        $client->setCellphone(5559898989);
        $client->setIsCustomer(true);
        $client->setReseller($reseller);

        $notificationNew->setId(self::COMPLAINT_STATUS['new']);
        $notificationNew->setName('complaint status new');
        $notificationNew->setTemplate("Added new entry (reseller id: #resellerId#).");
        $notificationChanged->setId(self::COMPLAINT_STATUS['changed']);
        $notificationChanged->setName('complaint status changed');
        $notificationChanged->setTemplate("Entry status changed (resseler id: #resellerId#): previous status: #differencesFrom#, current status: #differencesTo#");

        $employeeMapper
            ->method('getById')
            ->willReturnMap([
                [12, $creator],
                [7, $expert],
            ]);

        $employeeMapper
            ->method('getAllByReseller')
            ->with(86)
            ->willReturn([$employee1, $employee2]);

        $clientMapper
            ->method('getById')
            ->with(27)
            ->willReturn($client);

        $resellerMapper
            ->method('getById')
            ->with(86)
            ->willReturn($reseller);

        $notificationMapper
            ->method('getById')
            ->willReturnMap([
                [1, $notificationNew],
                [2, $notificationChanged],
            ]);

        $mapperFactory
            ->method('getMapper')
            ->willReturnMap([
                ['Reseller', $resellerMapper],
                ['Client', $clientMapper],
                ['Employee', $employeeMapper],
                ['Notification', $notificationMapper],
            ]);

        return $mapperFactory;
    }

    private function mockDependencies(): Dependencies
    {
        $mail = $this->mockMail();
        $sms = $this->mockSms();
        $dependency = $this->createMock(Dependencies::class);
        $dependency->method('getMailService')->willReturn($mail);
        $dependency->method('getSmsService')->willReturn($sms);
        return $dependency;
    }

    public function mockMail(): Messenger
    {
        $emailClient = new class extends EmailClient {
            //mock the class: always return true (isSuccess) for every try to send email
            public function send(Message $message): bool
            {
                return true;
            }
        };
        return new Messenger($emailClient);
    }

    private function mockSms(): Messenger
    {
        $smsClient = new class extends SmsClient {
            //mock the class: always return true (isSuccess) for every try to send email
            public function send(Message $message): bool
            {
                return true;
            }
        };
        return new Messenger($smsClient);
    }
}

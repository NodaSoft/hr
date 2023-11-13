<?php

namespace Tests\Unit\ReferencesOperation\Operation;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ComplaintMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Message\Client\EmailClient;
use NodaSoft\Message\Client\SmsClient;
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
        $mapperFactory = $this->mockMapperFactory();
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

    private function mockMapperFactory(): MapperFactory
    {
        $expert = new Employee(7, 'Boris', 'boris@mail.com');
        $creator = new Employee(12, 'Sarah', 'sarah@mail.com');
        $employees = new EmployeeCollection([
            new Employee(64, 'Mark', 'mark@mail.ru'),
            new Employee(65, 'Nana', 'nana@mail.ru'),
        ]);
        $reseller = new Reseller(86, 'Bob', 'reseller@mail.ru', 1234567890, $employees);
        $client = new Client(27, 'Anna', 'anna.27@gmail.com', 5559898989, true, $reseller);
        $notificationNew = new Notification(
            self::COMPLAINT_STATUS['new'],
            'complaint new',
            'Added new entry (reseller id: #resellerId#).'
        );
        $notificationChanged = new Notification(
            self::COMPLAINT_STATUS['changed'],
            'complaint status changed',
            'Status changed (#complaintId#): previous status: #differencesFrom#, current status: #differencesTo#'
        );
        $complaint = new Complaint(1, 'Test Complaint', $creator, $client, $expert, $reseller);


        $mapperFactory = $this->createMock(MapperFactory::class);
        $notificationMapper = $this->createMock(NotificationMapper::class);
        $complaintMapper = $this->createMock(ComplaintMapper::class);

        $complaintMapper
            ->method('getById')
            ->willReturnMap([
                [1, $complaint],
            ]);

        $notificationMapper
            ->method('getByName')
            ->willReturnMap([
                ['complaint new', $notificationNew],
                ['complaint status changed', $notificationChanged],
            ]);

        $mapperFactory
            ->method('getMapper')
            ->willReturnMap([
                ['Complaint', $complaintMapper],
                ['Notification', $notificationMapper],
            ]);

        return $mapperFactory;
    }

    private function mockDependencies(): Dependencies
    {
        $emailClient = $this->createMock(EmailClient::class);
        $emailClient->method('send')->withAnyParameters()->willReturn(true);
        $emailClient->method('isValid')->withAnyParameters()->willReturn(true);
        $mailMessenger = new Messenger($emailClient);

        $smsClient = $this->createMock(SmsClient::class);
        $smsClient->method('send')->withAnyParameters()->willReturn(true);
        $smsClient->method('isValid')->withAnyParameters()->willReturn(true);
        $smsMessenger = new Messenger($smsClient);

        $dependency = $this->createMock(Dependencies::class);
        $dependency->method('getMailService')->willReturn($mailMessenger);
        $dependency->method('getSmsService')->willReturn($smsMessenger);

        return $dependency;
    }
}

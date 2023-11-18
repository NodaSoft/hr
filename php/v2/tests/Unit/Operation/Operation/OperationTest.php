<?php

namespace Tests\Unit\Operation;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\DataMapper\Entity\ComplaintStatus;
use NodaSoft\DataMapper\Entity\Consumption;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\DataMapper\Mapper\ComplaintMapper;
use NodaSoft\DataMapper\Mapper\NotificationMapper;
use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Client\SmsClient;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\Operation;
use NodaSoft\Operation\Result\NotifyComplaintStatusChangedResult;
use NodaSoft\Request\HttpRequest;
use PHPUnit\Framework\TestCase;

/**
 * @group phpunit-excluded
 */
class OperationTest extends TestCase
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
        $_SERVER['REQUEST_URI'] = "/notify/complaint/status_changed";
        $dependencies = $this->mockDependencies();
        $operation = new Operation($dependencies);
        /** @var NotifyComplaintStatusChangedResult $result */
        $result = $operation->doOperation();
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
            'complaintId' => 1,
        ];
    }

    private function mockMapperFactory(): MapperFactory
    {
        $employees = $this->getEmployees();
        $creator = $employees['creator'];
        $expert = $employees['expert'];
        $reseller = new Reseller(
            33,
            'Dora',
            'dora@mail.ru',
            1234567890,
            new EmployeeCollection($employees)
        );
        $consumption = new Consumption(
            1,
            'foo client\'s consumption',
            'p12',
            'm17'

        );
        $client = new Client(
            11,
            'Anna',
            'anna@mail.ru',
            1234567890,
            true,
            $reseller,
            $consumption
        );

        $complaint = new Complaint(
            11,
            "Foo complaint",
            $creator,
            $client,
            $expert,
            $reseller,
            new ComplaintStatus(5, 'closed'),
            new ComplaintStatus(6, 'reopened'),
            'AO16578-g'
        );

        $notificationNew = new Notification(
            self::COMPLAINT_STATUS['new'],
            'complaint new',
            'Added new entry (reseller id: #resellerId#).',
            'Added new entry (reseller id: #resellerId#).'
        );

        $notificationChanged = new Notification(
            self::COMPLAINT_STATUS['changed'],
            'complaint status changed',
            'Status changed (#complaintId#): previous status: #previousStatusName#, status: #currentStatusName#',
            'Status changed (#complaintId#): previous status: #previousStatusName#, status: #currentStatusName#'
        );

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

    /**
     * @return array{creator: Employee, expert: Employee}
     */
    private function getEmployees(): array
    {
        $creator = new Employee(
            22,
            'Sarah',
            'sarah@mail.ru',
            1234567890
        );
        $expert = new Employee(
            21,
            'Bob',
            'bob@mail.ru',
            1234567890
        );
        return ['creator' => $creator, 'expert' => $expert];
    }

    private function mockDependencies(): Dependencies
    {
        $mapperFactory = $this->mockMapperFactory();

        $emailClient = $this->createMock(EmailClient::class);
        $emailClient->method('send')->withAnyParameters()->willReturn(true);
        $emailClient->method('isValid')->withAnyParameters()->willReturn(true);
        $emailService = new Messenger($emailClient);

        $smsClient = $this->createMock(SmsClient::class);
        $smsClient->method('send')->withAnyParameters()->willReturn(true);
        $smsClient->method('isValid')->withAnyParameters()->willReturn(true);
        $smsService = new Messenger($smsClient);

        $dependency = $this->createMock(Dependencies::class);
        $dependency->method('getEmailService')->willReturn($emailService);
        $dependency->method('getSmsService')->willReturn($smsService);
        $dependency->method('getMapperFactory')->willReturn($mapperFactory);

        $request = new HttpRequest();

        return new Dependencies($request, $emailService, $smsService, $mapperFactory);
    }
}

<?php

namespace Tests\Unit\Operation\Command;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Complaint;
use NodaSoft\DataMapper\Entity\ComplaintStatus;
use NodaSoft\DataMapper\Entity\Consumption;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\Command\NotifyComplaintNewCommand;
use NodaSoft\Operation\InitialData\NotifyComplaintNewInitialData;
use PHPUnit\Framework\TestCase;

class NotifyComplaintNewCommandTest extends TestCase
{
    public function testExecute(): void
    {
        $emailClient = $this->createMock(EmailClient::class);
        $emailClient->method('send')->withAnyParameters()->willReturn(true);
        $emailClient->method('isValid')->withAnyParameters()->willReturn(true);
        $command = new NotifyComplaintNewCommand();
        $command->setEmail(new Messenger($emailClient));
        $result = $command->execute($this->mockInitialData());
        $employees = $this->getEmployees();
        $this->assertSame($result->toArray(), [
            'employeeEmails' => [
                [
                    'isSent' => true,
                    'clientClass' => get_class($emailClient),
                    'errorMessage' => '',
                    'recipient' => $employees['creator']->toArray(),
                ],
                [
                    'isSent' => true,
                    'clientClass' => get_class($emailClient),
                    'errorMessage' => '',
                    'recipient' => $employees['expert']->toArray(),
                ],
            ]
        ]);
    }

    private function mockInitialData(): NotifyComplaintNewInitialData
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

        $notification = new Notification(
            '21',
            'complaint new',
            'reseller: #resellerId#, client: #clientId#',
            'reseller: #resellerId#, client: #clientId#'
        );

        $data = new NotifyComplaintNewInitialData();
        $data->setComplaint($complaint);
        $data->setNotification($notification);

        return $data;
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
}

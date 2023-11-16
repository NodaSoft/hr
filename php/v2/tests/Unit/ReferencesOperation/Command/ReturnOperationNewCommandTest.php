<?php

namespace Tests\Unit\ReferencesOperation\Command;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationNewMessageContentList;
use NodaSoft\GenericDto\Factory\GenericDtoFactory;
use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\Command\NotifyComplaintNewCommand;
use NodaSoft\Operation\InitialData\NotifyComplaintNewInitialData;
use NodaSoft\Operation\Result\NotifyComplaintNewResult;
use PHPUnit\Framework\TestCase;

class ReturnOperationNewCommandTest extends TestCase
{
    public function testExecute(): void
    {
        $emailClient = $this->createMock(EmailClient::class);
        $emailClient->method('send')->withAnyParameters()->willReturn(true);
        $emailClient->method('isValid')->withAnyParameters()->willReturn(true);
        $command = new NotifyComplaintNewCommand();
        $command->setMail(new Messenger($emailClient));
        $command->setInitialData($this->mockInitialData());
        $result = $command->execute();
        $this->assertSame($result->toArray(), [
            'employeeEmails' => [
                [
                    'isSent' => true,
                    'clientClass' => get_class($emailClient),
                    'errorMessage' => '',
                    'recipient' => [
                        'id' => 21,
                        'name' => 'Bob',
                        'email' => 'bob@mail.ru',
                        'cellphone' => 9876543210
                    ],
                ],
                [
                    'isSent' => true,
                    'clientClass' => get_class($emailClient),
                    'errorMessage' => '',
                    'recipient' => [
                        'id' => 23,
                        'name' => 'Mark',
                        'email' => 'mark@mailru',
                        'cellphone' => 1111111111
                    ],
                ],
            ]
        ]);
    }

    private function mockInitialData(): NotifyComplaintNewInitialData
    {
        $dtoFactory = new GenericDtoFactory();
        $list = $dtoFactory->fillDtoArray(
            new ReturnOperationNewMessageContentList(),
            [
                'complaintId' => 4343421,
                'complaintNumber' => '06.07.2008FV',
                'creatorId' => 27,
                'creatorName' => 'Alen',
                'expertId' => 21,
                'expertName' => 'Bob',
                'resellerId' => 31,
                'clientId' => 11,
                'clientName' => 'Anna',
                'consumptionId' => 2,
                'consumptionNumber' => 'AFG83',
                'agreementNumber' => 'PO67',
                'date' => '11.12.2023'
            ]
        );
        $data = new NotifyComplaintNewInitialData();
        $reseller = new Reseller(31, 'John', 'john@mail.ru', 1234567890);
        $data->setReseller($reseller);
        $data->setEmployees(new EmployeeCollection([
            new Employee(21, 'Bob', 'bob@mail.ru', 9876543210),
            new Employee(23, 'Mark', 'mark@mailru', 1111111111),
        ]));
        $data->setNotification(new Notification(
            1,
            'new',
            'reseller: #resellerId#, client: #clientId#, date: #date#',
            'reseller: #resellerId#, client: #clientId#, date: #date#')
        );
        $data->setMessageContentList($list);
        return $data;
    }
}

<?php

namespace Tests\Unit\ReferencesOperation\Command;

use NodaSoft\DataMapper\Collection\EmployeeCollection;
use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\DataMapper\Entity\Reseller;
use NodaSoft\GenericDto\Dto\ReturnOperationStatusChangedMessageBodyList;
use NodaSoft\GenericDto\Factory\GenericDtoFactory;
use NodaSoft\Messenger\Client\EmailClient;
use NodaSoft\Messenger\Client\SmsClient;
use NodaSoft\Messenger\Messenger;
use NodaSoft\Operation\Command\NotifyComplaintStatusChangedCommand;
use NodaSoft\Operation\InitialData\NotifyComplaintStatusChangedInitialData;
use NodaSoft\Operation\Result\ReturnOperationStatusChangedResult;
use PHPUnit\Framework\TestCase;

class ReturnOperationStatusChangedCommandTest extends TestCase
{
    public function testExecute(): void
    {
        $emailClient = $this->createMock(EmailClient::class);
        $emailClient->method('send')->withAnyParameters()->willReturn(true);
        $emailClient->method('isValid')->withAnyParameters()->willReturn(true);
        $smsClient = $this->createMock(SmsClient::class);
        $smsClient->method('send')->withAnyParameters()->willReturn(true);
        $smsClient->method('isValid')->withAnyParameters()->willReturn(true);
        $command = new NotifyComplaintStatusChangedCommand();
        $command->setMail(new Messenger($emailClient));
        $command->setSms(new Messenger($smsClient));
        $command->setResult(new ReturnOperationStatusChangedResult());
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
            ],
            'clientEmail' => [
                'isSent' => true,
                'clientClass' => get_class($emailClient),
                'errorMessage' => '',
                'recipient' => [
                    'id' => 11,
                    'name' => 'Anna',
                    'email' => 'anna@mail.ru',
                    'cellphone' => 2222222222,
                ],
            ],
            'clientSms' => [
                    'isSent' => true,
                    'clientClass' => get_class($smsClient),
                    'errorMessage' => '',
                    'recipient' => [
                        'id' => 11,
                        'name' => 'Anna',
                        'email' => 'anna@mail.ru',
                        'cellphone' => 2222222222,
                    ],
            ]
        ]);
    }

    private function mockInitialData(): NotifyComplaintStatusChangedInitialData
    {
        $dtoFactory = new GenericDtoFactory();
        $list = $dtoFactory->fillDtoArray(
            new ReturnOperationStatusChangedMessageBodyList(),
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
        $data = new NotifyComplaintStatusChangedInitialData();
        $reseller = new Reseller(31, 'John', 'john@mail.ru', 1234567890);
        $data->setReseller($reseller);
        $data->setEmployees(new EmployeeCollection([
            new Employee(21, 'Bob', 'bob@mail.ru', 9876543210),
            new Employee(23, 'Mark', 'mark@mailru', 1111111111),
        ]));
        $data->setClient(new Client(11, 'Anna', 'anna@mail.ru', 2222222222, true, $reseller));
        $data->setNotification(new Notification(
            1,
            'new',
            'reseller: #resellerId#, client: #clientId#, date: #date#',
            'reseller: #resellerId#, client: #clientId#, date: #date#'
        ));
        $data->setMessageTemplate($list);
        return $data;
    }
}

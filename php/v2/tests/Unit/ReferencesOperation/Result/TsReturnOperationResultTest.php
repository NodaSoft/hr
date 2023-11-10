<?php

namespace Tests\Unit\ReferencesOperation\Result;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\Message\Client\EmailClient;
use NodaSoft\Message\Client\SmsClient;
use NodaSoft\Message\Result;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;
use PHPUnit\Framework\TestCase;

class TsReturnOperationResultTest extends TestCase
{
    public function testToArray(): void
    {
        $origin = [
            'employeeEmails' => [
                [
                    'isSent' => true,
                    'clientClass' => EmailClient::class,
                    'errorMessage' => '',
                    'recipient' => [
                        'id' => 7,
                        'name' => "Bob",
                        'email' => "bob@mail.com",
                        'cellphone' => null,
                    ],
                ],
                [
                    'isSent' => false,
                    'clientClass' => EmailClient::class,
                    'errorMessage' => 'mail(): "sendmail_from" not set in php.ini or custom "From:" header missing.',
                    'recipient' => [
                        'id' => 12,
                        'name' => "Sarah",
                        'email' => "sarah@mail.com",
                        'cellphone' => null,
                    ],
                ],
            ],
            'clientEmail' => [
                'isSent' => true,
                'clientClass' => EmailClient::class,
                'errorMessage' => '',
                'recipient' => [
                    'id' => 345,
                    'name' => "Anna",
                    'email' => "anna@mail.com",
                    'cellphone' => 1234567890,
                ],
            ],
            'clientSms' => [
                'isSent' => true,
                'clientClass' => SmsClient::class,
                'errorMessage' => 'Foo Bar Baz',
                'recipient' => [
                    'id' => 345,
                    'name' => 'Anna',
                    'email' => 'anna@mail.com',
                    'cellphone' => 1234567890,
                ]

            ],
        ];

        $bob = new Employee();
        $bob->setId($origin['employeeEmails'][0]['recipient']['id']);
        $bob->setName($origin['employeeEmails'][0]['recipient']['name']);
        $bob->setEmail($origin['employeeEmails'][0]['recipient']['email']);
        $sarah = new Employee();
        $sarah->setId($origin['employeeEmails'][1]['recipient']['id']);
        $sarah->setName($origin['employeeEmails'][1]['recipient']['name']);
        $sarah->setEmail($origin['employeeEmails'][1]['recipient']['email']);
        $client = new Client();
        $client->setId($origin['clientEmail']['recipient']['id']);
        $client->setName($origin['clientEmail']['recipient']['name']);
        $client->setEmail($origin['clientEmail']['recipient']['email']);
        $client->setCellphone(1234567890);
        $bobEmailResult = new Result(
            $bob,
            EmailClient::class,
            $origin['employeeEmails'][0]['isSent'],
            $origin['employeeEmails'][0]['errorMessage']
        );
        $sarahEmailResult = new Result(
            $sarah,
            EmailClient::class,
            $origin['employeeEmails'][1]['isSent'],
            $origin['employeeEmails'][1]['errorMessage']
        );
        $clientEmailResult = new Result(
            $client,
            EmailClient::class,
            $origin['clientEmail']['isSent'],
            $origin['clientEmail']['errorMessage']
        );
        $clientSmsResult = new Result(
            $client,
            SmsClient::class,
            $origin['clientSms']['isSent'],
            $origin['clientSms']['errorMessage']
        );
        $result = new TsReturnOperationResult();
        $result->addEmployeeEmailResult($bobEmailResult);
        $result->addEmployeeEmailResult($sarahEmailResult);
        $result->setClientEmailResult($clientEmailResult);
        $result->setClientSmsResult($clientSmsResult);
        $resultArray = $result->toArray();
        $this->assertSame($origin, $resultArray);
    }
}

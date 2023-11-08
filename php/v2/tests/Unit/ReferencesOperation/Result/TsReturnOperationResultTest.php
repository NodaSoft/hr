<?php

namespace Tests\Unit\ReferencesOperation\Result;

use NodaSoft\DataMapper\Entity\Client;
use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\Mail\Result;
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
                    'errorMessage' => '',
                    'recipient' => [
                        'id' => 7,
                        'name' => "Bob",
                        'email' => "bob@mail.com",
                    ],
                ],
                [
                    'isSent' => false,
                    'errorMessage' => 'mail(): "sendmail_from" not set in php.ini or custom "From:" header missing.',
                    'recipient' => [
                        'id' => 12,
                        'name' => "Sarah",
                        'email' => "sarah@mail.com",
                    ],
                ],
            ],
            'clientEmail' => [
                'isSent' => true,
                'errorMessage' => '',
                'recipient' => [
                    'id' => 345,
                    'name' => "Anna",
                    'email' => "anna@mail.com",
                ],
            ],
            'clientSms' => [
                'isSent' => true,
                'errorMessage' => 'Foo Bar Baz',
                'recipient' => null //todo: add recipient
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
        $bobEmailResult = new Result(
            $bob,
            $origin['employeeEmails'][0]['isSent'],
            $origin['employeeEmails'][0]['errorMessage']
        );
        $sarahEmailResult = new Result(
            $sarah,
            $origin['employeeEmails'][1]['isSent'],
            $origin['employeeEmails'][1]['errorMessage']
        );
        $clientEmailResult = new Result(
            $client,
            $origin['clientEmail']['isSent'],
            $origin['clientEmail']['errorMessage']
        );
        $result = new TsReturnOperationResult();
        $result->addEmployeeEmailResult($bobEmailResult);
        $result->addEmployeeEmailResult($sarahEmailResult);
        $result->setClientEmailResult($clientEmailResult);
        $result->markClientSmsSent();
        $result->setClientSmsErrorMessage($origin['clientSms']['errorMessage']);
        $resultArray = $result->toArray();
        $this->assertSame($origin, $resultArray);
    }
}

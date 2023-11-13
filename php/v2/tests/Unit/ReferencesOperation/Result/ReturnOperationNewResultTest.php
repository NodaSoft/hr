<?php

namespace Tests\Unit\ReferencesOperation\Result;

use NodaSoft\DataMapper\Entity\Employee;
use NodaSoft\Message\Client\EmailClient;
use NodaSoft\Message\Result;
use NodaSoft\ReferencesOperation\Result\ReturnOperationNewResult;
use PHPUnit\Framework\TestCase;

class ReturnOperationNewResultTest extends TestCase
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
        ];

        $bob = new Employee();
        $bob->setId($origin['employeeEmails'][0]['recipient']['id']);
        $bob->setName($origin['employeeEmails'][0]['recipient']['name']);
        $bob->setEmail($origin['employeeEmails'][0]['recipient']['email']);
        $sarah = new Employee();
        $sarah->setId($origin['employeeEmails'][1]['recipient']['id']);
        $sarah->setName($origin['employeeEmails'][1]['recipient']['name']);
        $sarah->setEmail($origin['employeeEmails'][1]['recipient']['email']);
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
        $result = new ReturnOperationNewResult();
        $result->addEmployeeEmailResult($bobEmailResult);
        $result->addEmployeeEmailResult($sarahEmailResult);
        $resultArray = $result->toArray();
        $this->assertSame($origin, $resultArray);
    }
}

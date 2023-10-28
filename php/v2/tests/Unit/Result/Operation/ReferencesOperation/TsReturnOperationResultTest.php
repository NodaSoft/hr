<?php

namespace Tests\Unit\Result\Operation\ReferencesOperation;

use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;
use PHPUnit\Framework\TestCase;

class TsReturnOperationResultTest extends TestCase
{
    public function testToArray(): void
    {
        $origin = [
            'employeeEmail' => true,
            'clientEmail' => false,
            'clientSms' => [
                'isSent' => true,
                'errorMessage' => 'Foo Bar Baz',
            ],
        ];
        $result = new TsReturnOperationResult();
        $result->markEmployeeEmailSent();
        $result->markClientSmsSent();
        $result->setClientSmsErrorMessage($origin['clientSms']['errorMessage']);
        $resultArray = $result->toArray();
        $this->assertSame($origin, $resultArray);
    }
}

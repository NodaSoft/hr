<?php

namespace Tests\Unit\ReferencesOperation\Operation;

use NodaSoft\ReferencesOperation\Factory\TsReturnOperationFactory;
use NodaSoft\ReferencesOperation\Operation\ReferencesOperation;
use NodaSoft\Request\HttpRequest;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;
use PHPUnit\Framework\TestCase;

/**
 * @group phpunit-excluded
 */
class ReferencesOperationTest extends TestCase
{
    /**
     * @dataProvider operationDataProvider
     */
    public function testDoOperation(
        array $data,
        bool $isEmployeeNotifiedEmail,
        bool $isClientNotifiedEmail,
        bool $isClientNotifiedSms,
        bool $errorMessage
    ): void {
        $_REQUEST['data'] = $data;
        $request = new HttpRequest();
        $factory = new TsReturnOperationFactory();
        $tsReturnOperation = new ReferencesOperation($factory, $request);
        /** @var TsReturnOperationResult $result */
        $result = $tsReturnOperation->doOperation();

        $this->assertSame($isEmployeeNotifiedEmail, $result->getEmployeeEmail()->isSent());
        $this->assertSame($isClientNotifiedEmail, $result->getClientEmail()->isSent());
        $this->assertSame($isClientNotifiedSms, $result->getClientSms()->isSent());
        $this->assertSame($errorMessage, $result->getClientSms()->getErrorMessage());
    }

    public function operationDataProvider(): \Generator
    {
        yield 'main' => [$this->getMainData(), false, false, false, ''];
    }

    private function getMainData(): array
    {
        return [
            'resellerId' => 1, //int
            'notificationType' => 2, //int
            'clientId' => 1, //int
            'creatorId' => 1, //int
            'expertId' => 1, //int
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
}

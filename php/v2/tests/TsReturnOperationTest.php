<?php

use NW\WebService\References\Operations\Notification\Contracts\MessagesClientInterface;
use NW\WebService\References\Operations\Notification\Contracts\NotificationManagerInterface;
use NW\WebService\References\Operations\Notification\TsReturnOperation;
use PHPUnit\Framework\TestCase;

class TsReturnOperationTest extends TestCase
{
    private $messagesClientMock;
    private $notificationManagerMock;

    protected function setUp(): void
    {
        $this->messagesClientMock = $this->createMock(MessagesClientInterface::class);
        $this->notificationManagerMock = $this->createMock(NotificationManagerInterface::class);
    }

    public function testDoOperationSuccess()
    {
        $_REQUEST['data'] = [
            'resellerId' => 1,
            'notificationType' => TsReturnOperation::TYPE_NEW,
            'clientId' => 1,
            'creatorId' => 1,
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->messagesClientMock->expects($this->once())
            ->method('sendMessage');

        $result = $operation->doOperation();

        $this->assertTrue($result['notificationEmployeeByEmail']);
    }

    public function testDoOperationMissingResellerId()
    {
        $_REQUEST['data'] = [
            'notificationType' => TsReturnOperation::TYPE_NEW,
            'clientId' => 1,
            'creatorId' => 1,
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('Empty resellerId');

        $operation->doOperation();
    }

    public function testDoOperationMissingClient()
    {
        $_REQUEST['data'] = [
            'resellerId' => 1,
            'notificationType' => TsReturnOperation::TYPE_NEW,
            'clientId' => 999, // неверный ID клиента
            'creatorId' => 1,
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('Client not found');

        $operation->doOperation();
    }

    public function testDoOperationMissingCreator()
    {
        $_REQUEST['data'] = [
            'resellerId' => 1,
            'notificationType' => TsReturnOperation::TYPE_NEW,
            'clientId' => 1,
            'creatorId' => null, // отсутствует ID создателя
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('Creator not found');

        $operation->doOperation();
    }

    public function testDoOperationInvalidNotificationType()
    {
        $_REQUEST['data'] = [
            'resellerId' => 1,
            'notificationType' => 999, // неверный тип уведомления
            'clientId' => 1,
            'creatorId' => 1,
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('Invalid notificationType');

        $operation->doOperation();
    }

    public function testDoOperationMissingNotificationType()
    {
        $_REQUEST['data'] = [
            'resellerId' => 1,
            'clientId' => 1,
            'creatorId' => 1,
            'expertId' => 1,
            'complaintId' => 1,
            'complaintNumber' => 'C123',
            'consumptionId' => 1,
            'consumptionNumber' => 'CN123',
            'agreementNumber' => 'AN123',
            'date' => '2023-06-25',
        ];

        $operation = new TsReturnOperation($this->messagesClientMock, $this->notificationManagerMock);

        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('Empty notificationType');

        $operation->doOperation();
    }
}


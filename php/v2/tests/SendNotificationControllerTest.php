<?php

namespace Tests;

use Nodasoft\Testapp\DTO\SendNotificationDTO;
use Nodasoft\Testapp\Http\Controllers\SendNotificationController;
use Nodasoft\Testapp\Services\SendNotification\Base\ReferencesOperation;
use Nyholm\Psr7\Factory\Psr17Factory;
use PHPUnit\Framework\MockObject\Exception;
use PHPUnit\Framework\TestCase;
use Psr\Http\Message\ResponseInterface;

class SendNotificationControllerTest extends TestCase
{
    private ReferencesOperation $referencesOperation;
    private SendNotificationController $controller;

    /**
     * @throws Exception
     */
    protected function setUp(): void
    {
        $this->referencesOperation = $this->createMock(ReferencesOperation::class);
        $this->controller = new SendNotificationController($this->referencesOperation);
    }

    public function test_controller_can_handle_request_to_send_notification(): void
    {
        $request = (new Psr17Factory())->createServerRequest('POST', '/notify')
            ->withParsedBody([
                'reseller_id' => 1,
                'notification_type' => 1,
                'client_id' => 1,
                'creator_id' => 1,
                'expert_id' => 1,
                'complaint_id' => 1,
                'complaint_number' => '12345',
                'consumption_id' => 1,
                'consumption_number' => '67890',
                'agreement_number' => '54321',
                'date' => '2023-01-01',
                'differences' => [
                    'from' => 1,
                    'to' => 2
                ]
            ]);

        $this->referencesOperation->expects($this->once())
            ->method('doOperation')
            ->with($this->isInstanceOf(SendNotificationDTO::class))
            ->willReturn([]);

        $response = $this->controller->send($request);

        $this->assertInstanceOf(ResponseInterface::class, $response);
        $this->assertEquals(200, $response->getStatusCode());
        $this->assertJsonStringEqualsJsonString(
            json_encode(['message' => 'handled successfully!', 'data' => []]),
            (string)$response->getBody()
        );
    }
}

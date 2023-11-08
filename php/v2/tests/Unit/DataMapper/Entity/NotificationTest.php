<?php

namespace Tests\Unit\DataMapper\Entity;

use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;
use PHPUnit\Framework\TestCase;

class NotificationTest extends TestCase
{
    public function testComposeMessage(): void
    {
        $expected = "Dear, Bob there is a problem with you order №14765.";
        $template = "Dear, #name# there is a problem with you order №#number#.";
        $params = $this->createMock(ReferencesOperationParams::class);
        $params
            ->method('toArray')
            ->willReturn(['name' => 'Bob', 'number' => 14765]);
        $notification = new Notification();
        $notification->setTemplate($template);
        $message = $notification->composeMessage($params);
        $this->assertSame($expected, $message);
    }
}

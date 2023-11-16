<?php

namespace Tests\Unit\DataMapper\Entity;

use NodaSoft\DataMapper\Entity\Notification;
use NodaSoft\GenericDto\Dto\Dto;
use PHPUnit\Framework\TestCase;

class NotificationTest extends TestCase
{
    public function testComposeMessage(): void
    {
        $expected = "Dear, Bob there is a problem with you order №14765.";
        $template = "Dear, #name# there is a problem with you order №#number#.";
        $params = $this->fakeDto();
        $notification = new Notification();
        $message = $notification->fillTemplate($template, $params);
        $this->assertSame($expected, $message);
    }

    public function testComposes(): void
    {
        $expected = "Dear, Bob there is a problem with you order №14765.";
        $template = "Dear, #name# there is a problem with you order №#number#.";
        $params = $this->fakeDto();
        $notification = new Notification(1, 'Fake', $template, $template);
        $this->assertSame($expected, $notification->composeMessageSubject($params));
        $this->assertSame($expected, $notification->composeMessageBody($params));
    }

    private function fakeDto(): Dto
    {
        return new class ('Bob', 14765) implements Dto
        {
            /** @var string */
            private $number;

            /** @var string */
            private $name;

            public function __construct(string $name, int $number)
            {
                $this->name = $name;
                $this->number = $number;
            }

            public function isValid(): bool
            {
                foreach ($this as $item) {
                    if (is_null($item)) {
                        return false;
                    }
                }
                return true;
            }

            /**
             * @return array<string, mixed>
             */
            public function toArray(): array
            {
                $array = [];
                foreach ($this as $key => $item) {
                    $array[$key] = $item;
                }
                return $array;
            }

            public function getName(): string
            {
                return $this->name;
            }

            public function setName(string $name): void
            {
                $this->name = $name;
            }

            public function getNumber(): string
            {
                return $this->number;
            }

            public function setNumber(string $number): void
            {
                $this->number = $number;
            }
        };
    }
}

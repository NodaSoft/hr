<?php

namespace Nodasoft\Testapp\Events;

use Nodasoft\Testapp\Events\Base\EventInterface;
use Nodasoft\Testapp\Listeners\ChangeReturnStatusEventListener;
use Nodasoft\Testapp\Listeners\MyCustomListenerToMakeLifeEasy;

readonly class ChangeReturnStatusEvent implements EventInterface
{
    public function __construct(
        private array $data,
    )
    {
    }

    public function getData(): array
    {
        return $this->data;
    }

    public function listeners(): array
    {
        return [
            ChangeReturnStatusEventListener::class,
            MyCustomListenerToMakeLifeEasy::class
        ];
    }
}
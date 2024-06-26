<?php

namespace Nodasoft\Testapp\Listeners;


use Nodasoft\Testapp\Events\Base\EventInterface;
use Nodasoft\Testapp\Listeners\Base\ListenerInterface;

class MyCustomListenerToMakeLifeEasy implements ListenerInterface
{

    public function handle(EventInterface $event): void
    {
        // TODO: handle given case.
    }
}
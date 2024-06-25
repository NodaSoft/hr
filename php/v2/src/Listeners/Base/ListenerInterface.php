<?php

namespace Nodasoft\Testapp\Listeners\Base;


use Nodasoft\Testapp\Events\Base\EventInterface;

interface ListenerInterface
{
    public function handle(EventInterface $event): void;
}
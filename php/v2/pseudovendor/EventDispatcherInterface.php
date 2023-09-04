<?php

declare(strict_types=1);

namespace pseudovendor;

interface EventDispatcherInterface
{
    /**
     * @param mixed $eventName
     * @param Event $event
     * @return Event
     */
    public function dispatch(mixed $eventName, Event $event): Event;
}

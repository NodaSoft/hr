<?php

namespace Nodasoft\Testapp\Events\Base;

use Nodasoft\Testapp\Listeners\Base\ListenerInterface;
use PHPUnit\Logging\Exception;
use ReflectionClass;
use ReflectionException;

class EventDispatcher
{
    /**
     * @throws ReflectionException
     */
    public static function dispatch(EventInterface $event): void
    {
        // use any DI service to init the listener from container
        // i will use reflection for testing purpose :)
        foreach ($event->listeners() as $listener) {
            if (class_exists($listener)) {
                $reflectionClass = new ReflectionClass($listener);
                if ($reflectionClass->implementsInterface(ListenerInterface::class)) {
                    $listenerInstance = $reflectionClass->newInstance();
                    $listenerInstance->handle($event);
                } else {
                    throw new Exception("Listener $listener does not implement ListenerInterface.");
                }
            } else {
                throw new Exception("Listener class {$listener} not found.");
            }
        }
    }
}
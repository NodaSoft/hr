<?php

use PHPUnit\Framework\TestCase;

class HelloWorldTest extends TestCase
{
    public function testHelloWorld(): void
    {
        $this->assertTrue(true, "Hello, Test!");
    }
}

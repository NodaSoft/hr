<?php

namespace Tests\Unit\Request;

use NodaSoft\Request\HttpRequest;
use PHPUnit\Framework\TestCase;

class HttpRequestTest extends TestCase
{
    public function testComposeOperationFactoryClassName(): void
    {
        $request = new HttpRequest();
        $this->assertSame(
            'NodaSoft\Operation\Factory\FooBarBazQuzQuuzFactory',
            $request->composeOperationFactoryClassName("/foo/bar_baz/quz_quuz")
        );
    }
}

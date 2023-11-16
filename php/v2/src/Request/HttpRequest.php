<?php

namespace NodaSoft\Request;

class HttpRequest implements Request
{
    /** @var array<string, mixed> */
    private $request;

    public function __construct()
    {
        $this->request = $_REQUEST;
    }

    /**
     * @param string $key
     * @return mixed
     */
    public function getData(string $key)
    {
        return $this->request['data'][$key] ?? null;
    }
}

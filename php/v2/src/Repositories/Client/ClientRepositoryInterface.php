<?php

namespace Nodasoft\Testapp\Repositories\Client;


use Nodasoft\Testapp\Entities\Client\Client;

interface ClientRepositoryInterface
{
    public function getById(int $id): Client;
}
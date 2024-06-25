<?php

namespace Nodasoft\Testapp\Repositories\Status;


use Nodasoft\Testapp\Entities\Status\Status;

interface StatusRepositoryInterface
{

    public function getById(int $id): Status;
}
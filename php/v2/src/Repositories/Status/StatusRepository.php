<?php

namespace Nodasoft\Testapp\Repositories\Status;

use Exception;
use Nodasoft\Testapp\Entities\Status\Status;
use Nodasoft\Testapp\Entities\StatusMockData;
use Nodasoft\Testapp\Traits\CanGetByKey;

class StatusRepository implements StatusRepositoryInterface
{
    use CanGetByKey;

    /**
     * @throws Exception
     */
    public function getById(int $id): Status
    {
        $data = StatusMockData::get();
        $item = $this->getByKeyOrThrow($data, $id);

        return new status(
            $item['id'],
            $item['name'],
        );
    }
}
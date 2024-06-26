<?php

namespace Nodasoft\Testapp\Repositories\Client;

use Exception;
use Nodasoft\Testapp\Entities\Client\Client;
use Nodasoft\Testapp\Entities\ClientMockData;
use Nodasoft\Testapp\Enums\ContactorType;
use Nodasoft\Testapp\Traits\CanGetByKey;

class ClientRepository implements ClientRepositoryInterface
{
    use CanGetByKey;

    /**
     * @throws Exception
     */
    public function getById(int $id): Client
    {
        $record = $this->getByKeyOrThrow(ClientMockData::get(), $id);

        return new Client(
            $record['id'],
            $record['name'],
            $record['seller_id'],
            $record['email'],
            $record['mobile'],
            ContactorType::tryFrom($record['type'])
        );
    }
}
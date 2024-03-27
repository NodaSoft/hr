<?php

namespace Israil\V2\Repositories;

use Israil\V2\Stubs\Contractor;
use Israil\V2\Exceptions\NotFoundException;

class ClientRepository
{
	/**
	 * @throws NotFoundException
	 */
	public function findOrFail(int $id)
	{
		return Contractor::getById($id) ?? throw new NotFoundException('Client not found');
	}
}
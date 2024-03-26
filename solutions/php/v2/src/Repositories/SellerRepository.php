<?php

namespace Israil\V2\Repositories;

use Israil\V2\Stubs\Contractor;
use Israil\V2\Exceptions\NotFoundException;
use Israil\V2\Stubs\Seller;
use function NW\WebService\References\Operations\Notification\getResellerEmailFrom;

class SellerRepository
{
	/**
	 * @throws NotFoundException
	 */
	public function findOrFail(int $id): Contractor
	{
		return Seller::getById($id) ?? throw new NotFoundException('Seller not found');
	}

	public function getEmailFrom(int $id): string
	{
		return getResellerEmailFrom($id);
	}
}
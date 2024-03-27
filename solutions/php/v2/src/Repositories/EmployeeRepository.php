<?php

namespace Israil\V2\Repositories;

use Israil\V2\Stubs\Contractor;
use Israil\V2\Exceptions\NotFoundException;
use Israil\V2\Stubs\Employee;

class EmployeeRepository
{
	/**
	 * @throws NotFoundException
	 */
	public function findOrFail(int $id): Contractor
	{
		return Employee::getById($id) ?? throw new NotFoundException('Creator not found');
	}
}
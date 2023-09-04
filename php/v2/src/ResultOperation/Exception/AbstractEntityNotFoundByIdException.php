<?php

declare(strict_types=1);

namespace ResultOperation\Exception;

use Exception;

abstract class AbstractEntityNotFoundByIdException extends Exception
{
    public function __construct(mixed $id)
    {
        /**
         * Может упасть, если в {@var $id} пихнут какой-нибудь неконвертируемый в string ужас
         */
        parent::__construct(
            sprintf(
                'Cant find entity by id: %s',
                $id
            )
        );
    }
}

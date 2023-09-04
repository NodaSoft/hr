<?php

declare(strict_types=1);

namespace pseudovendor;

class BaseRepository
{
    /**
     * @param int $id
     * @return ?BaseEntity
     */
    public function get(int $id): ?BaseEntity
    {
        return new BaseEntity(); // stub
    }
}

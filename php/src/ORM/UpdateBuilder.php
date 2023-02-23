<?php

namespace App\ORM;

class UpdateBuilder extends AbstractBuilder
{
    protected function preBuild(): string
    {
        $sql = "UPDATE `{$this->queryBuilder->getFrom()}` SET ";
        $keys = array_keys($this->queryBuilder->getParams());
        $sql .= join(', ', array_map(fn($key) => "$key = :$key", $keys));
        return $sql;
    }
}
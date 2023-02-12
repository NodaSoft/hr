<?php

namespace App\ORM;

class UpdateBuilder extends AbstractBuilder
{
    protected function preBuild(): string
    {
        $sql = "UPDATE `{$this->queryBuilder->getFrom()}` SET ";

        foreach(array_keys($this->queryBuilder->getParams()) as $key) {
            $col = substr($key, 1);
            $sql .= " $col = $key ";
        }

        return $sql;
    }
}
<?php

namespace App\ORM;

class InsertBuilder extends AbstractBuilder
{
    protected function preBuild(): string
    {
        $cols = [];
        $keys = [];

        foreach(array_keys($this->queryBuilder->getParams()) as $col) {
            $cols[] = $col;
            $keys[] = ":$col";
        }

        $cols = join(', ', $cols);
        $keys = join(', ', $keys);
        return "INSERT INTO `{$this->queryBuilder->getFrom()}` ($cols) VALUES ($keys)";
    }
}